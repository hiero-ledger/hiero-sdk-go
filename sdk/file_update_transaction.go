package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// FileUpdateTransaction
// Modify the metadata and/or contents of a file. If a field is not set in the transaction body, the
// corresponding file attribute will be unchanged. This transaction must be signed by all the keys
// in the top level of a key list (M-of-M) of the file being updated. If the keys themselves are
// being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
type FileUpdateTransaction struct {
	*Transaction[*FileUpdateTransaction]
	fileID         *FileID
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
	maxChunks      uint64
	chunkSize      int
}

// NewFileUpdateTransaction creates a FileUpdateTransaction which modifies the metadata and/or contents of a file.
// If a field is not set in the transaction body, the corresponding file attribute will be unchanged.
// tx transaction must be signed by all the keys in the top level of a key list (M-of-M) of the file being updated.
// If the keys themselves are being updated, then the transaction must also be signed by all the new keys. If the keys contain
// additional KeyList or ThresholdKey then M-of-M secondary KeyList or ThresholdKey signing
// requirements must be meet
func NewFileUpdateTransaction() *FileUpdateTransaction {
	tx := &FileUpdateTransaction{
		maxChunks: 20,
		chunkSize: 2048,
	}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _FileUpdateTransactionFromProtobuf(tx Transaction[*FileUpdateTransaction], pb *services.TransactionBody) FileUpdateTransaction {
	var keys *KeyList
	if pb.GetFileUpdate().GetKeys() != nil {
		keysVal, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
		keys = &keysVal
	}
	var expiration *time.Time
	if pb.GetFileUpdate().GetExpirationTime() != nil {
		expirationVal := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())
		expiration = &expirationVal
	}

	fileUpdateTransaction := FileUpdateTransaction{
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           keys,
		expirationTime: expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
		maxChunks:      20,
		chunkSize:      2048,
	}

	tx.childTransaction = &fileUpdateTransaction
	fileUpdateTransaction.Transaction = &tx
	return fileUpdateTransaction
}

// SetFileID Sets the FileID to be updated
func (tx *FileUpdateTransaction) SetFileID(fileID FileID) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID to be updated
func (tx *FileUpdateTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// SetKeys Sets the new list of keys that can modify or delete the file
func (tx *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	if tx.keys == nil {
		tx.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	tx.keys = keyList

	return tx
}

func (tx *FileUpdateTransaction) GetKeys() KeyList {
	if tx.keys != nil {
		return *tx.keys
	}

	return KeyList{}
}

// SetExpirationTime Sets the new expiry time
func (tx *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

// GetExpirationTime returns the new expiry time
func (tx *FileUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// SetContents Sets the new contents that should overwrite the file's current contents
func (tx *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.contents = contents
	return tx
}

// GetContents returns the new contents that should overwrite the file's current contents
func (tx *FileUpdateTransaction) GetContents() []byte {
	return tx.contents
}

// SetMaxChunkSize sets the maximum size of each chunk used by ExecuteAll when the contents are
// larger than a single transaction can hold. Defaults to 2048 bytes, matching FileAppendTransaction.
func (tx *FileUpdateTransaction) SetMaxChunkSize(size int) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.chunkSize = size
	return tx
}

// GetMaxChunkSize returns the maximum size of each chunk used by ExecuteAll.
func (tx *FileUpdateTransaction) GetMaxChunkSize() int {
	return tx.chunkSize
}

// SetMaxChunks sets the maximum number of chunks ExecuteAll is allowed to split the contents into.
// Defaults to 20, matching FileAppendTransaction. ExecuteAll returns ErrMaxChunksExceeded if the
// contents require more chunks than this.
func (tx *FileUpdateTransaction) SetMaxChunks(size uint64) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.maxChunks = size
	return tx
}

// GetMaxChunks returns the maximum number of chunks ExecuteAll is allowed to split the contents into.
func (tx *FileUpdateTransaction) GetMaxChunks() uint64 {
	return tx.maxChunks
}

// SetFileMemo Sets the new memo to be associated with the file (UTF-8 encoding max 100 bytes)
func (tx *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo

	return tx
}

// GeFileMemo
// Deprecated
// use GetFileMemo()
func (tx *FileUpdateTransaction) GeFileMemo() string {
	return tx.memo
}

func (tx *FileUpdateTransaction) GetFileMemo() string {
	return tx.memo
}

// ExecuteAll updates the file, transparently splitting contents larger than a single transaction
// can hold. When the contents fit in one chunk it behaves like Execute and returns a single-element
// slice; otherwise the first chunk overwrites the file via the FileUpdate and the remainder is
// appended with a FileAppendTransaction, so a file larger than the ~6 KiB single-transaction limit
// can be updated in one call. The returned slice holds the FileUpdate response followed by one
// response per append chunk.
//
// Each sub-transaction is charged its own fee and is signed with the operator only. If the file
// requires additional keys, or you need explicit control of the FileID, per-transaction fees, error
// recovery between chunks, or scheduling, use the manual FileUpdateTransaction +
// FileAppendTransaction two-step instead. Execute keeps its single-transaction semantics and still
// returns TRANSACTION_OVERSIZE for oversized contents.
func (tx *FileUpdateTransaction) ExecuteAll(client *Client) ([]TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return nil, errNoClientProvided
	}
	if tx.freezeError != nil {
		return nil, tx.freezeError
	}

	chunkSize := tx.chunkSize
	if chunkSize <= 0 {
		chunkSize = 2048
	}

	chunks := uint64((len(tx.contents) + chunkSize - 1) / chunkSize)
	if chunks == 0 {
		chunks = 1
	}
	if chunks > tx.maxChunks {
		return nil, ErrMaxChunksExceeded{Chunks: chunks, MaxChunks: tx.maxChunks}
	}

	// Fits in a single transaction: behave exactly like Execute.
	if chunks <= 1 {
		resp, err := tx.Execute(client)
		return []TransactionResponse{resp}, err
	}

	// Chunking rebuilds the bodies and appends by FileID, so the transaction must be unfrozen and
	// carry a FileID.
	if tx.IsFrozen() {
		return nil, errFileUpdateChunkingRequiresUnfrozen
	}
	if tx.fileID == nil {
		return nil, errFileUpdateChunkingRequiresFileID
	}

	fullContents := tx.contents
	firstChunkEnd := min(chunkSize, len(fullContents))

	// The first chunk overwrites the file's contents via the update itself; wait for its receipt
	// before appending the rest to preserve ordering.
	tx.SetContents(fullContents[:firstChunkEnd])
	updateResponse, err := tx.Execute(client)
	if err != nil {
		return []TransactionResponse{updateResponse}, err
	}
	if _, err := updateResponse.SetValidateStatus(true).GetReceipt(client); err != nil {
		return []TransactionResponse{updateResponse}, err
	}

	appendTx := NewFileAppendTransaction().
		SetFileID(*tx.fileID).
		SetContents(fullContents[firstChunkEnd:]).
		SetMaxChunkSize(chunkSize).
		SetMaxChunks(tx.maxChunks)
	if nodeAccountIDs := tx.GetNodeAccountIDs(); len(nodeAccountIDs) > 0 {
		appendTx.SetNodeAccountIDs(nodeAccountIDs)
	}

	appendResponses, err := appendTx.ExecuteAll(client)
	responses := append([]TransactionResponse{updateResponse}, appendResponses...)
	return responses, err
}

// Schedule creates a ScheduleCreateTransaction for this FileUpdateTransaction. Chunked contents
// (more than one chunk) cannot be scheduled, mirroring FileAppendTransaction.Schedule.
func (tx *FileUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	chunkSize := tx.chunkSize
	if chunkSize <= 0 {
		chunkSize = 2048
	}
	chunks := uint64((len(tx.contents) + chunkSize - 1) / chunkSize)
	if chunks > 1 {
		return &ScheduleCreateTransaction{}, ErrMaxChunksExceeded{
			Chunks:    chunks,
			MaxChunks: 1,
		}
	}

	return tx.Transaction.Schedule()
}

// ----------- Overridden functions ----------------

func (tx FileUpdateTransaction) getName() string {
	return "FileUpdateTransaction"
}
func (tx FileUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}
func (tx FileUpdateTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_FileUpdate{
		FileUpdate: tx.buildProtoBody(),
	}

	return body
}
func (tx FileUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	body := tx.buildSchedulableTransactionBody()
	body.Data = &services.SchedulableTransactionBody_FileUpdate{
		FileUpdate: tx.buildProtoBody(),
	}

	return body, nil
}
func (tx FileUpdateTransaction) buildProtoBody() *services.FileUpdateTransactionBody {
	body := &services.FileUpdateTransactionBody{
		Memo: &wrapperspb.StringValue{Value: tx.memo},
	}
	if tx.fileID != nil {
		body.FileID = tx.fileID._ToProtobuf()
	}

	if tx.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.keys != nil {
		body.Keys = tx.keys._ToProtoKeyList()
	}

	if tx.contents != nil {
		body.Contents = tx.contents
	}

	return body
}

func (tx FileUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}

func (tx FileUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx FileUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
