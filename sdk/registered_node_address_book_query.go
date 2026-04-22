package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	registeredNodeQueryRetryDelay = 200 * time.Millisecond
	registeredNodeMaxPages        = 1000
)

type RegisteredNode struct {
	AdminKey         Key
	Description      string
	RegisteredNodeID uint64
	ServiceEndpoints []RegisteredServiceEndpoint
	CreatedTimestamp string
}

type RegisteredNodeAddressBook struct {
	RegisteredNodes []RegisteredNode
}

// RegisteredNodeAddressBookQuery defines the contract for querying registered nodes.
type RegisteredNodeAddressBookQuery struct {
	registeredNodeId *uint64
	maxAttempts      uint64
	limit            int32
}

func NewRegisteredNodeAddressBookQuery() *RegisteredNodeAddressBookQuery {
	return &RegisteredNodeAddressBookQuery{}
}

// SetRegisteredNodeId filters the query to a specific registered node ID.
func (q *RegisteredNodeAddressBookQuery) SetRegisteredNodeId(id uint64) *RegisteredNodeAddressBookQuery {
	q.registeredNodeId = &id
	return q
}

// GetRegisteredNodeId returns the registered node ID filter, or 0 if not set.
func (q *RegisteredNodeAddressBookQuery) GetRegisteredNodeId() uint64 {
	if q.registeredNodeId == nil {
		return 0
	}
	return *q.registeredNodeId
}

// SetLimit sets the maximum number of registered nodes to return.
// Zero (the default) leaves the limit up to the mirror node.
func (q *RegisteredNodeAddressBookQuery) SetLimit(limit int32) *RegisteredNodeAddressBookQuery {
	q.limit = limit
	return q
}

// GetLimit returns the current limit, or 0 if unset.
func (q *RegisteredNodeAddressBookQuery) GetLimit() int32 {
	return q.limit
}

// SetMaxAttempts sets the total number of attempts (initial try + retries).
// Zero (the default) is treated as a single attempt with no retries.
func (q *RegisteredNodeAddressBookQuery) SetMaxAttempts(maxAttempts uint64) *RegisteredNodeAddressBookQuery {
	q.maxAttempts = maxAttempts
	return q
}

// GetMaxAttempts returns the configured retry budget.
func (q *RegisteredNodeAddressBookQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

func (q *RegisteredNodeAddressBookQuery) Execute(client *Client) (RegisteredNodeAddressBook, error) {
	if client == nil {
		return RegisteredNodeAddressBook{}, errNoClientProvided
	}

	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return RegisteredNodeAddressBook{}, fmt.Errorf("mirror node is not set")
	}

	mirrorUrl, err := client.GetMirrorRestApiBaseUrl()
	if err != nil {
		return RegisteredNodeAddressBook{}, fmt.Errorf("failed to get mirror REST API base URL: %w", err)
	}

	if strings.Contains(mirrorUrl, "localhost") || strings.Contains(mirrorUrl, "127.0.0.1") {
		mirrorUrl = "http://localhost:8084/api/v1"
	}

	mirrorBase, err := url.Parse(mirrorUrl)
	if err != nil {
		return RegisteredNodeAddressBook{}, fmt.Errorf("invalid mirror REST API base URL %q: %w", mirrorUrl, err)
	}

	attempts := q.maxAttempts
	if attempts == 0 {
		if clientMax := client.GetMaxAttempts(); clientMax > 0 {
			attempts = uint64(clientMax)
		} else {
			attempts = 1
		}
	}

	endpoint := q.buildURL(mirrorUrl)
	allNodes := make([]RegisteredNode, 0)

	for range registeredNodeMaxPages {
		body, err := fetchRegisteredNodesPage(endpoint, attempts)
		if err != nil {
			return RegisteredNodeAddressBook{}, err
		}

		nodes, next, err := parseRegisteredNodes(body)
		if err != nil {
			return RegisteredNodeAddressBook{}, err
		}
		allNodes = append(allNodes, nodes...)

		if next == nil || *next == "" {
			return RegisteredNodeAddressBook{RegisteredNodes: allNodes}, nil
		}

		resolved, err := resolveNextURL(mirrorBase, *next)
		if err != nil {
			return RegisteredNodeAddressBook{}, fmt.Errorf("invalid pagination next link %q: %w", *next, err)
		}
		endpoint = resolved
	}

	return RegisteredNodeAddressBook{}, fmt.Errorf("exceeded pagination cap of %d pages", registeredNodeMaxPages)
}

// buildURL composes the mirror node REST URL together with any query
// parameters configured on the query.
func (q *RegisteredNodeAddressBookQuery) buildURL(mirrorBaseURL string) string {
	endpoint := fmt.Sprintf("%s/network/registered-nodes", mirrorBaseURL)

	params := url.Values{}
	if q.registeredNodeId != nil {
		params.Set("registerednode.id", strconv.FormatUint(*q.registeredNodeId, 10))
	}
	if q.limit > 0 {
		params.Set("limit", strconv.FormatInt(int64(q.limit), 10))
	}

	if encoded := params.Encode(); encoded != "" {
		endpoint = endpoint + "?" + encoded
	}
	return endpoint
}

// fetchRegisteredNodes issues a single GET against the mirror node and
// returns the response body together with the HTTP status code.
func fetchRegisteredNodes(endpoint string) ([]byte, int, error) {
	resp, err := http.Get(endpoint) // #nosec
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, resp.StatusCode, nil
}

// fetchRegisteredNodesPage wraps fetchRegisteredNodes with the per-page retry.
func fetchRegisteredNodesPage(endpoint string, attempts uint64) ([]byte, error) {
	var lastErr error
	for attempt := range attempts {
		if attempt > 0 {
			time.Sleep(registeredNodeQueryRetryDelay)
		}

		body, status, err := fetchRegisteredNodes(endpoint)
		if err != nil {
			lastErr = err
			continue
		}

		if status >= 500 {
			lastErr = fmt.Errorf("received non-200 response from mirror node: %d, details: %s", status, body)
			continue
		}

		if status != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response from mirror node: %d, details: %s", status, body)
		}

		return body, nil
	}
	return nil, fmt.Errorf("failed after %d attempt(s): %w", attempts, lastErr)
}

// resolveNextURL resolves a pagination next link against the mirror base URL.
func resolveNextURL(base *url.URL, next string) (string, error) {
	parsed, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(parsed).String(), nil
}

func parseRegisteredNodes(body []byte) ([]RegisteredNode, *string, error) {
	var raw registeredNodesResponseJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	nodes := make([]RegisteredNode, 0, len(raw.RegisteredNodes))
	for _, rn := range raw.RegisteredNodes {
		node, err := registeredNodeFromJSON(rn)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert registered node %d: %w", rn.RegisteredNodeID, err)
		}
		nodes = append(nodes, node)
	}

	var next *string
	if raw.Links != nil {
		next = raw.Links.Next
	}
	return nodes, next, nil
}

type registeredNodesResponseJSON struct {
	RegisteredNodes []registeredNodeJSON `json:"registered_nodes"`
	Links           *linksJSON           `json:"links"`
}

type linksJSON struct {
	Next *string `json:"next"`
}

type registeredNodeJSON struct {
	AdminKey         *adminKeyJSON         `json:"admin_key"`
	CreatedTimestamp string                `json:"created_timestamp"`
	Description      string                `json:"description"`
	RegisteredNodeID uint64                `json:"registered_node_id"`
	ServiceEndpoints []serviceEndpointJSON `json:"service_endpoints"`
}

type adminKeyJSON struct {
	Type string `json:"_type"`
	Key  string `json:"key"`
}

type serviceEndpointJSON struct {
	BlockNode      *blockNodeJSON      `json:"block_node"`
	GeneralService *generalServiceJSON `json:"general_service"`
	DomainName     *string             `json:"domain_name"`
	IPAddress      *string             `json:"ip_address"`
	Port           uint32              `json:"port"`
	RequiresTls    bool                `json:"requires_tls"`
	Type           string              `json:"type"`
}

type blockNodeJSON struct {
	EndpointApis []string `json:"endpoint_apis"`
}

type generalServiceJSON struct {
	Description string `json:"description"`
}

func registeredNodeFromJSON(raw registeredNodeJSON) (RegisteredNode, error) {
	node := RegisteredNode{
		Description:      raw.Description,
		RegisteredNodeID: raw.RegisteredNodeID,
		CreatedTimestamp: raw.CreatedTimestamp,
	}

	if raw.AdminKey != nil {
		key, err := adminKeyFromJSON(*raw.AdminKey)
		if err != nil {
			return RegisteredNode{}, fmt.Errorf("failed to parse admin key: %w", err)
		}
		node.AdminKey = key
	}

	endpoints := make([]RegisteredServiceEndpoint, 0, len(raw.ServiceEndpoints))
	for i, ep := range raw.ServiceEndpoints {
		endpoint, err := serviceEndpointFromJSON(ep)
		if err != nil {
			return RegisteredNode{}, fmt.Errorf("failed to parse service endpoint %d: %w", i, err)
		}
		endpoints = append(endpoints, endpoint)
	}
	node.ServiceEndpoints = endpoints

	return node, nil
}

func adminKeyFromJSON(raw adminKeyJSON) (Key, error) {
	switch strings.ToUpper(raw.Type) {
	case "ED25519":
		return PublicKeyFromStringEd25519(raw.Key)
	case "ECDSA_SECP256K1":
		return PublicKeyFromStringECDSA(raw.Key)
	default:
		return PublicKeyFromString(raw.Key)
	}
}

func serviceEndpointFromJSON(raw serviceEndpointJSON) (RegisteredServiceEndpoint, error) {
	base := registeredEndpointBase{
		port:        raw.Port,
		requiresTls: raw.RequiresTls,
	}

	if raw.IPAddress != nil && *raw.IPAddress != "" {
		ip := net.ParseIP(*raw.IPAddress)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", *raw.IPAddress)
		}
		if v4 := ip.To4(); v4 != nil {
			base.ipAddress = v4
		} else {
			base.ipAddress = ip.To16()
		}
	}

	if raw.DomainName != nil && *raw.DomainName != "" {
		base.domainName = *raw.DomainName
	}

	switch strings.ToUpper(raw.Type) {
	case "BLOCK_NODE":
		endpoint := &BlockNodeServiceEndpoint{
			registeredEndpointBase: base,
		}
		if raw.BlockNode != nil {
			apis := make([]BlockNodeApi, 0, len(raw.BlockNode.EndpointApis))
			for _, apiStr := range raw.BlockNode.EndpointApis {
				apis = append(apis, blockNodeApiFromString(apiStr))
			}
			endpoint.endpointApis = apis
		}
		return endpoint, nil

	case "MIRROR_NODE":
		return &MirrorNodeServiceEndpoint{
			registeredEndpointBase: base,
		}, nil

	case "RPC_RELAY":
		return &RpcRelayServiceEndpoint{
			registeredEndpointBase: base,
		}, nil

	case "GENERAL_SERVICE":
		endpoint := &GeneralServiceEndpoint{
			registeredEndpointBase: base,
		}
		if raw.GeneralService != nil {
			endpoint.description = raw.GeneralService.Description
		}
		return endpoint, nil

	default:
		return nil, fmt.Errorf("unknown endpoint type: %s", raw.Type)
	}
}

func blockNodeApiFromString(s string) BlockNodeApi {
	switch strings.ToUpper(s) {
	case "STATUS":
		return BlockNodeApiStatus
	case "PUBLISH":
		return BlockNodeApiPublish
	case "SUBSCRIBE_STREAM":
		return BlockNodeApiSubscribeStream
	case "STATE_PROOF":
		return BlockNodeApiStateProof
	default:
		return BlockNodeApiOther
	}
}
