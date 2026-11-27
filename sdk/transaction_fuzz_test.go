//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func FuzzTransactionFromBytes(f *testing.F) {
	for _, seed := range transactionFromBytesSeeds(f) {
		f.Add(seed)
	}
	for _, seed := range transactionBodySeeds(f) {
		f.Add(seed)
	}

	// Empty and malformed inputs
	f.Add([]byte{})
	f.Add([]byte{0x0a})
	f.Add([]byte{0x0a, 0x00})
	f.Add([]byte{0x0a, 0x02, 0x2a, 0x00})
	f.Add([]byte{0x2a, 0x02, 0x0a, 0x00})
	f.Add([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	f.Fuzz(func(t *testing.T, data []byte) {
		tx, err := TransactionFromBytes(data)
		if err != nil {
			return
		}

		// What a caller does with a parsed transaction must not panic either.
		// Errors are fine: e.g. an unfrozen transaction has no hash.
		_, _ = TransactionString(tx)
		_, _ = TransactionToBytes(tx)
		_, _ = TransactionGetSignatures(tx)
		_, _ = TransactionGetTransactionHash(tx)
		_, _ = TransactionGetTransactionHashPerNode(tx)
		_, _ = TransactionGetTransactionID(tx)
		_, _ = TransactionGetNodeAccountIDs(tx)
		_, _ = TransactionGetMaxTransactionFee(tx)
		_, _ = TransactionGetTransactionMemo(tx)
		_, _ = TransactionGetTransactionValidDuration(tx)
	})
}

// transactionFromBytesSeeds serializes transactions of different types and in the
// different forms TransactionFromBytes accepts, so the fuzzer starts from inputs that
// reach every branch.
func transactionFromBytesSeeds(f *testing.F) [][]byte {
	f.Helper()

	nodes := []AccountID{{Account: 3}, {Account: 4}}
	txID := testTransactionID

	var seeds [][]byte
	add := func(tx interface{ ToBytes() ([]byte, error) }) {
		f.Helper()
		b, err := tx.ToBytes()
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, b)
	}
	freeze := func(tx TransactionInterface, err error) TransactionInterface {
		f.Helper()
		if err != nil {
			f.Fatal(err)
		}
		return tx
	}
	toBytes := func(tx TransactionInterface) {
		f.Helper()
		b, err := TransactionToBytes(tx)
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, b)
	}

	// Signed with an ED25519 key, one node (utilities_for_test.go)
	mockTx, err := _NewMockTransaction()
	if err != nil {
		f.Fatal(err)
	}
	add(mockTx)

	// Not frozen: the list holds unsigned BodyBytes
	add(NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)))

	// Signed with an ECDSA key, two nodes
	transfer, err := NewTransferTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetTransactionMemo("memo").
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		Freeze()
	if err != nil {
		f.Fatal(err)
	}
	add(transfer.Sign(privateKeyECDSA))

	// Keys and custom fees (batch_transaction_unit_test.go)
	toBytes(freeze(NewTokenCreateTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetTokenName("name").
		SetTokenSymbol("SYM").
		SetAdminKey(privateKeyED25519.PublicKey()).
		SetSupplyKey(privateKeyECDSA.PublicKey()).
		SetCustomFees([]Fee{NewCustomFixedFee().SetAmount(1).SetFeeCollectorAccountID(AccountID{Account: 5})}).
		Freeze()))
	add(spawnTestTransactionAccountCreate())

	// Inner transactions, each parsed by TransactionFromBytes again
	add(spawnTestBatchTransaction())

	// A scheduled transaction body
	scheduled, err := NewScheduleCreateTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetScheduledTransaction(NewTransferTransaction().
			AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)))
	if err != nil {
		f.Fatal(err)
	}
	toBytes(freeze(scheduled.Freeze()))

	// Chunked transactions: one entry per chunk and node
	toBytes(freeze(NewFileAppendTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetFileID(FileID{File: 5}).
		SetMaxChunkSize(16).
		SetContents(bytes.Repeat([]byte{0x01}, 40)).
		Freeze()))
	toBytes(freeze(NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetTopicID(TopicID{Topic: 5}).
		SetChunkSize(16).
		SetMessage(bytes.Repeat([]byte{0x01}, 40)).
		Freeze()))

	// Ethereum and contract transactions
	ethereumData, err := (&EthereumEIP1559Transaction{}).
		SetChainId(298).
		SetGasLimit(21_000).
		SetTo(make([]byte, 20)).
		SetValue(big.NewInt(1)).
		Sign(privateKeyECDSA)
	if err != nil {
		f.Fatal(err)
	}
	toBytes(freeze(NewEthereumTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetEthereumData(ethereumData).
		Freeze()))
	toBytes(freeze(NewContractCreateTransaction().
		SetNodeAccountIDs(nodes).
		SetTransactionID(txID).
		SetBytecode([]byte{0x60, 0x80}).
		SetGas(100_000).
		SetAutoRenewPeriod(time.Hour).
		Freeze()))

	// A single Transaction instead of a TransactionList, in the deprecated form with
	// BodyBytes and SigMap, and in the current form with SignedTransactionBytes
	body, err := protobuf.Marshal(transfer.build())
	if err != nil {
		f.Fatal(err)
	}
	sigMap := &services.SignatureMap{SigPair: []*services.SignaturePair{{
		PubKeyPrefix: privateKeyED25519.PublicKey().BytesRaw(),
		Signature:    &services.SignaturePair_Ed25519{Ed25519: privateKeyED25519.Sign(body)},
	}}}
	signed, err := protobuf.Marshal(&services.SignedTransaction{BodyBytes: body, SigMap: sigMap})
	if err != nil {
		f.Fatal(err)
	}
	for _, msg := range []protobuf.Message{
		&services.Transaction{BodyBytes: body, SigMap: sigMap}, // nolint
		&services.Transaction{SignedTransactionBytes: signed},
		&sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: signed}, {}}},
		&sdk.TransactionList{TransactionList: []*services.Transaction{{}}},
	} {
		b, err := protobuf.Marshal(msg)
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, b)
	}

	return seeds
}

// transactionBodySeeds returns a signed transaction for every transaction type in the
// TransactionBody data oneof, with all fields of the type's body set, so the fuzzer
// starts from inputs that reach every field each type's converter reads. Types added
// to the protobufs later are covered without changing this function.
func transactionBodySeeds(f *testing.F) [][]byte {
	f.Helper()

	var seeds [][]byte
	data := (&services.TransactionBody{}).ProtoReflect().Descriptor().Oneofs().ByName("data")
	for i := 0; i < data.Fields().Len(); i++ {
		field := data.Fields().Get(i)

		body := &services.TransactionBody{
			TransactionID: testTransactionID._ToProtobuf(),
			NodeAccountID: AccountID{Account: 3}._ToProtobuf(),
		}
		value := body.ProtoReflect().NewField(field)
		fillProtoMessage(value.Message(), 4)
		body.ProtoReflect().Set(field, value)

		bodyBytes, err := protobuf.Marshal(body)
		if err != nil {
			f.Fatal(err)
		}
		signed, err := protobuf.Marshal(&services.SignedTransaction{BodyBytes: bodyBytes, SigMap: &services.SignatureMap{}})
		if err != nil {
			f.Fatal(err)
		}
		list, err := protobuf.Marshal(&sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: signed}}})
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, list)
	}

	return seeds
}

// fillProtoMessage sets every field of m, recursing into nested messages until depth
// runs out. Of each oneof only the first field is set, and lists get one element.
func fillProtoMessage(m protoreflect.Message, depth int) {
	if depth == 0 {
		return
	}

	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if oneof := field.ContainingOneof(); oneof != nil && m.WhichOneof(oneof) != nil {
			continue
		}

		switch {
		case field.IsMap():
			continue
		case field.IsList():
			list := m.Mutable(field).List()
			element := list.NewElement()
			if field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind {
				fillProtoMessage(element.Message(), depth-1)
			} else {
				element = protoScalarSeed(field)
			}
			list.Append(element)
		case field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind:
			value := m.NewField(field)
			fillProtoMessage(value.Message(), depth-1)
			m.Set(field, value)
		default:
			m.Set(field, protoScalarSeed(field))
		}
	}
}

func protoScalarSeed(field protoreflect.FieldDescriptor) protoreflect.Value {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(true)
	case protoreflect.EnumKind:
		values := field.Enum().Values()
		return protoreflect.ValueOfEnum(values.Get(values.Len() - 1).Number())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(1)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(1)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(1)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(1)
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(1)
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(1)
	case protoreflect.StringKind:
		return protoreflect.ValueOfString("1")
	default:
		return protoreflect.ValueOfBytes([]byte{1})
	}
}
