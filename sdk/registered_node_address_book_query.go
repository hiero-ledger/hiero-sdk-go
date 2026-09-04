package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const registeredNodeMaxPages = 1000

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

	baseURL, err := q.resolveBaseURL(client)
	if err != nil {
		return RegisteredNodeAddressBook{}, err
	}

	path, err := q.buildPath()
	if err != nil {
		return RegisteredNodeAddressBook{}, err
	}

	options := client.mirrorHttpSettings()
	options.maxAttempts = int(q.resolveAttempts(client))

	return q.walkPages(client.mirrorRestClientForBaseURL(baseURL, options), path)
}

// resolveBaseURL returns the REST base. A local node serves this endpoint on 8084 rather than
// on the client's mirror REST port.
func (q *RegisteredNodeAddressBookQuery) resolveBaseURL(client *Client) (string, error) {
	mirrorUrl, err := mirrorNodeRestBaseURL(client)
	if err != nil {
		return "", err
	}

	if strings.Contains(mirrorUrl, "localhost") || strings.Contains(mirrorUrl, "127.0.0.1") {
		mirrorUrl = "http://localhost:8084/api/v1"
	}

	return mirrorUrl, nil
}

// resolveEndpoint returns the initial query URL.
func (q *RegisteredNodeAddressBookQuery) resolveEndpoint(client *Client) (string, error) {
	baseURL, err := q.resolveBaseURL(client)
	if err != nil {
		return "", err
	}

	return q.buildURL(baseURL), nil
}

// resolveAttempts picks the per-page retry budget: query setting first,
// client default second, single attempt as the final fallback.
func (q *RegisteredNodeAddressBookQuery) resolveAttempts(client *Client) uint64 {
	if q.maxAttempts > 0 {
		return q.maxAttempts
	}
	if clientMax := client.GetMaxAttempts(); clientMax > 0 {
		return uint64(clientMax)
	}
	return 1
}

// walkPages follows links.next until exhausted (or the page cap trips). Each page goes through
// the shared mirror HTTP layer, so pagination inherits the same retry policy as the first
// request instead of the fixed-delay loop this used to carry.
//
// A next link is turned back into a path, so a page URL can only ever address the configured
// mirror node — see nextPagePath.
func (q *RegisteredNodeAddressBookQuery) walkPages(restClient *mirrorHttpClient, startPath mirrorRestPath) (RegisteredNodeAddressBook, error) {
	allNodes := make([]RegisteredNode, 0)
	path := startPath

	for range registeredNodeMaxPages {
		body, err := fetchRegisteredNodesPage(restClient, path)
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

		path, err = nextPagePath(*next)
		if err != nil {
			return RegisteredNodeAddressBook{}, fmt.Errorf("invalid pagination next link %q: %w", *next, err)
		}
	}

	return RegisteredNodeAddressBook{}, fmt.Errorf("exceeded pagination cap of %d pages", registeredNodeMaxPages)
}

// buildPath composes the endpoint path together with any query parameters configured on the
// query. It is a path, never a URL, so no call site above this can name a host.
func (q *RegisteredNodeAddressBookQuery) buildPath() (mirrorRestPath, error) {
	path := "/network/registered-nodes"

	params := url.Values{}
	if q.registeredNodeId != nil {
		params.Set("registerednode.id", strconv.FormatUint(*q.registeredNodeId, 10))
	}
	if q.limit > 0 {
		params.Set("limit", strconv.FormatInt(int64(q.limit), 10))
	}

	if encoded := params.Encode(); encoded != "" {
		path = path + "?" + encoded
	}

	return newMirrorRestPath(path)
}

// buildURL resolves buildPath against a mirror base URL.
func (q *RegisteredNodeAddressBookQuery) buildURL(mirrorBaseURL string) string {
	path, err := q.buildPath()
	if err != nil {
		return ""
	}

	return resolveMirrorPath(mirrorBaseURL, path)
}

// fetchRegisteredNodesPage requests one page. Retry, backoff and classification all belong to
// the mirror HTTP layer now; this used to carry its own loop with a fixed 200 ms delay.
func fetchRegisteredNodesPage(restClient *mirrorHttpClient, path mirrorRestPath) ([]byte, error) {
	resp, err := restClient.get(context.Background(), path)
	if err != nil {
		if errors.Is(err, errMirrorHttpRetriesExhausted) {
			return nil, mirrorNodeStatusError(resp)
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if resp.statusCode != http.StatusOK {
		return nil, mirrorNodeStatusError(resp)
	}

	return resp.body, nil
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
	base, err := endpointBaseFromJSON(raw)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(raw.Type) {
	case "BLOCK_NODE":
		return blockNodeEndpointFromJSON(base, raw.BlockNode)
	case "MIRROR_NODE":
		return &MirrorNodeServiceEndpoint{registeredEndpointBase: base}, nil
	case "RPC_RELAY":
		return &RpcRelayServiceEndpoint{registeredEndpointBase: base}, nil
	case "GENERAL_SERVICE":
		return generalServiceEndpointFromJSON(base, raw.GeneralService), nil
	default:
		return nil, fmt.Errorf("unknown endpoint type: %s", raw.Type)
	}
}

// endpointBaseFromJSON extracts the IP / domain / port / TLS fields shared by
// every endpoint subtype.
func endpointBaseFromJSON(raw serviceEndpointJSON) (registeredEndpointBase, error) {
	base := registeredEndpointBase{
		port:        raw.Port,
		requiresTls: raw.RequiresTls,
	}

	if raw.IPAddress != nil && *raw.IPAddress != "" {
		ip := net.ParseIP(*raw.IPAddress)
		if ip == nil {
			return base, fmt.Errorf("invalid IP address: %s", *raw.IPAddress)
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

	return base, nil
}

func blockNodeEndpointFromJSON(base registeredEndpointBase, raw *blockNodeJSON) (*BlockNodeServiceEndpoint, error) {
	endpoint := &BlockNodeServiceEndpoint{registeredEndpointBase: base}
	if raw == nil {
		return endpoint, nil
	}

	endpoint.endpointApis = make([]BlockNodeApi, 0, len(raw.EndpointApis))
	for _, apiStr := range raw.EndpointApis {
		api, err := blockNodeApiFromString(apiStr)
		if err != nil {
			return nil, err
		}
		endpoint.endpointApis = append(endpoint.endpointApis, api)
	}
	return endpoint, nil
}

func generalServiceEndpointFromJSON(base registeredEndpointBase, raw *generalServiceJSON) *GeneralServiceEndpoint {
	endpoint := &GeneralServiceEndpoint{registeredEndpointBase: base}
	if raw != nil {
		endpoint.description = raw.Description
	}
	return endpoint
}
