package hedera

import (
	"bytes"
	"crypto/ed25519"
	"github.com/stretchr/testify/assert"
	"strings"

	"testing"
)

const testPrivateKeyStr = "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"

const testPublicKeyStr = "302a300506032b6570032100e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"

// generated by hedera-keygen-java, not used anywhere
const testMnemonic = "inmate flip alley wear offer often piece magnet surge toddler submit right radio absent pear floor belt raven price stove replace reduce plate home"
const testMnemonicKey = "302e020100300506032b657004220420853f15aecd22706b105da1d709b4ac05b4906170c2b9c7495dff9af49e1391da"

// backup phrase generated by the iOS wallet, not used anywhere
const iosMnemonicString = "tiny denial casual grass skull spare awkward indoor ethics dash enough flavor good daughter early hard rug staff capable swallow raise flavor empty angle"

// private key for "default account", should be index 0
const iosDefaultPrivateKey = "5f66a51931e8c99089472e0d70516b6272b94dd772b967f8221e1077f966dbda2b60cf7ee8cf10ecd5a076bffad9a7c7b97df370ad758c0f1dd4ef738e04ceb6"

// backup phrase generated by the Android wallet, also not used anywhere
const androidMnemonicString = "ramp april job flavor surround pyramid fish sea good know blame gate village viable include mixed term draft among monitor swear swing novel track"

// private key for "default account", should be index 0
const androidDefaultPrivateKey = "c284c25b3a1458b59423bc289e83703b125c8eefec4d5aa1b393c2beb9f2bae66188a344ba75c43918ab12fa2ea4a92960eca029a2320d8c6a1c3b94e06c9985"

// test pem key contests for the above testPrivateKeyStr
const pemString = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEINtIS4KOZLLY8SzjwKDpOguMznrxu485yXcyOUSCU44Q
-----END PRIVATE KEY-----
`

// const encryptedPem = `-----BEGIN ENCRYPTED PRIVATE KEY-----
// MIGbMFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAi8WY7Gy2tThQICCAAw
// DAYIKoZIhvcNAgkFADAdBglghkgBZQMEAQIEEOq46NPss58chbjUn20NoK0EQG1x
// R88hIXcWDOECttPTNlMXWJt7Wufm1YwBibrxmCq1QykIyTYhy1TZMyxyPxlYW6aV
// 9hlo4YEh3uEaCmfJzWM=
// -----END ENCRYPTED PRIVATE KEY-----`

const encryptedPem = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIGbMFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAi8WY7Gy2tThQICCAAw
DAYIKoZIhvcNAgkFADAdBglghkgBZQMEAQIEEOq46NPss58chbjUn20NoK0EQG1x
R88hIXcWDOECttPTNlMXWJt7Wufm1YwBibrxmCq1QykIyTYhy1TZMyxyPxlYW6aV
9hlo4YEh3uEaCmfJzWM=
-----END ENCRYPTED PRIVATE KEY-----
`

const pemPassphrase = "this is a passphrase"

func TestPrivateKeyGenerate(t *testing.T) {
	key, err := GeneratePrivateKey()

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(key.String(), PrivateKeyDerPrefix))
}

func TestPrivateKeyExternalSerialization(t *testing.T) {
	key, err := PrivateKeyFromString(testPrivateKeyStr)

	assert.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestPrivateKeyExternalSerializationForConcatenatedHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"
	key, err := PrivateKeyFromString(keyStr)

	assert.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestPrivateKeyExternalSerializationForRawHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
	key, err := PrivateKeyFromString(keyStr)

	assert.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestPublicKeyExternalSerializationForDerEncodedHex(t *testing.T) {
	key, err := PublicKeyFromString(testPublicKeyStr)

	assert.NoError(t, err)
	assert.Equal(t, testPublicKeyStr, key.String())
}

func TestPublicKeyExternalSerializationForRawHex(t *testing.T) {
	keyStr := "e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"
	key, err := PublicKeyFromString(keyStr)

	assert.NoError(t, err)
	assert.Equal(t, testPublicKeyStr, key.String())
}

func TestPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	assert.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")

	assert.NoError(t, err)
	assert.Equal(t, testMnemonicKey, key.String())
}

func TestIOSPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(iosMnemonicString)
	assert.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")
	assert.NoError(t, err)

	derivedKey, err := key.Derive(0)
	assert.NoError(t, err)

	expectedKey, err := PrivateKeyFromString(iosDefaultPrivateKey)
	assert.NoError(t, err)

	assert.Equal(t, expectedKey.keyData, derivedKey.keyData)
}

func TestAndroidPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(androidMnemonicString)
	assert.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")
	assert.NoError(t, err)

	derivedKey, err := key.Derive(0)
	assert.NoError(t, err)

	expectedKey, err := PrivateKeyFromString(androidDefaultPrivateKey)
	assert.NoError(t, err)

	assert.Equal(t, expectedKey.keyData, derivedKey.keyData)
}

func TestSigning(t *testing.T) {
	priKey, err := PrivateKeyFromString(testPrivateKeyStr)
	pubKey, err := PublicKeyFromString(testPublicKeyStr)
	testSignData := []byte("this is the test data to sign")
	signature := priKey.Sign(testSignData)

	assert.NoError(t, err)
	assert.True(t, ed25519.Verify(pubKey.Bytes(), []byte("this is the test data to sign"), signature))
}

func TestGenerated24MnemonicToWorkingPrivateKey(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()

	assert.NoError(t, err)

	privateKey, err := mnemonic.ToPrivateKey("")

	assert.NoError(t, err)

	message := []byte("this is a test message")

	signature := privateKey.Sign(message)

	assert.True(t, ed25519.Verify(privateKey.PublicKey().Bytes(), message, signature))
}

func TestGenerated12MnemonicToWorkingPrivateKey(t *testing.T) {
	mnemonic, err := GenerateMnemonic12()

	assert.NoError(t, err)

	privateKey, err := mnemonic.ToPrivateKey("")

	assert.NoError(t, err)

	message := []byte("this is a test message")

	signature := privateKey.Sign(message)

	assert.True(t, ed25519.Verify(privateKey.PublicKey().Bytes(), message, signature))
}

func TestPrivateKeyFromKeystore(t *testing.T) {
	privatekey, err := PrivateKeyFromKeystore([]byte(testKeystore), passphrase)
	assert.NoError(t, err)

	actualPrivateKey, err := PrivateKeyFromString(testKeystoreKeyString)
	assert.NoError(t, err)

	assert.Equal(t, actualPrivateKey.keyData, privatekey.keyData)
}

func TestPrivateKey_Keystore(t *testing.T) {
	privateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	assert.NoError(t, err)

	keystore, err := privateKey.Keystore(passphrase)
	assert.NoError(t, err)

	ksPrivateKey, err := parseKeystore(keystore, passphrase)
	assert.NoError(t, err)

	assert.Equal(t, privateKey.keyData, ksPrivateKey.keyData)
}

func TestPrivateKey_ReadKeystore(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromString(testKeystoreKeyString)
	assert.NoError(t, err)

	keystoreReader := bytes.NewReader([]byte(testKeystore))

	privateKey, err := PrivateKeyReadKeystore(keystoreReader, passphrase)
	assert.NoError(t, err)

	assert.Equal(t, actualPrivateKey.keyData, privateKey.keyData)
}

func TestPrivateKey_FromPem(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromPem([]byte(pemString), "")
	assert.NoError(t, err)

	assert.Equal(t, actualPrivateKey, privateKey)
}

func TestPrivateKey_FromPemWithPassphrase(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromPem([]byte(encryptedPem), pemPassphrase)
	assert.NoError(t, err)

	assert.Equal(t, actualPrivateKey, privateKey)
}

func TestSetKeyUsesAnyKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		if err != nil {
			panic(err)
		}

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	thresholdKey := KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	_, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetKey(newKey).
		SetKey(thresholdKey).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

}
