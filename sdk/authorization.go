package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// Authorization represents a single EIP-7702 authorization tuple
// [chainId, address, nonce, yParity, r, s], authorizing the delegation of an
// EOA's code to a contract address. Introduced for SDKs in HIP-1340.
type Authorization struct {
	chainId []byte
	address []byte
	nonce   uint64
	yParity uint32
	r       []byte
	s       []byte
}

// NewAuthorization creates a new Authorization with the provided fields.
func NewAuthorization(chainId, address []byte, nonce uint64, yParity uint32, r, s []byte) Authorization {
	return Authorization{
		chainId: chainId,
		address: address,
		nonce:   nonce,
		yParity: yParity,
		r:       r,
		s:       s,
	}
}

// GetChainId returns the chain id for this authorization.
func (a *Authorization) GetChainId() []byte { return a.chainId }

// SetChainId sets the chain id for this authorization.
func (a *Authorization) SetChainId(chainId []byte) *Authorization {
	a.chainId = chainId
	return a
}

// GetAddress returns the delegated contract address.
func (a *Authorization) GetAddress() []byte { return a.address }

// SetAddress sets the delegated contract address.
func (a *Authorization) SetAddress(address []byte) *Authorization {
	a.address = address
	return a
}

// GetNonce returns the authorization nonce.
func (a *Authorization) GetNonce() uint64 { return a.nonce }

// SetNonce sets the authorization nonce.
func (a *Authorization) SetNonce(nonce uint64) *Authorization {
	a.nonce = nonce
	return a
}

// GetYParity returns the signature y-parity (recovery id).
func (a *Authorization) GetYParity() uint32 { return a.yParity }

// SetYParity sets the signature y-parity (recovery id).
func (a *Authorization) SetYParity(yParity uint32) *Authorization {
	a.yParity = yParity
	return a
}

// GetR returns the R signature component.
func (a *Authorization) GetR() []byte { return a.r }

// SetR sets the R signature component.
func (a *Authorization) SetR(r []byte) *Authorization {
	a.r = r
	return a
}

// GetS returns the S signature component.
func (a *Authorization) GetS() []byte { return a.s }

// SetS sets the S signature component.
func (a *Authorization) SetS(s []byte) *Authorization {
	a.s = s
	return a
}

// String returns a string representation of the Authorization.
func (a Authorization) String() string {
	return fmt.Sprintf("{ChainId: %s, Address: %s, Nonce: %d, YParity: %d, R: %s, S: %s}",
		hex.EncodeToString(a.chainId),
		hex.EncodeToString(a.address),
		a.nonce,
		a.yParity,
		hex.EncodeToString(a.r),
		hex.EncodeToString(a.s),
	)
}

// _toRLPItem encodes the authorization as a 6-element RLP list, with nonce and
// yParity in canonical minimal big-endian form.
func (a Authorization) _toRLPItem() *RLPItem {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(a.chainId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(a.address))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(_uint64ToEthBytes(a.nonce)))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(_uint64ToEthBytes(uint64(a.yParity))))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(a.r))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(a.s))
	return item
}

// _authorizationFromRLPItem decodes a 6-element RLP list into an Authorization.
func _authorizationFromRLPItem(item *RLPItem) (Authorization, error) {
	if item.itemType != LIST_TYPE {
		return Authorization{}, errors.New("invalid authorization list entry: must be a list")
	}
	if len(item.childItems) != 6 {
		return Authorization{}, errors.New("invalid authorization list entry: must be [chainId, address, nonce, yParity, r, s]")
	}
	return NewAuthorization(
		item.childItems[0].itemValue,
		item.childItems[1].itemValue,
		_ethBytesToUint64(item.childItems[2].itemValue),
		uint32(_ethBytesToUint64(item.childItems[3].itemValue)),
		item.childItems[4].itemValue,
		item.childItems[5].itemValue,
	), nil
}
