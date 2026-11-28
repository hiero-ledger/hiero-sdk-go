//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestUnitContractCreateFlowTemporaryFileCleanup(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		bytecode []byte
		appends  [][]byte
	}{
		{name: "small", bytecode: []byte{1, 2, 3}},
		{
			name:     "chunked",
			bytecode: bytes.Repeat([]byte{7}, 6145),
			appends: [][]byte{
				bytes.Repeat([]byte{7}, 2048),
				bytes.Repeat([]byte{7}, 2048),
				{7},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			trace := newContractCreateFlowCleanupTrace()
			steps := trace.creationSteps(test.appends, services.ResponseCodeEnum_OK, services.ResponseCodeEnum_SUCCESS)
			steps = append(steps,
				trace.transaction("delete", services.ResponseCodeEnum_OK),
				trace.receipt("delete receipt", services.ResponseCodeEnum_SUCCESS, nil),
			)
			client, server := NewMockClientAndServer([][]interface{}{{}, steps})
			defer server.Close()
			defer client.Close()
			client.SetMaxAttempts(1)

			response, err := NewContractCreateFlow().
				SetBytecode(test.bytecode).
				SetGas(100000).
				SetNodeAccountIDs([]AccountID{{Account: 4}}).
				Execute(client)
			trace.mu.Lock()
			defer trace.mu.Unlock()
			require.NoError(t, err)

			trace.assertCreation(t, test.appends, test.bytecode, client)
			trace.assertResponse(t, response)
			require.Equal(t, append(trace.creationEvents(test.appends), "delete", "delete receipt"), trace.events)
			require.True(t, proto.Equal(trace.fileID, trace.transactions["delete"].GetFileDelete().GetFileID()))
			trace.assertReceipt(t, "delete receipt", "delete")
		})
	}
}

func TestUnitContractCreateFlowCleanupFailureRetainsContractResponse(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		precheck      services.ResponseCodeEnum
		receiptStatus services.ResponseCodeEnum
	}{
		{
			name:          "submission",
			precheck:      services.ResponseCodeEnum_INVALID_SIGNATURE,
			receiptStatus: services.ResponseCodeEnum_SUCCESS,
		},
		{
			name:          "receipt",
			precheck:      services.ResponseCodeEnum_OK,
			receiptStatus: services.ResponseCodeEnum_INVALID_FILE_ID,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			trace := newContractCreateFlowCleanupTrace()
			steps := trace.creationSteps(nil, services.ResponseCodeEnum_OK, services.ResponseCodeEnum_SUCCESS)
			steps = append(steps, trace.transaction("delete", test.precheck))
			expected := append(trace.creationEvents(nil), "delete")
			if test.precheck == services.ResponseCodeEnum_OK {
				steps = append(steps, trace.receipt("delete receipt", test.receiptStatus, nil))
				expected = append(expected, "delete receipt")
			}
			client, server := NewMockClientAndServer([][]interface{}{{}, steps})
			defer server.Close()
			defer client.Close()
			client.SetMaxAttempts(1)

			response, err := NewContractCreateFlow().
				SetBytecode([]byte{1, 2, 3}).
				SetGas(100000).
				SetNodeAccountIDs([]AccountID{{Account: 4}}).
				Execute(client)
			trace.mu.Lock()
			defer trace.mu.Unlock()
			require.Error(t, err)
			if test.precheck != services.ResponseCodeEnum_OK {
				var precheckErr ErrHederaPreCheckStatus
				require.ErrorAs(t, err, &precheckErr)
				require.Equal(t, StatusInvalidSignature, precheckErr.Status)
			} else {
				var receiptErr ErrHederaReceiptStatus
				require.ErrorAs(t, err, &receiptErr)
				require.Equal(t, StatusInvalidFileID, receiptErr.Status)
				trace.assertReceipt(t, "delete receipt", "delete")
			}
			trace.assertCreation(t, nil, []byte{1, 2, 3}, client)
			trace.assertResponse(t, response)
			require.Equal(t, expected, trace.events)
			require.True(t, proto.Equal(trace.fileID, trace.transactions["delete"].GetFileDelete().GetFileID()))
		})
	}
}

func TestUnitContractCreateFlowDoesNotDeleteFileAfterContractFailure(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		precheck      services.ResponseCodeEnum
		receiptStatus services.ResponseCodeEnum
	}{
		{"submission", services.ResponseCodeEnum_INVALID_SIGNATURE, services.ResponseCodeEnum_SUCCESS},
		{"receipt", services.ResponseCodeEnum_OK, services.ResponseCodeEnum_CONTRACT_REVERT_EXECUTED},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			trace := newContractCreateFlowCleanupTrace()
			steps := trace.creationSteps(nil, test.precheck, test.receiptStatus)
			steps = append(steps,
				trace.transaction("unexpected delete", services.ResponseCodeEnum_OK),
				trace.receipt("unexpected delete receipt", services.ResponseCodeEnum_SUCCESS, nil),
			)
			client, server := NewMockClientAndServer([][]interface{}{{}, steps})
			defer server.Close()
			defer client.Close()
			client.SetMaxAttempts(1)

			response, err := NewContractCreateFlow().
				SetBytecode([]byte{1, 2, 3}).
				SetGas(100000).
				SetNodeAccountIDs([]AccountID{{Account: 4}}).
				Execute(client)
			trace.mu.Lock()
			defer trace.mu.Unlock()
			require.Error(t, err)
			require.Equal(t, TransactionResponse{}, response)
			if test.precheck != services.ResponseCodeEnum_OK {
				var precheckErr ErrHederaPreCheckStatus
				require.ErrorAs(t, err, &precheckErr)
				require.Equal(t, StatusInvalidSignature, precheckErr.Status)
				require.Equal(t, []string{"file create", "file receipt", "contract create"}, trace.events)
			} else {
				var receiptErr ErrHederaReceiptStatus
				require.ErrorAs(t, err, &receiptErr)
				require.Equal(t, StatusContractRevertExecuted, receiptErr.Status)
				trace.assertCreation(t, nil, []byte{1, 2, 3}, client)
				require.Equal(t, trace.creationEvents(nil), trace.events)
			}
		})
	}
}

// The mock records wire requests rather than replacing the SDK transaction builders.
// A mutex also makes the observations safe when the tests run in parallel.
type contractCreateFlowCleanupTrace struct {
	mu           sync.Mutex
	fileID       *services.FileID
	events       []string
	transactions map[string]*services.TransactionBody
	signed       map[string]*services.SignedTransaction
	queries      map[string]*services.Query
	decodeErrors []error
}

func newContractCreateFlowCleanupTrace() *contractCreateFlowCleanupTrace {
	return &contractCreateFlowCleanupTrace{
		fileID:       &services.FileID{ShardNum: 1, RealmNum: 2, FileNum: 3456},
		transactions: make(map[string]*services.TransactionBody),
		signed:       make(map[string]*services.SignedTransaction),
		queries:      make(map[string]*services.Query),
	}
}

func (trace *contractCreateFlowCleanupTrace) transaction(name string, status services.ResponseCodeEnum) interface{} {
	return func(request *services.Transaction) *services.TransactionResponse {
		trace.mu.Lock()
		defer trace.mu.Unlock()
		trace.events = append(trace.events, name)
		signed := &services.SignedTransaction{}
		body := &services.TransactionBody{}
		if err := proto.Unmarshal(request.GetSignedTransactionBytes(), signed); err != nil {
			trace.decodeErrors = append(trace.decodeErrors, err)
		}
		if err := proto.Unmarshal(signed.GetBodyBytes(), body); err != nil {
			trace.decodeErrors = append(trace.decodeErrors, err)
		}
		trace.transactions[name] = body
		trace.signed[name] = signed
		return &services.TransactionResponse{NodeTransactionPrecheckCode: status}
	}
}

func (trace *contractCreateFlowCleanupTrace) receipt(name string, status services.ResponseCodeEnum, fileID *services.FileID) interface{} {
	return func(request *services.Query) *services.Response {
		trace.mu.Lock()
		defer trace.mu.Unlock()
		trace.events = append(trace.events, name)
		trace.queries[name] = request
		return &services.Response{Response: &services.Response_TransactionGetReceipt{
			TransactionGetReceipt: &services.TransactionGetReceiptResponse{
				Header: &services.ResponseHeader{
					NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
					ResponseType:                services.ResponseType_ANSWER_ONLY,
				},
				Receipt: &services.TransactionReceipt{
					Status: status, FileID: fileID,
					ContractID: &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 7890}},
				},
			},
		}}
	}
}

func (trace *contractCreateFlowCleanupTrace) creationSteps(appends [][]byte, precheck, receiptStatus services.ResponseCodeEnum) []interface{} {
	steps := []interface{}{
		trace.transaction("file create", services.ResponseCodeEnum_OK),
		trace.receipt("file receipt", services.ResponseCodeEnum_SUCCESS, trace.fileID),
	}
	for i := range appends {
		name := fmt.Sprintf("append %d", i+1)
		steps = append(steps, trace.transaction(name, services.ResponseCodeEnum_OK), trace.receipt(name+" receipt", services.ResponseCodeEnum_SUCCESS, nil))
	}
	if len(appends) > 0 {
		// ExecuteAll gets each chunk's receipt, then the flow gets the first again.
		steps = append(steps, trace.receipt("flow append receipt", services.ResponseCodeEnum_SUCCESS, nil))
	}
	steps = append(steps, trace.transaction("contract create", precheck))
	if precheck == services.ResponseCodeEnum_OK {
		steps = append(steps, trace.receipt("contract receipt", receiptStatus, nil))
	}
	return steps
}

func (trace *contractCreateFlowCleanupTrace) creationEvents(appends [][]byte) []string {
	events := []string{"file create", "file receipt"}
	for i := range appends {
		name := fmt.Sprintf("append %d", i+1)
		events = append(events, name, name+" receipt")
	}
	if len(appends) > 0 {
		events = append(events, "flow append receipt")
	}
	return append(events, "contract create", "contract receipt")
}

func (trace *contractCreateFlowCleanupTrace) assertReceipt(t *testing.T, receipt, transaction string) {
	t.Helper()
	require.True(t, proto.Equal(trace.transactions[transaction].GetTransactionID(), trace.queries[receipt].GetTransactionGetReceipt().GetTransactionID()), receipt)
}

func (trace *contractCreateFlowCleanupTrace) assertCreation(t *testing.T, appends [][]byte, bytecode []byte, client *Client) {
	t.Helper()
	require.Empty(t, trace.decodeErrors)
	trace.assertReceipt(t, "file receipt", "file create")
	trace.assertReceipt(t, "contract receipt", "contract create")
	require.True(t, proto.Equal(trace.fileID, trace.transactions["contract create"].GetContractCreateInstance().GetFileID()))
	contents := append([]byte(nil), trace.transactions["file create"].GetFileCreate().GetContents()...)
	for i, expected := range appends {
		name := fmt.Sprintf("append %d", i+1)
		body := trace.transactions[name].GetFileAppend()
		require.True(t, proto.Equal(trace.fileID, body.GetFileID()))
		require.Equal(t, expected, body.GetContents())
		contents = append(contents, body.GetContents()...)
		trace.assertReceipt(t, name+" receipt", name)
	}
	if len(appends) > 0 {
		trace.assertReceipt(t, "flow append receipt", "append 1")
	}
	require.Equal(t, bytecode, contents)
	operatorAccountID := client.GetOperatorAccountID()
	operatorPublicKey := client.GetOperatorPublicKey()
	for name, body := range trace.transactions {
		require.Equal(t, int64(4), body.GetNodeAccountID().GetAccountNum(), name)
		require.True(t, proto.Equal(operatorAccountID._ToProtobuf(), body.GetTransactionID().GetAccountID()), name)
		pairs := trace.signed[name].GetSigMap().GetSigPair()
		require.Len(t, pairs, 1, name)
		require.True(t, operatorPublicKey.VerifySignedMessage(trace.signed[name].GetBodyBytes(), pairs[0].GetEd25519()), name)
	}
}

func (trace *contractCreateFlowCleanupTrace) assertResponse(t *testing.T, response TransactionResponse) {
	t.Helper()
	require.True(t, proto.Equal(trace.transactions["contract create"].GetTransactionID(), response.TransactionID._ToProtobuf()))
	require.Equal(t, AccountID{Account: 4}, response.NodeID)
	require.NotEmpty(t, response.Hash)
}
