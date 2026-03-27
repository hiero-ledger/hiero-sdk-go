package hiero

// SPDX-License-Identifier: Apache-2.0

// BlockNodeApi is an enumeration of well-known block node endpoint APIs.
type BlockNodeApi int32

const (
	// BlockNodeApiOther represents an unknown or custom block node API type.
	BlockNodeApiOther BlockNodeApi = 0
	// BlockNodeApiStatus represents the Block Node Status API.
	BlockNodeApiStatus BlockNodeApi = 1
	// BlockNodeApiPublish represents the Block Node Publish API.
	BlockNodeApiPublish BlockNodeApi = 2
	// BlockNodeApiSubscribeStream represents the Block Node Subscribe Stream API.
	BlockNodeApiSubscribeStream BlockNodeApi = 3
	// BlockNodeApiStateProof represents the Block Node State Proof API.
	BlockNodeApiStateProof BlockNodeApi = 4
)

// String returns the string representation of the BlockNodeApi.
func (api BlockNodeApi) String() string {
	switch api {
	case BlockNodeApiOther:
		return "OTHER"
	case BlockNodeApiStatus:
		return "STATUS"
	case BlockNodeApiPublish:
		return "PUBLISH"
	case BlockNodeApiSubscribeStream:
		return "SUBSCRIBE_STREAM"
	case BlockNodeApiStateProof:
		return "STATE_PROOF"
	}

	return unknownString
}
