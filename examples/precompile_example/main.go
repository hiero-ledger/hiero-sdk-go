package main

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2/examples/contract_helper"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

//go:embed PrecompileExample.json
var precompileExample []byte

type AbiObject struct {
	ByteCode string `json:"bytecode"`
}

func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
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

	alicePrivateKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating Alice's private key", err))
	}
	alicePublicKey := alicePrivateKey.PublicKey()
	accountCreateResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(alicePublicKey).
		SetInitialBalance(hiero.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating Alice's account", err))
	}

	accountCreateReceipt, err := accountCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account create receipt", err))
	}

	var aliceAccountID hiero.AccountID
	if accountCreateReceipt.AccountID != nil {
		aliceAccountID = *accountCreateReceipt.AccountID
	} else {
		panic("Alice's account id from receipt is nil")
	}

	var abiObject AbiObject
	err = json.Unmarshal(precompileExample, &abiObject)
	if err != nil {
		panic("error reading from json")
	}

	contractFunctionParameters, err := hiero.NewContractFunctionParameters().
		AddAddress(client.GetOperatorAccountID().ToEvmAddress())
	if err != nil {
		panic(fmt.Sprintf("%v : error making contract function parameters", err))
	}

	contractFunctionParameters, err = contractFunctionParameters.
		AddAddress(aliceAccountID.ToEvmAddress())
	if err != nil {
		panic(fmt.Sprintf("%v : error adding alice's address to contract function parameters", err))
	}

	contractHelper := contract_helper.NewContractHelper([]byte(abiObject.ByteCode), *contractFunctionParameters, client)

	keyList := hiero.KeyListWithThreshold(1).Add(operatorKey.PublicKey()).Add(contractHelper.ContractID)

	tx, err := hiero.NewAccountUpdateTransaction().SetAccountID(operatorAccountID).SetKey(keyList).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating alice's account", err))
	}
	_, err = tx.GetReceipt(client)
	if err != nil {
		panic(err)
	}
	keyList = hiero.KeyListWithThreshold(1).Add(alicePublicKey).Add(contractHelper.ContractID)

	frozenTxn, err := hiero.NewAccountUpdateTransaction().SetAccountID(aliceAccountID).SetKey(keyList).FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating alice's account", err))
	}
	tx, err = frozenTxn.Sign(alicePrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating alice's account", err))
	}
	_, err = tx.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	tokenUpdateFunction := func(address string) {
		id, err := hiero.TokenIDFromEvmAddress(0, 0, address)
		if err != nil {
			panic(err)
		}
		frozenTxn, err := hiero.NewTokenUpdateTransaction().SetTokenID(id).SetAdminKey(keyList).SetSupplyKey(keyList).FreezeWith(client)
		if err != nil {
			panic(err)
		}
		tx, err := frozenTxn.Sign(alicePrivateKey).Execute(client)
		if err != nil {
			panic(err)
		}
		_, err = tx.GetReceipt(client)
		if err != nil {
			panic(err)
		}
	}

	contractHelper.
		SetResultValidatorForStep(0, func(contractFunctionResult hiero.ContractFunctionResult) bool {
			println("getPseudoRandomSeed() returned " + hex.EncodeToString(contractFunctionResult.GetBytes32(0)))
			return true
		}).
		SetPayableAmountForStep(1, hiero.NewHbar(20)).
		// step 3 associates Alice with the token, which requires Alice's signature
		AddSignerForStep(3, alicePrivateKey).
		AddSignerForStep(5, alicePrivateKey).
		SetParameterSupplierForStep(11, func() *hiero.ContractFunctionParameters {
			return hiero.NewContractFunctionParameters().
				// when contracts work with a public key, they handle the raw bytes of the public key
				AddBytes(alicePublicKey.BytesRaw())
		}).
		SetPayableAmountForStep(11, hiero.NewHbar(40)).
		// Because we're setting the adminKey for the created NFT token to Alice's key,
		// Alice must sign the ContractExecuteTransaction.
		AddSignerForStep(11, alicePrivateKey).
		SetStepLogic(11, tokenUpdateFunction).
		// and Alice must sign for minting because her key is the supply key.
		AddSignerForStep(12, alicePrivateKey).
		SetParameterSupplierForStep(12, func() *hiero.ContractFunctionParameters {
			return hiero.NewContractFunctionParameters().
				// add three metadatas
				AddBytesArray([][]byte{{0x01b}, {0x02b}, {0x03b}})
		}). // and alice must sign to become associated with the token.
		AddSignerForStep(13, alicePrivateKey).
		// Alice must sign to burn the token because her key is the supply key
		AddSignerForStep(16, alicePrivateKey)

	// step 0 tests pseudo random number generator (PRNG)
	// step 1 creates a fungible token
	// step 2 mints it
	// step 3 associates Alice with it
	// step 4 transfers it to Alice.
	// step 5 approves an allowance of the fungible token with operator as the owner and Alice as the spender
	// steps 6 - 10 test misc functions on the fungible token (see PrecompileExample.sol for details).
	// step 11 creates an NFT token with a custom fee, and with the admin and supply set to Alice's key
	// step 12 mints some NFTs
	// step 13 associates Alice with the NFT token
	// step 14 transfers some NFTs to Alice
	// step 15 approves an NFT allowance with operator as the owner and Alice as the spender
	// step 16 burn some NFTs

	// TODO there is currently possible bug in services causing this operation to fail, should be investigated
	// _, err = contractHelper.
	// 	ExecuteSteps( /* from step */ 0 /* to step */, 16, client)
	// if err != nil {
	// 	panic(fmt.Sprintf("%v : error executing steps", err))
	// }
}
