package main

import (
	"fmt"
	"net"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// Registered Node Lifecycle
//
// Walks through the full lifecycle of a registered node:
//  1. Generate an admin key.
//  2. Create a registered node with a block node service endpoint.
//  3. Query the RegisteredNodeAddressBook to confirm it appears.
//  4. Update the registered node with a new description and a second endpoint.
//  5. Associate the registered node with an existing consensus node.
//  6. Delete the registered node.
//
// Environment:
//
//	HEDERA_NETWORK          — e.g. "previewnet", "testnet", or "mainnet"
//	OPERATOR_ID             — e.g. "0.0.2"
//	OPERATOR_KEY            — DER-encoded private key
//	CONSENSUS_NODE_ID       — numeric consensus node ID to associate (optional; step 5 is skipped if unset)
//	CONSENSUS_NODE_ADMIN_KEY — DER-encoded admin key for that consensus node (optional; required if CONSENSUS_NODE_ID is set)
func main() {
	client, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Operator: defaults to account 0.0.2 + genesis Ed25519 key (the standard
	// local hedera-services bootstrap, also what the HIP-1137 e2e tests use).
	// Override either via GENESIS_OPERATOR_ID / GENESIS_OPERATOR_KEY if your
	// network has a different system account.
	operatorIDStr := os.Getenv("GENESIS_OPERATOR_ID")
	if operatorIDStr == "" {
		operatorIDStr = "0.0.2"
	}
	operatorAccountID, err := hiero.AccountIDFromString(operatorIDStr)
	if err != nil {
		panic(fmt.Sprintf("%v : error parsing operator ID", err))
	}

	operatorKeyStr := os.Getenv("GENESIS_OPERATOR_KEY")
	if operatorKeyStr == "" {
		operatorKeyStr = "302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137"
	}
	operatorKey, err := hiero.PrivateKeyFromString(operatorKeyStr)
	if err != nil {
		panic(fmt.Sprintf("%v : error parsing operator key", err))
	}

	client.SetOperator(operatorAccountID, operatorKey)

	// Step 1 — generate the registered node admin key.
	adminKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating admin key", err))
	}
	fmt.Printf("Generated registered-node admin key: %s\n", adminKey.PublicKey().String())

	// Step 2 — build the initial block node service endpoint (IP + TLS + two APIs).
	primaryEndpoint := &hiero.BlockNodeServiceEndpoint{}
	primaryEndpoint.
		SetIPAddress(net.IPv4(192, 168, 1, 1).To4()).
		SetPort(50211).
		SetRequiresTls(true).
		SetEndpointApis([]hiero.BlockNodeApi{
			hiero.BlockNodeApiSubscribeStream,
			hiero.BlockNodeApiStatus,
		})

	// Step 3 — submit the create transaction.
	createTx, err := hiero.NewRegisteredNodeCreateTransaction().
		SetAdminKey(adminKey).
		SetDescription("My Block Node").
		SetServiceEndpoints([]hiero.RegisteredServiceEndpoint{primaryEndpoint}).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing create tx", err))
	}

	createResp, err := createTx.Sign(adminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing create tx", err))
	}

	createReceipt, err := createResp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error fetching create receipt", err))
	}

	// Step 4 — verify the assigned registeredNodeId.
	if createReceipt.RegisteredNodeId == nil || *createReceipt.RegisteredNodeId == 0 {
		panic("expected a non-zero registeredNodeId on the receipt")
	}
	registeredNodeId := *createReceipt.RegisteredNodeId
	fmt.Printf("Created registered node with id: %d\n", registeredNodeId)

	// Wait for it
	time.Sleep(time.Second * 5)

	// Step 5 — query the address book and confirm the new node is present.
	book, err := hiero.NewRegisteredNodeAddressBookQuery().
		SetRegisteredNodeId(registeredNodeId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing address book query", err))
	}

	found := false
	for _, n := range book.RegisteredNodes {
		if n.RegisteredNodeID == registeredNodeId {
			found = true
			fmt.Printf("Address book returned registered node: id=%d description=%q\n",
				n.RegisteredNodeID, n.Description)
			break
		}
	}
	if !found {
		panic("registered node was not found in the address book")
	}

	// Step 6 — build the second endpoint (domain name, TLS, STATUS only).
	secondaryEndpoint := &hiero.BlockNodeServiceEndpoint{}
	secondaryEndpoint.
		SetDomainName("block.example.com").
		SetPort(50212).
		SetRequiresTls(true).
		SetEndpointApis([]hiero.BlockNodeApi{hiero.BlockNodeApiStatus})

	// Step 7 — update the registered node with a new description and both endpoints.
	updateTx, err := hiero.NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(registeredNodeId).
		SetDescription("My Updated Block Node").
		SetServiceEndpoints([]hiero.RegisteredServiceEndpoint{primaryEndpoint, secondaryEndpoint}).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing update tx", err))
	}

	updateResp, err := updateTx.Sign(adminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing update tx", err))
	}

	updateReceipt, err := updateResp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error fetching update receipt", err))
	}
	fmt.Printf("Update receipt status: %s\n", updateReceipt.Status)

	// Step 8 — associate with an existing consensus node (optional).
	if consensusNodeIDStr := os.Getenv("CONSENSUS_NODE_ID"); consensusNodeIDStr != "" {
		var consensusNodeID uint64
		if _, err := fmt.Sscanf(consensusNodeIDStr, "%d", &consensusNodeID); err != nil {
			panic(fmt.Sprintf("%v : error parsing CONSENSUS_NODE_ID", err))
		}

		consensusAdminKey, err := hiero.PrivateKeyFromString(os.Getenv("CONSENSUS_NODE_ADMIN_KEY"))
		if err != nil {
			panic(fmt.Sprintf("%v : error parsing CONSENSUS_NODE_ADMIN_KEY", err))
		}

		nodeUpdateTx, err := hiero.NewNodeUpdateTransaction().
			SetNodeID(consensusNodeID).
			AddAssociatedRegisteredNode(registeredNodeId).
			FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error freezing node update tx", err))
		}

		nodeUpdateResp, err := nodeUpdateTx.Sign(consensusAdminKey).Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error executing node update tx", err))
		}

		nodeUpdateReceipt, err := nodeUpdateResp.SetValidateStatus(true).GetReceipt(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error fetching node update receipt", err))
		}
		fmt.Printf("Consensus node %d updated with associated registered node %d: %s\n",
			consensusNodeID, registeredNodeId, nodeUpdateReceipt.Status)
	} else {
		fmt.Println("CONSENSUS_NODE_ID not set — skipping consensus-node association step")
	}

	// Step 9 — delete the registered node.
	deleteTx, err := hiero.NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(registeredNodeId).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing delete tx", err))
	}

	deleteResp, err := deleteTx.Sign(adminKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing delete tx", err))
	}

	deleteReceipt, err := deleteResp.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error fetching delete receipt", err))
	}
	fmt.Printf("Delete receipt status: %s\n", deleteReceipt.Status)

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
}
