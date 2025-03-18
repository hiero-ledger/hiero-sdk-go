package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {
	var client *hiero.Client
	var err error

	/*
	 * Step 1: Initialize the client.
	 * Note: By default, the first network address book update will be executed now
	 * and subsequent updates will occur every 24 hours.
	 * This is controlled by `networkUpdatePeriod`, which defaults to 86400000 milliseconds or 24 hours.
	 */
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	networkUpdateDuration := client.GetNetworkUpdatePeriod()
	fmt.Printf("The current default network update period is: %.2f minutes or %.2f hours.\n", networkUpdateDuration.Minutes(), networkUpdateDuration.Hours())

	/*
	 * Step 2: Change network update period to 1 hour
	 */
	fmt.Println("Changing network update period to 1 hour...")
	client.SetNetworkUpdatePeriod(time.Duration(1) * time.Hour)
	networkUpdateDuration = client.GetNetworkUpdatePeriod()

	fmt.Printf("The current network update period is: %.2f minutes or %.2f hours.\n", networkUpdateDuration.Minutes(), networkUpdateDuration.Hours())

	/*
	 * Step 3: Create client without scheduling network update
	 */
	const testClientJSON string = `{
    "network": {
		"35.237.200.180:50211": "0.0.3",
		"35.186.191.247:50211": "0.0.4",
		"35.192.2.25:50211": "0.0.5",
		"35.199.161.108:50211": "0.0.6",
		"35.203.82.240:50211": "0.0.7",
		"35.236.5.219:50211": "0.0.8",
		"35.197.192.225:50211": "0.0.9",
		"35.242.233.154:50211": "0.0.10",
		"35.240.118.96:50211": "0.0.11",
		"35.204.86.32:50211": "0.0.12"
    },
    "mirrorNetwork": "testnet"
}`
	fmt.Println("Creating client without scheduling network update...")
	client, err = hiero.ClientFromConfigWithoutScheduleNetworkUpdate([]byte(testClientJSON))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	networkUpdateDuration = client.GetNetworkUpdatePeriod()
	fmt.Printf("The current network update period is: %.2f minutes or %.2f hours.\n", networkUpdateDuration.Minutes(), networkUpdateDuration.Hours())

	client.Close()
}
