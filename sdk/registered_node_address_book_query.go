package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
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

// RegisteredNodeAddressBookQuery defines the contract for querying
// registered nodes. Transport, endpoint path, pagination and filtering depend
// on the mirror node implementation and may change as it evolves.
type RegisteredNodeAddressBookQuery struct {
	registeredNodeId *uint64
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

	isLocalHost := strings.Contains(mirrorUrl, "localhost") || strings.Contains(mirrorUrl, "127.0.0.1")
	if isLocalHost {
		mirrorUrl = "http://localhost:8084/api/v1"
	}

	url := fmt.Sprintf("%s/network/registered-nodes", mirrorUrl)
	if q.registeredNodeId != nil {
		url = fmt.Sprintf("%s?registerednode.id=%d", url, *q.registeredNodeId)
	}

	resp, err := http.Get(url) // #nosec
	if err != nil {
		return RegisteredNodeAddressBook{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisteredNodeAddressBook{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return RegisteredNodeAddressBook{}, fmt.Errorf("received non-200 response from mirror node: %d, details: %s", resp.StatusCode, body)
	}

	var raw registeredNodesResponseJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return RegisteredNodeAddressBook{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	nodes := make([]RegisteredNode, 0, len(raw.RegisteredNodes))
	for _, rn := range raw.RegisteredNodes {
		node, err := registeredNodeFromJSON(rn)
		if err != nil {
			return RegisteredNodeAddressBook{}, fmt.Errorf("failed to convert registered node %d: %w", rn.RegisteredNodeID, err)
		}
		nodes = append(nodes, node)
	}

	return RegisteredNodeAddressBook{
		RegisteredNodes: nodes,
	}, nil
}

type registeredNodesResponseJSON struct {
	RegisteredNodes []registeredNodeJSON `json:"registered_nodes"`
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
	BlockNode   *blockNodeJSON `json:"block_node"`
	DomainName  *string        `json:"domain_name"`
	IPAddress   *string        `json:"ip_address"`
	Port        uint32         `json:"port"`
	RequiresTls bool           `json:"requires_tls"`
	Type        string         `json:"type"`
}

type blockNodeJSON struct {
	EndpointApis []string `json:"endpoint_apis"`
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
