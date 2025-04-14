package main

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/hiero-ledger/hiero-sdk-go/tck/methods"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
	"github.com/joho/godotenv"
)

func main() {
	// Load dotenv
	_ = godotenv.Load()

	// Initialize the services
	sdkService := new(methods.SDKService)

	accountService := new(methods.AccountService)
	accountService.SetSdkService(sdkService)

	tokenService := new(methods.TokenService)
	tokenService.SetSdkService(sdkService)

	// Create a new RPC server
	assigner := handler.Map{
		"setup":                  postHandler(HandleError, handler.New(sdkService.Setup)),
		"reset":                  postHandler(HandleError, handler.New(sdkService.Reset)),
		"createAccount":          postHandler(HandleError, handler.New(accountService.CreateAccount)),
		"updateAccount":          postHandler(HandleError, handler.New(accountService.UpdateAccount)),
		"deleteAccount":          postHandler(HandleError, handler.New(accountService.DeleteAccount)),
		"approveAllowance":       postHandler(HandleError, handler.New(accountService.ApproveAllowance)),
		"deleteAllowance":        postHandler(HandleError, handler.New(accountService.DeleteAllowance)),
		"transferCrypto":         postHandler(HandleError, handler.New(accountService.TransferCrypto)),
		"createToken":            postHandler(HandleError, handler.New(tokenService.CreateToken)),
		"updateToken":            postHandler(HandleError, handler.New(tokenService.UpdateToken)),
		"deleteToken":            postHandler(HandleError, handler.New(tokenService.DeleteToken)),
		"updateTokenFeeSchedule": postHandler(HandleError, handler.New(tokenService.UpdateTokenFeeSchedule)),
		"associateToken":         postHandler(HandleError, handler.New(tokenService.AssociateToken)),
		"dissociateToken":        postHandler(HandleError, handler.New(tokenService.DissociatesToken)),
		"pauseToken":             postHandler(HandleError, handler.New(tokenService.PauseToken)),
		"unpauseToken":           postHandler(HandleError, handler.New(tokenService.UnpauseToken)),
		"freezeToken":            postHandler(HandleError, handler.New(tokenService.FreezeToken)),
		"unfreezeToken":          postHandler(HandleError, handler.New(tokenService.UnfreezeToken)),
		"grantTokenKyc":          postHandler(HandleError, handler.New(tokenService.GrantTokenKyc)),
		"revokeTokenKyc":         postHandler(HandleError, handler.New(tokenService.RevokeTokenKyc)),
		"mintToken":              postHandler(HandleError, handler.New(tokenService.MintToken)),
		"burnToken":              postHandler(HandleError, handler.New(tokenService.BurnToken)),
		"wipeToken":              postHandler(HandleError, handler.New(tokenService.WipeToken)),
		"generateKey":            postHandler(HandleError, handler.New(methods.GenerateKey)),
	}

	bridge := jhttp.NewBridge(assigner, nil)

	// Listen and redirect to bridge
	http.HandleFunc("/", bridge.ServeHTTP)
	port := os.Getenv("TCK_PORT")
	if port == "" {
		port = "8544"
	}
	log.Println("Server is listening on port: " + port)

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,  // Limit time to read the request
		WriteTimeout: 10 * time.Second, // Limit time to write the response
		IdleTimeout:  60 * time.Second, // Limit idle connections
	}

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for the termination signal
	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err)
	}

	log.Println("Server shutdown complete.")
}

// Handler is a function that handles errors reported by a method handler.
type Handler func(context.Context, *jrpc2.Request, error) error

func HandleError(_ context.Context, request *jrpc2.Request, err error) error {
	if err != nil {
		// jrpc generic error
		if jrpcError, ok := err.(*jrpc2.Error); ok {
			return jrpcError
		}
		// hiero specific errors
		if hieroErr, ok := err.(hiero.ErrHederaReceiptStatus); ok {
			return response.NewHederaReceiptError(hieroErr)
		}
		if hieroErr, ok := err.(hiero.ErrHederaPreCheckStatus); ok {
			return response.NewHederaPrecheckError(hieroErr)
		}
		// other errors
		return response.InternalError
	}
	return nil
}

// this wraps the jrpc2.Handler as it invokes the ErrorHandler func if error is returned
func postHandler(handler Handler, h jrpc2.Handler) jrpc2.Handler {
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		res, err := h(ctx, req)
		if err != nil {
			log.Printf("Error occurred processing JSON-RPC request: %s, Response error: %s", req, err)
			return nil, handler(ctx, req, err)
		}
		return res, nil
	}
}
