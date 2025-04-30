package hiero

// SPDX-License-Identifier: Apache-2.0

type NetworkName string

const (
	NetworkNameMainnet    NetworkName = "mainnet"
	NetworkNameTestnet    NetworkName = "testnet"
	NetworkNamePreviewnet NetworkName = "previewnet"
	NetworkNameOther      NetworkName = "other"
)

// Deprecated
func (networkName NetworkName) String() string { //nolint
	switch networkName {
	case NetworkNameMainnet:
		return "mainnet" //nolint
	case NetworkNameTestnet:
		return "testnet" //nolint
	case NetworkNamePreviewnet:
		return "previewnet" //nolint
	case NetworkNameOther:
		return "other"
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}

// Deprecated
func NetworkNameFromString(s string) NetworkName { //nolint
	switch s {
	case "mainnet": //nolint
		return NetworkNameMainnet
	case "testnet": //nolint
		return NetworkNameTestnet
	case "previewnet": //nolint
		return NetworkNamePreviewnet
	case "other": //nolint
		return NetworkNameOther
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}
