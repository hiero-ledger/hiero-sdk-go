package main

import (
	"encoding/hex"
	"fmt"

	"github.com/Arculus-Holdings-L-L-C/hedera-sdk-go"
)

func main() {

	// Decode a single transactions
	input := "0a440a150a0c08df9aff9a0610a1c197e502120518cb889c17120218091880c2d72f22020878721e0a1c0a0c0a0518e2f18d17108084af5f0a0c0a0518cb889c1710ff83af5f12660a640a205d3a70d08b2beafb72c7a68986b3ff819a306078b8c359d739e4966e82e6d40e1a40612589c3b15f1e3ed6084b5a3a5b1b81751578cac8d6c922f31731b3982a5bac80a22558b2197276f5bae49b62503a4d39448ceddbc5ef3ba9bee4c0f302f70c"
	data, err := hex.DecodeString(input)
	transaction, err := hedera.CreateTransferTransactionFromBytes(data)
	if err != nil {
		println(err.Error(), "Error parsing signed transaction")
		fmt.Printf("Transaction %v\n", transaction)
		return
	}
	fmt.Printf("Transaction %v\n", transaction)
}
