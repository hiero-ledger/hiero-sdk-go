package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	"github.com/youmark/pkcs8"
)

// _Ed25519PrivateKey is an ed25519 private key.
type _Ed25519PrivateKey struct {
	keyData   []byte
	chainCode []byte
}

func _GenerateEd25519PrivateKey() (*_Ed25519PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	return &_Ed25519PrivateKey{
		keyData: privateKey,
	}, nil
}

// Deprecated
func GeneratePrivateKey() (PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: &_Ed25519PrivateKey{
			keyData: privateKey,
		},
	}, nil
}

// _Ed25519PrivateKeyFromBytes constructs an _Ed25519PrivateKey from a raw slice of either 32 or 64 bytes.
func _Ed25519PrivateKeyFromBytes(bytes []byte) (*_Ed25519PrivateKey, error) {
	length := len(bytes)
	switch length {
	case 48:
		return _Ed25519PrivateKeyFromBytesDer(bytes)
	case 32:
		return _Ed25519PrivateKeyFromBytesRaw(bytes)
	case 64:
		return _Ed25519PrivateKeyFromBytesRaw(bytes)
	default:
		return &_Ed25519PrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(bytes))
	}
}

// _Ed25519PrivateKeyFromBytes constructs an _Ed25519PrivateKey from a raw slice of either 32 or 64 bytes.
func _Ed25519PrivateKeyFromBytesRaw(bytes []byte) (*_Ed25519PrivateKey, error) {
	return &_Ed25519PrivateKey{
		keyData: ed25519.NewKeyFromSeed(bytes[0:32]),
	}, nil
}

// _Ed25519PrivateKeyFromBytes constructs an _Ed25519PrivateKey from a raw slice of either 32 or 64 bytes.
func _Ed25519PrivateKeyFromBytesDer(byt []byte) (*_Ed25519PrivateKey, error) {
	given := hex.EncodeToString(byt)
	result := strings.ReplaceAll(given, _Ed25519PrivateKeyPrefix, "")
	decoded, err := hex.DecodeString(result)
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	length := len(decoded)
	if length != 32 && length != 64 {
		return &_Ed25519PrivateKey{}, _NewErrBadKeyf("invalid private key length: %v byt", len(byt))
	}

	return &_Ed25519PrivateKey{
		keyData: ed25519.NewKeyFromSeed(decoded[0:32]),
	}, nil
}

func _Ed25519PrivateKeyFromSeed(seed []byte) (*_Ed25519PrivateKey, error) {
	h := hmac.New(sha512.New, []byte("ed25519 seed"))

	_, err := h.Write(seed)
	if err != nil {
		return nil, err
	}

	digest := h.Sum(nil)

	keyBytes := digest[0:32]
	chainCode := digest[32:]
	privateKey, err := _Ed25519PrivateKeyFromBytes(keyBytes)

	if err != nil {
		return nil, err
	}
	privateKey.chainCode = chainCode

	return privateKey, nil
}

// Deprecated
// PrivateKeyFromMnemonic recovers an _Ed25519PrivateKey from a valid 24 word length mnemonic phrase and a
// passphrase.
//
// An empty string can be passed for passPhrase If the mnemonic phrase wasn't generated with a passphrase. This is
// required to recover a private key from a mnemonic generated by the Android and iOS wallets.
func _Ed25519PrivateKeyFromMnemonic(mnemonic Mnemonic, passPhrase string) (*_Ed25519PrivateKey, error) {
	seed := mnemonic._ToSeed(passPhrase)
	keyFromSeed, err := _Ed25519PrivateKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	keyBytes := keyFromSeed.keyData
	chainCode := keyFromSeed.chainCode

	// note the index is for derivation, not the index of the slice
	for _, index := range []uint32{44, 3030, 0, 0} {
		keyBytes, chainCode, err = _DeriveEd25519ChildKey(keyBytes, chainCode, index)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := _Ed25519PrivateKeyFromBytes(keyBytes)

	if err != nil {
		return nil, err
	}

	privateKey.chainCode = chainCode

	return privateKey, nil
}

func _Ed25519PrivateKeyFromString(s string) (*_Ed25519PrivateKey, error) {
	bytes, err := hex.DecodeString(strings.ToLower(s))
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	return _Ed25519PrivateKeyFromBytes(bytes)
}

// PrivateKeyFromKeystore recovers an _Ed25519PrivateKey from an encrypted _Keystore encoded as a byte slice.
func _Ed25519PrivateKeyFromKeystore(ks []byte, passphrase string) (*_Ed25519PrivateKey, error) {
	key, err := _ParseKeystore(ks, passphrase)
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	if key.ed25519PrivateKey != nil {
		return key.ed25519PrivateKey, nil
	}

	return &_Ed25519PrivateKey{}, errors.New("only ed25519 keys are currently supported")
}

// PrivateKeyReadKeystore recovers an _Ed25519PrivateKey from an encrypted _Keystore file.
func _Ed25519PrivateKeyReadKeystore(source io.Reader, passphrase string) (*_Ed25519PrivateKey, error) {
	keystoreBytes, err := io.ReadAll(source)
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	return _Ed25519PrivateKeyFromKeystore(keystoreBytes, passphrase)
}

func _Ed25519PrivateKeyFromPem(bytes []byte, passphrase string) (*_Ed25519PrivateKey, error) {
	var blockType string

	if len(passphrase) == 0 {
		blockType = "PRIVATE KEY"
	} else {
		// the pem is encrypted
		blockType = "ENCRYPTED PRIVATE KEY"
	}

	var pk *pem.Block
	for block, rest := pem.Decode(bytes); block != nil; {
		if block.Type == blockType {
			pk = block
			break
		}

		bytes = rest
		if len(bytes) == 0 {
			// no key was found
			return &_Ed25519PrivateKey{}, _NewErrBadKeyf("pem file did not contain a private key")
		}
	}

	if pk == nil {
		// no key was found
		return &_Ed25519PrivateKey{}, _NewErrBadKeyf("no PEM data is found")
	}

	if len(passphrase) == 0 {
		// key does not need decrypted, end here
		key, err := PrivateKeyFromString(hex.EncodeToString(pk.Bytes))
		if err != nil {
			return &_Ed25519PrivateKey{}, err
		}
		return key.ed25519PrivateKey, nil
	}

	keyI, err := pkcs8.ParsePKCS8PrivateKey(pk.Bytes, []byte(passphrase))
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	return _Ed25519PrivateKeyFromBytes(keyI.(ed25519.PrivateKey))
}

func _Ed25519PrivateKeyReadPem(source io.Reader, passphrase string) (*_Ed25519PrivateKey, error) {
	// note: Passphrases are currently not supported, but included in the function definition to avoid breaking
	// changes in the future.

	pemFileBytes, err := io.ReadAll(source)
	if err != nil {
		return &_Ed25519PrivateKey{}, err
	}

	return _Ed25519PrivateKeyFromPem(pemFileBytes, passphrase)
}

// _Ed25519PublicKey returns the _Ed25519PublicKey associated with this _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _PublicKey() *_Ed25519PublicKey {
	return &_Ed25519PublicKey{
		keyData: sk.keyData[32:],
	}
}

// String returns the text-encoded representation of the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) String() string {
	return sk._StringRaw()
}

// String returns the text-encoded representation of the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _StringDer() string {
	return fmt.Sprint(_Ed25519PrivateKeyPrefix, hex.EncodeToString(sk.keyData[:32]))
}

// String returns the text-encoded representation of the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _StringRaw() string {
	return hex.EncodeToString(sk.keyData[:32])
}

// _BytesRaw returns the byte slice representation of the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _BytesRaw() []byte {
	return sk.keyData[0:32]
}

func (sk _Ed25519PrivateKey) _BytesDer() []byte {
	type PrivateKeyInfo struct {
		Version             int
		PrivateKeyAlgorithm pkix.AlgorithmIdentifier
		PrivateKey          asn1.RawValue
	}

	// AlgorithmIdentifier for Ed25519 keys
	ed25519OID := asn1.ObjectIdentifier{1, 3, 101, 112}
	privateKeyBytes, err := asn1.Marshal(sk.keyData[:32])
	if err != nil {
		return nil
	}
	privateKeyInfo := PrivateKeyInfo{
		Version: 0,
		PrivateKeyAlgorithm: pkix.AlgorithmIdentifier{
			Algorithm: ed25519OID,
		},
		PrivateKey: asn1.RawValue{
			Tag:   asn1.TagOctetString,
			Class: asn1.ClassUniversal,
			Bytes: privateKeyBytes,
		},
	}

	derBytes, err := asn1.Marshal(privateKeyInfo)
	if err != nil {
		return nil
	}

	return derBytes
}

// Keystore returns an encrypted _Keystore containing the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _Keystore(passphrase string) ([]byte, error) {
	return _NewKeystore(sk.keyData, passphrase)
}

// WriteKeystore writes an encrypted _Keystore containing the _Ed25519PrivateKey to the provided destination.
func (sk _Ed25519PrivateKey) _WriteKeystore(destination io.Writer, passphrase string) error {
	keystore, err := sk._Keystore(passphrase)
	if err != nil {
		return err
	}

	_, err = destination.Write(keystore)

	return err
}

// Sign signs the provided message with the _Ed25519PrivateKey.
func (sk _Ed25519PrivateKey) _Sign(message []byte) []byte {
	return ed25519.Sign(sk.keyData, message)
}

// SupportsDerivation returns true if the _Ed25519PrivateKey supports derivation.
func (sk _Ed25519PrivateKey) _SupportsDerivation() bool {
	return sk.chainCode != nil
}

// Derive a child key compatible with the iOS and Android wallets using a provided wallet/account index. Use index 0 for
// the default account.
//
// This will fail if the key does not support derivation which can be checked by calling SupportsDerivation()
func (sk _Ed25519PrivateKey) _Derive(index uint32) (*_Ed25519PrivateKey, error) {
	if !sk._SupportsDerivation() {
		return nil, _NewErrBadKeyf("child key cannot be derived from this key")
	}

	derivedKeyBytes, chainCode, err := _DeriveEd25519ChildKey(sk._BytesRaw(), sk.chainCode, index)
	if err != nil {
		return nil, err
	}

	derivedKey, err := _Ed25519PrivateKeyFromBytes(derivedKeyBytes)

	if err != nil {
		return nil, err
	}

	derivedKey.chainCode = chainCode

	return derivedKey, nil
}

func (sk _Ed25519PrivateKey) _LegacyDerive(index int64) (*_Ed25519PrivateKey, error) {
	keyData, err := _DeriveLegacyChildKey(sk._BytesRaw(), index)
	if err != nil {
		return nil, err
	}

	return _Ed25519PrivateKeyFromBytes(keyData)
}

func (sk _Ed25519PrivateKey) _ToProtoKey() *services.Key {
	return sk._PublicKey()._ToProtoKey()
}

func (sk _Ed25519PrivateKey) _SignTransaction(tx *Transaction) ([]byte, error) {
	tx._RequireOneNodeAccountID()

	if tx.signedTransactions._Length() == 0 {
		return make([]byte, 0), errTransactionRequiresSingleNodeAccountID
	}

	signature := sk._Sign(tx.signedTransactions._Get(0).(*services.SignedTransaction).GetBodyBytes())

	publicKey := sk._PublicKey()
	if publicKey == nil {
		return []byte{}, errors.New("public key is nil")
	}

	wrappedPublicKey := PublicKey{
		ed25519PublicKey: publicKey,
	}

	if tx._KeyAlreadySigned(wrappedPublicKey) {
		return []byte{}, nil
	}

	tx.transactions = _NewLockableSlice()
	tx.publicKeys = append(tx.publicKeys, wrappedPublicKey)
	tx.transactionSigners = append(tx.transactionSigners, nil)
	tx.transactionIDs.locked = true

	for index := 0; index < tx.signedTransactions._Length(); index++ {
		temp := tx.signedTransactions._Get(index).(*services.SignedTransaction)
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		tx.signedTransactions._Set(index, temp)
	}

	return signature, nil
}
