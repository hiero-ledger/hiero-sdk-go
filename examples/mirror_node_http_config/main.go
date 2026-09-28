package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// loggingTransport wraps another transport and prints every request.
type loggingTransport struct {
	next hiero.HttpTransport
}

func (t loggingTransport) RoundTrip(req hiero.HttpRequest, cancellation hiero.Cancellation) (hiero.HttpResponse, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req, cancellation)
	if err != nil {
		fmt.Printf("  %s %s failed after %v: %v\n", req.GetMethod(), req.GetURL(), time.Since(start).Round(time.Millisecond), err)
		return resp, err
	}
	fmt.Printf("  %s %s -> %d in %v\n", req.GetMethod(), req.GetURL(), resp.GetStatusCode(), time.Since(start).Round(time.Millisecond))
	return resp, nil
}

func (t loggingTransport) Close(closeTimeout time.Duration) {
	t.next.Close(closeTimeout)
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
	defer client.Close()

	// SetMirrorNodeHttpConfig replaces the whole configuration, so start from the current one.
	config := client.GetMirrorNodeHttpConfig()

	// Step 1: Tighten timeouts and retries. Each With* returns a copy with one field changed.
	client.SetMirrorNodeHttpConfig(config.
		WithTransportConfiguration(config.GetTransportConfiguration().WithConnectTimeout(2*time.Second)).
		WithRetryPolicy(config.GetRetryPolicy().
			WithMaxAttempts(3).
			WithPerAttemptTimeout(5*time.Second).
			WithTotalDeadline(15*time.Second)).
		WithRequestHeader("x-example-request", "mirror-node-http-config"))

	balance, err := hiero.NewMirrorNodeAccountBalanceQuery().SetAccountID(operatorAccountID).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing mirror node account balance query", err))
	}
	fmt.Printf("balance with tighter bounds = %v\n", balance.Hbars)

	// Step 2: Supply the transport. The application owns it, so Client.Close leaves it open.
	transport := loggingTransport{next: hiero.NewDefaultHttpTransport(hiero.DefaultHttpTransportConfiguration())}
	defer transport.Close(5 * time.Second)

	client.SetMirrorNodeHttpConfig(client.GetMirrorNodeHttpConfig().WithTransport(transport))

	fmt.Println("balance through the application's transport:")
	balance, err = hiero.NewMirrorNodeAccountBalanceQuery().SetAccountID(operatorAccountID).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing mirror node account balance query", err))
	}
	fmt.Printf("balance = %v\n", balance.Hbars)

	// GetMirrorNodeHttpConfig returns the configuration as set.
	fmt.Printf("max attempts = %d, total deadline = %v\n",
		client.GetMirrorNodeHttpConfig().GetRetryPolicy().GetMaxAttempts(),
		client.GetMirrorNodeHttpConfig().GetRetryPolicy().GetTotalDeadline())
}
