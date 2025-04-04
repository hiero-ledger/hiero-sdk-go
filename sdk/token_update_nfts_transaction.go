package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TokenUpdateNfts struct {
	*Transaction[*TokenUpdateNfts]
	tokenID       *TokenID
	serialNumbers []int64
	metadata      *[]byte
}

func NewTokenUpdateNftsTransaction() *TokenUpdateNfts {
	tx := &TokenUpdateNfts{}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _TokenUpdateNftsTransactionFromProtobuf(tx Transaction[*TokenUpdateNfts], pb *services.TransactionBody) TokenUpdateNfts {
	tokenUpdateNfts := TokenUpdateNfts{
		tokenID:       _TokenIDFromProtobuf(pb.GetTokenUpdateNfts().GetToken()),
		serialNumbers: append([]int64{}, pb.GetTokenUpdateNfts().GetSerialNumbers()...),
	}

	tx.childTransaction = &tokenUpdateNfts
	tokenUpdateNfts.Transaction = &tx
	return tokenUpdateNfts
}

// Getter for tokenID
func (tx *TokenUpdateNfts) GetTokenID() *TokenID {
	return tx.tokenID
}

// Setter for tokenID
func (tx *TokenUpdateNfts) SetTokenID(tokenID TokenID) *TokenUpdateNfts {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// Getter for serialNumbers
func (tx *TokenUpdateNfts) GetSerialNumbers() []int64 {
	return tx.serialNumbers
}

// Setter for serialNumbers
func (tx *TokenUpdateNfts) SetSerialNumbers(serialNumbers []int64) *TokenUpdateNfts {
	tx._RequireNotFrozen()
	tx.serialNumbers = serialNumbers
	return tx
}

// Getter for metadata
func (tx *TokenUpdateNfts) GetMetadata() *[]byte {
	return tx.metadata
}

// Setter for metadata
func (tx *TokenUpdateNfts) SetMetadata(metadata []byte) *TokenUpdateNfts {
	tx._RequireNotFrozen()
	tx.metadata = &metadata
	return tx
}

// ----------- Overridden functions ----------------

func (tx TokenUpdateNfts) getName() string {
	return "TokenUpdateNfts"
}

func (tx TokenUpdateNfts) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}
	return nil
}

func (tx TokenUpdateNfts) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUpdateNfts{
			TokenUpdateNfts: tx.buildProtoBody(),
		},
	}
}

func (tx TokenUpdateNfts) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUpdateNfts{
			TokenUpdateNfts: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenUpdateNfts) buildProtoBody() *services.TokenUpdateNftsTransactionBody {
	body := &services.TokenUpdateNftsTransactionBody{}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}
	serialNumbers := make([]int64, 0)
	if len(tx.serialNumbers) != 0 {
		for _, serialNumber := range tx.serialNumbers {
			serialNumbers = append(serialNumbers, serialNumber)
			body.SerialNumbers = serialNumbers
		}
	}
	if tx.metadata != nil {
		body.Metadata = wrapperspb.Bytes(*tx.metadata)
	}
	return body
}

func (tx TokenUpdateNfts) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateNfts,
	}
}

func (tx TokenUpdateNfts) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenUpdateNfts) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
