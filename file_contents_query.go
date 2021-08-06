package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	Query
	fileID FileID
}

// NewFileContentsQuery creates a FileContentsQuery query which can be used to construct and execute a
// File Get Contents Query.
func NewFileContentsQuery() *FileContentsQuery {
	return &FileContentsQuery{
		Query: newQuery(true),
	}
}

// SetFileID sets the FileID of the file whose contents are requested.
func (query *FileContentsQuery) SetFileID(id FileID) *FileContentsQuery {
	query.fileID = id
	return query
}

func (query *FileContentsQuery) GetFileID(id FileID) FileID {
	return query.fileID
}

func (query *FileContentsQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.fileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *FileContentsQuery) build() *proto.Query_FileGetContents {
	body := &proto.FileGetContentsQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.fileID.isZero() {
		body.FileID = query.fileID.toProtobuf()
	}

	return &proto.Query_FileGetContents{
		FileGetContents: body,
	}
}

func (query *FileContentsQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.FileGetContents.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.FileGetContents.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *FileContentsQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.FileGetContents.Header.Payment = paymentTransaction
	pb.FileGetContents.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *FileContentsQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		fileContentsQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		fileContentsQuery_getMethod,
		fileContentsQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetFileGetContents().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func fileContentsQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetFileGetContents().Header.NodeTransactionPrecheckCode))
}

func fileContentsQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetFileGetContents().Header.NodeTransactionPrecheckCode),
	}
}

func fileContentsQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getFile().GetFileContent,
	}
}

func (query *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return make([]byte, 0), errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
	}

	query.build()

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return []byte{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []byte{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "FileContentsQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []byte{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		fileContentsQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		fileContentsQuery_getMethod,
		fileContentsQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return []byte{}, err
	}

	return resp.query.GetFileGetContents().FileContents.Contents, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *FileContentsQuery) SetNodeAccountIDs(accountID []AccountID) *FileContentsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileContentsQuery) SetMaxRetry(count int) *FileContentsQuery {
	query.Query.SetMaxRetry(count)
	return query
}
