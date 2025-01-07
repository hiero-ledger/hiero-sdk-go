package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {

	operatorID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		fmt.Println("Error parsing OPERATOR_ID")
		return
	}

	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		fmt.Println("Error parsing OPERATOR_KEY")
		return
	}

	// Step 0: Create and configure the SDK Client
	node8 := map[string]hiero.AccountID{
		"35.155.49.147:50211": hiero.AccountID{Account: 8},
	}
	client := hiero.ClientForNetwork(node8)
	client.SetOperator(operatorID, operatorKey)

	trimmedNodes := map[string]hiero.AccountID{
		"34.94.106.61:50211":   hiero.AccountID{Account: 3},
		"3.212.6.13:50211":     hiero.AccountID{Account: 4},
		"35.245.27.193:50211":  hiero.AccountID{Account: 5},
		"34.83.112.116:50211":  hiero.AccountID{Account: 6},
		"34.94.160.4:50211":    hiero.AccountID{Account: 7},
		"34.133.197.230:50211": hiero.AccountID{Account: 9},
	}

	clientTwo := hiero.ClientForNetwork(trimmedNodes)
	clientTwo.SetOperator(operatorID, operatorKey)

	for {
		time.Sleep(2 * time.Second)

		transfer, err := hiero.NewTransferTransaction().
			AddHbarTransfer(operatorID, hiero.HbarFrom(1, hiero.HbarUnits.Hbar).Negated()).
			AddHbarTransfer(hiero.AccountID{Account: 5172898}, hiero.HbarFrom(1, hiero.HbarUnits.Hbar)).
			FreezeWith(client)
		if err != nil {
			fmt.Println("Error freezing transaction:", err)
			continue
		}

		transfer.Sign(operatorKey)
		fmt.Println(transfer.GetNodeAccountIDs())

		resp, err := transfer.Execute(clientTwo)
		if err != nil {
			fmt.Println("Error executing transaction:", err)
			continue
		}

		_, err = resp.GetReceipt(clientTwo)
		if err != nil {
			fmt.Println("Error getting receipt:", err)
			continue
		}

		fmt.Println("Transfer executed against node", resp.NodeID)
	}
}
