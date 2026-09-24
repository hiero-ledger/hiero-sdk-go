package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {
	var client *hiero.Client
	var err error

	// Initialize the client with the mirror node of the network named by HEDERA_NETWORK. This will also get the
	// address book from the mirror node and use it to populate the Client's consensus network.
	network := os.Getenv("HEDERA_NETWORK")
	mirrorNodes := map[string]string{
		"mainnet":    "mainnet-public.mirrornode.hedera.com:443",
		"testnet":    "testnet.mirrornode.hedera.com:443",
		"previewnet": "previewnet.mirrornode.hedera.com:443",
		"localhost":  "127.0.0.1:5600",
	}
	mirrorNode, ok := mirrorNodes[network]
	if !ok {
		panic(fmt.Sprintf("HEDERA_NETWORK %q has no known mirror node", network))
	}

	client, err = hiero.ClientForMirrorNetwork([]string{mirrorNode})
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}
	fmt.Printf("Consensus nodes from the %s address book: %v\n", network, client.GetNetwork())

	// A local network advertises the addresses its nodes have inside the cluster, which this machine cannot reach, so
	// talk to the node through its local port instead.
	if network == "localhost" {
		if err := client.SetNetwork(map[string]hiero.AccountID{"127.0.0.1:50211": {Account: 3}}); err != nil {
			panic(fmt.Sprintf("%v : error setting the local network", err))
		}
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	privateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey()

	txResponse, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(1)).
		SetKeyWithoutAlias(publicKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create transaction", err))
	}

	receipt, err := txResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	fmt.Printf("New account id, %s", receipt.AccountID.String())
}
