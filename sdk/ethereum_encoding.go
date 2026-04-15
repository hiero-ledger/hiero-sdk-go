package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/binary"
	"math/big"
)

// Helpers for Ethereum RLP canonical integer encoding: big-endian, no
// leading zeros, zero is the empty byte string.

// _uint64ToEthBytes encodes v as canonical minimal big-endian. Zero → empty.
func _uint64ToEthBytes(v uint64) []byte {
	if v == 0 {
		return []byte{}
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	i := 0
	for i < len(buf) && buf[i] == 0 {
		i++
	}
	return buf[i:]
}

// _ethBytesToUint64 decodes canonical bytes to uint64. Empty → 0.
// Input longer than 8 bytes is truncated to the low 8.
func _ethBytesToUint64(b []byte) uint64 {
	if len(b) == 0 {
		return 0
	}
	if len(b) > 8 {
		b = b[len(b)-8:]
	}
	var buf [8]byte
	copy(buf[8-len(b):], b)
	return binary.BigEndian.Uint64(buf[:])
}

// _bigIntToEthBytes encodes v as canonical minimal big-endian.
// Nil and zero → empty. Negative values encode their absolute value.
func _bigIntToEthBytes(v *big.Int) []byte {
	if v == nil || v.Sign() == 0 {
		return []byte{}
	}
	if v.Sign() < 0 {
		return new(big.Int).Abs(v).Bytes()
	}
	return v.Bytes()
}

// _ethBytesToBigInt decodes canonical bytes to *big.Int. Empty → 0.
func _ethBytesToBigInt(b []byte) *big.Int {
	if len(b) == 0 {
		return new(big.Int)
	}
	return new(big.Int).SetBytes(b)
}
