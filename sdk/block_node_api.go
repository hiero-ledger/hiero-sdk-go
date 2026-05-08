package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
)

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

// blockNodeApiFromString parses a BlockNodeApi enum name. Returns an error
// for anything outside the known set, including "UNKNOWN".
func blockNodeApiFromString(s string) (BlockNodeApi, error) {
	switch strings.ToUpper(s) {
	case "OTHER":
		return BlockNodeApiOther, nil
	case "STATUS":
		return BlockNodeApiStatus, nil
	case "PUBLISH":
		return BlockNodeApiPublish, nil
	case "SUBSCRIBE_STREAM":
		return BlockNodeApiSubscribeStream, nil
	case "STATE_PROOF":
		return BlockNodeApiStateProof, nil
	default:
		return BlockNodeApiOther, fmt.Errorf("unknown block node API: %q", s)
	}
}
