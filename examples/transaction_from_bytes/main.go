package main

import (
	"encoding/hex"
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {

	// Decode a single transactions
	// input := "0a440a150a0c08df9aff9a0610a1c197e502120518cb889c17120218091880c2d72f22020878721e0a1c0a0c0a0518e2f18d17108084af5f0a0c0a0518cb889c1710ff83af5f12660a640a205d3a70d08b2beafb72c7a68986b3ff819a306078b8c359d739e4966e82e6d40e1a40612589c3b15f1e3ed6084b5a3a5b1b81751578cac8d6c922f31731b3982a5bac80a22558b2197276f5bae49b62503a4d39448ceddbc5ef3ba9bee4c0f302f70c"
	input := "0a410a150a0c08b3dd81ca0610e5849a8003120518d3bd83011202182318f5fa0522020878721c0a1a0a0b0a0518e389c2041080897a0a0b0a0518d3bd830110ff887a12660a640a202a428bf2a03fa76bc62452bad6f64c015e543093f3dc77e7b84742ff0f57333f1a40910f5d5410a69258c91905a642c78826a71b371a8c3b88a875f4dec349fd1493a9e4c9806f8ead60683421faa05ecb37c7de35b6b2e678af834a06837b57850d"
	data, err := hex.DecodeString(input)
	transaction, err := hiero.CreateTransferTransactionFromBytes(data)
	if err != nil {
		println(err.Error(), "Error parsing signed transaction")
		fmt.Printf("Transaction %v\n", transaction)
		return
	}
	fmt.Printf("Transaction %v\n", transaction)
}
