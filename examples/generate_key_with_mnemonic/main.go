package main

import (
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to generate ECDSA keys from BIP-39 mnemonics.
func main() {
	fmt.Println("Generate ECDSA Key With Mnemonic Phrase Example Start!")

	// 1. 24-word mnemonic → ECDSA private key → public key
	fmt.Println("Generating random 24-word mnemonic from the BIP-39 standard English word list...")
	mnemonic24, err := hiero.GenerateMnemonic24()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 24 word mnemonic", err))
	}
	fmt.Printf("Generated 24-word mnemonic: %v\n", mnemonic24)

	fmt.Println("Recovering an ECDSA private key from the 24-word mnemonic phrase above...")
	privateKey24, err := mnemonic24.ToStandardECDSAsecp256k1PrivateKey( /* passphrase */ "", /* index */ 0)
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 24 word mnemonic to PrivateKey", err))
	}
	fmt.Printf("Recovered ECDSA private key: %v\n", privateKey24)

	fmt.Println("Deriving a public key from the above private key...")
	fmt.Printf("Public key: %v\n", privateKey24.PublicKey())

	fmt.Println("---")

	// 2. 12-word mnemonic → ECDSA private key → public key
	fmt.Println("Generating random 12-word mnemonic from the BIP-39 standard English word list...")
	mnemonic12, err := hiero.GenerateMnemonic12()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 12 word mnemonic", err))
	}
	fmt.Printf("Generated 12-word mnemonic: %v\n", mnemonic12)

	fmt.Println("Recovering an ECDSA private key from the 12-word mnemonic phrase above...")
	privateKey12, err := mnemonic12.ToStandardECDSAsecp256k1PrivateKey( /* passphrase */ "", /* index */ 0)
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 12 word mnemonic to PrivateKey", err))
	}
	fmt.Printf("Recovered ECDSA private key: %v\n", privateKey12)

	fmt.Println("Deriving a public key from the above private key...")
	fmt.Printf("Public key: %v\n", privateKey12.PublicKey())

	fmt.Println("Generate ECDSA Key With Mnemonic Phrase Example Complete!")
}
