package hiero

// SPDX-License-Identifier: Apache-2.0

import "fmt"

// ClprSignatureScheme is the scheme of a CLPR ownership key (HIP-1535).
type ClprSignatureScheme int32

const (
	ClprSignatureSchemeEcdsaSecp256k1 ClprSignatureScheme = 0
	ClprSignatureSchemeEd25519        ClprSignatureScheme = 1
)

// String returns the wire name of the scheme, or "UNKNOWN (n)" for a code this SDK does not know.
func (scheme ClprSignatureScheme) String() string {
	switch scheme {
	case ClprSignatureSchemeEcdsaSecp256k1:
		return "ECDSA_SECP256K1"
	case ClprSignatureSchemeEd25519:
		return "ED25519"
	}

	return fmt.Sprintf("UNKNOWN (%d)", int32(scheme))
}
