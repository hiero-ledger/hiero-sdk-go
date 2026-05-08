package main

import (
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to generate ED25519 keys from BIP-39 mnemonics, plus a legacy 22-word phrase.
func main() {
	fmt.Println("Generate ED25519 Key With Mnemonic Phrase Example Start!")

	// 1. 24-word mnemonic → ED25519 private key → public key
	fmt.Println("Generating random 24-word mnemonic from the BIP-39 standard English word list...")
	mnemonic24, err := hiero.GenerateMnemonic24()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 24 word mnemonic", err))
	}
	fmt.Printf("Generated 24-word mnemonic: %v\n", mnemonic24)

	fmt.Println("Recovering an ED25519 private key from the 24-word mnemonic phrase above...")
	privateKey24, err := mnemonic24.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 24 word mnemonic to PrivateKey", err))
	}
	fmt.Printf("Recovered ED25519 private key: %v\n", privateKey24)

	fmt.Println("Deriving a public key from the above private key...")
	fmt.Printf("Public key: %v\n", privateKey24.PublicKey())

	fmt.Println("---")

	// 2. 12-word mnemonic → ED25519 private key → public key
	fmt.Println("Generating random 12-word mnemonic from the BIP-39 standard English word list...")
	mnemonic12, err := hiero.GenerateMnemonic12()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating 12 word mnemonic", err))
	}
	fmt.Printf("Generated 12-word mnemonic: %v\n", mnemonic12)

	fmt.Println("Recovering an ED25519 private key from the 12-word mnemonic phrase above...")
	privateKey12, err := mnemonic12.ToPrivateKey( /* passphrase */ "")
	if err != nil {
		panic(fmt.Sprintf("%v : error converting 12 word mnemonic to PrivateKey", err))
	}
	fmt.Printf("Recovered ED25519 private key: %v\n", privateKey12)

	fmt.Println("Deriving a public key from the above private key...")
	fmt.Printf("Public key: %v\n", privateKey12.PublicKey())

	fmt.Println("---")

	// 3. Legacy 22-word phrase → legacy private key → public key
	// Note: the legacy mnemonic format does not support a passphrase.
	legacyString := "jolly kidnap tom lawn drunk chick optic lust mutter mole bride galley dense member sage neural widow decide curb aboard margin manure"
	fmt.Printf("Parsing the hardcoded 22-word legacy phrase: %v\n", legacyString)
	mnemonicLegacy, err := hiero.MnemonicFromString(legacyString)
	if err != nil {
		panic(fmt.Sprintf("%v : error parsing legacy mnemonic from string", err))
	}

	fmt.Println("Deriving the legacy private key from the legacy phrase...")
	privateLegacy, err := mnemonicLegacy.ToLegacyPrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error converting legacy mnemonic to PrivateKey", err))
	}
	fmt.Printf("Legacy private key: %v\n", privateLegacy)

	fmt.Println("Deriving a public key from the legacy private key...")
	fmt.Printf("Public key: %v\n", privateLegacy.PublicKey())

	fmt.Println("Generate ED25519 Key With Mnemonic Phrase Example Complete!")
}
