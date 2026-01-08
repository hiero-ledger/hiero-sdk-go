package param

// SPDX-License-Identifier: Apache-2.0

type EndpointParams struct {
	Address    *string `json:"ipAddressV4,omitempty"`
	Port       *int32  `json:"port,omitempty"`
	DomainName *string `json:"domainName,omitempty"`
}

type CreateNodeParams struct {
	BaseParams
	AccountId            *string           `json:"accountId,omitempty"`
	Description          *string           `json:"description,omitempty"`
	GossipEndpoints      *[]EndpointParams `json:"gossipEndpoints,omitempty"`
	ServiceEndpoints     *[]EndpointParams `json:"serviceEndpoints,omitempty"`
	GossipCaCertificate  *string           `json:"gossipCaCertificate,omitempty"`
	GrpcCertificateHash  *string           `json:"grpcCertificateHash,omitempty"`
	AdminKey             *string           `json:"adminKey,omitempty"`
	DeclineReward        *bool             `json:"declineReward,omitempty"`
	GrpcWebProxyEndpoint *EndpointParams   `json:"grpcWebProxyEndpoint,omitempty"`
}
type UpdateNodeParams struct {
	BaseParams
	NodeId               *string           `json:"nodeId,omitempty"`
	AccountId            *string           `json:"accountId,omitempty"`
	Description          *string           `json:"description,omitempty"`
	GossipEndpoints      *[]EndpointParams `json:"gossipEndpoints,omitempty"`
	ServiceEndpoints     *[]EndpointParams `json:"serviceEndpoints,omitempty"`
	GossipCaCertificate  *string           `json:"gossipCaCertificate,omitempty"`
	GrpcCertificateHash  *string           `json:"grpcCertificateHash,omitempty"`
	AdminKey             *string           `json:"adminKey,omitempty"`
	DeclineReward        *bool             `json:"declineReward,omitempty"`
	GrpcWebProxyEndpoint *EndpointParams   `json:"grpcWebProxyEndpoint,omitempty"`
}

type DeleteNodeParams struct {
	BaseParams
	NodeId *string `json:"nodeId,omitempty"`
}
