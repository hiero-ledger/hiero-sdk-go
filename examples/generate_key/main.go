package main

import (
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to generate an ECDSA (secp256k1) key pair and derive the EVM address from the public key.
//
// ECDSA keys with an EVM address derived from the public key are the recommended choice for new
// Hedera accounts when you want compatibility with Ethereum tooling (wallets, Hardhat, ethers.js,
// etc.).
func main() {
	fmt.Println("Generate ECDSA key pair and EVM address example start")

	fmt.Println("Generating an ECDSA (secp256k1) private key...")
	privateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating ECDSA PrivateKey", err))
	}
	fmt.Printf("Private key: %v\n", privateKey)

	fmt.Println("Deriving the public key from the private key")
	publicKey := privateKey.PublicKey()
	fmt.Printf("Public key: %v\n", publicKey)

	fmt.Println("Deriving the EVM address (last 20 bytes of Keccak-256 of the uncompressed public key)")
	evmAddress := publicKey.ToEvmAddress()
	fmt.Printf("EVM address: 0x%s\n", evmAddress)

	fmt.Println("Generate ECDSA key pair and EVM address example complete")
}
