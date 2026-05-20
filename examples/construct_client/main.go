package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to construct and configure a client in different ways.
//
// A client has a network and an operator.
//
// A Hedera network is made up of nodes — individual servers who participate
// in the process of reaching consensus on the order and validity of transactions
// on the network. Three networks you likely know of are previewnet, testnet, and mainnet.
//
// For the purpose of connecting to it, each node has an IP address or URL and a port number.
// Each node also has an AccountID used to refer to that node for several purposes,
// including the paying of fees to that node when a client submits requests to it.
//
// You can configure what network you want a client to use — in other words, you can specify
// a list of URLs and port numbers with associated account IDs, and
// when that client is used to execute queries and transactions, the client will
// submit requests only to nodes in that list.
//
// A Client has an operator, which has an AccountID and a PublicKey, and which can
// sign requests. A client's operator can also be configured.
func main() {
	fmt.Println("Construct Client Example Start!")

	// Here's the simplest way to construct a client.
	// These clients' networks are filled with default lists of nodes that are baked into the SDK.
	// Their operators are not yet set, and trying to use them now will result in exceptions.
	testnetClient := hiero.ClientForTestnet()
	previewnetClient := hiero.ClientForPreviewnet()
	mainnetClient := hiero.ClientForMainnet()

	// We can also construct a client for testnet, previewnet or mainnet depending on the value of a
	// network name string. If, for example, the input string equals "testnet", this client will be
	// configured to connect to testnet.
	namedNetworkClient, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client for name", err))
	}

	// Let's set the operator on testnetClient.
	// (The AccountID and PrivateKey here are fake, this is just an example.)
	id, err := hiero.AccountIDFromString("0.0.3")
	if err != nil {
		panic(fmt.Sprintf("%v : error creating AccountID from string", err))
	}
	key, err := hiero.PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	if err != nil {
		panic(fmt.Sprintf("%v : error creating PrivateKey from string", err))
	}
	testnetClient.SetOperator(id, key)

	// Let's create a client with a custom network.
	customNetwork := map[string]hiero.AccountID{
		"2.testnet.hedera.com:50211": {Account: 5},
		"3.testnet.hedera.com:50211": {Account: 6},
	}
	customClient, err := hiero.ClientForNetworkV2(customNetwork)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client for network", err))
	}

	// Since our customClient's network is in this case a subset of testnet, we should set the
	// ledger ID to testnet. If we don't do this, checksum validation won't work.
	customClient.SetLedgerID(*hiero.NewLedgerIDTestnet())

	// Let's generate a client from a config.json file.
	// A config file may specify a network by name, or it may provide a custom network
	// in the form of a list of nodes.
	// The config file should specify the operator, so you can use a client constructed
	// using ClientFromConfigFile() immediately.
	if os.Getenv("CONFIG_FILE") != "" {
		configClient, err := hiero.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))
		if err != nil {
			panic(fmt.Sprintf("%v : error creating Client from config file", err))
		}
		err = configClient.Close()
		if err != nil {
			panic(fmt.Sprintf("%v : error closing configClient", err))
		}
	}

	// Always close a client when you're done with it.
	if err = previewnetClient.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing previewnetClient", err))
	}
	if err = testnetClient.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing testnetClient", err))
	}
	if err = mainnetClient.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing mainnetClient", err))
	}
	if err = namedNetworkClient.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing namedNetworkClient", err))
	}
	if err = customClient.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing customClient", err))
	}

	fmt.Println("Construct Client Example Complete!")
}
