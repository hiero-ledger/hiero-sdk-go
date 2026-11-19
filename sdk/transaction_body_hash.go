package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/sha512"
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/protobuf/encoding/protowire"
	protobuf "google.golang.org/protobuf/proto"
)

// GetTransactionBodyHash returns a 48-byte SHA-384 commitment to every prepared
// node and chunk of a frozen transaction, excluding outer and batch-child
// signature maps and other SignedTransaction envelope metadata. It uses original
// body bytes without re-encoding them. Independently encoded protobuf messages
// with equivalent content need not have the same hash.
//
// This implements the candidate profile in sdk-collaboration-hub PR #288.
// Framing, import ordering and rejection rules require maintainer agreement
// before SDK PRs merge; this is not an approved cross-SDK contract.
//
// The getter performs no signing or network calls and preserves iteration
// positions, bodies and signatures. Success locks transaction IDs and node
// membership against execution regeneration or replacement, while allowing
// subsequent signing. Failure leaves the transaction unchanged.
// Selected ID/node membership must match the complete prepared layout.
// Pending common-field edits that differ from prepared bodies are rejected.
// Freeze prepares configured common fields without signing or serialization.
// Unchanged editable values retain imported nested unknown bytes. Raw wire
// bounds apply to bodies, identifiers, batch wrappers and child envelopes;
// operation-specific submessages retain existing SDK decoding checks.
// After success, body-changing common setters report the existing frozen error,
// and builders reject any indirect edit that would rewrite a reviewed body.
func (tx *Transaction[T]) GetTransactionBodyHash() ([]byte, error) {
	if !tx.IsFrozen() {
		return nil, errTransactionIsNotFrozen
	}

	// Keep grouping local: storage order and the build/signing cursor are unrelated
	// to the canonical chunk-major traversal of this commitment.
	chunks := make([]map[string][sha512.Size384]byte, 0)
	chunkIndices := make(map[_TransactionIDValue]int)
	var nodes []string
	bodies := make([]*services.TransactionBody, 0, tx.signedTransactions._Length())
	for _, value := range tx.signedTransactions.slice {
		signed, ok := value.(*services.SignedTransaction)
		if !ok || signed == nil {
			return nil, fmt.Errorf("invalid prepared signed transaction")
		}
		entry, body, err := _transactionBodyHashEntry(signed.BodyBytes, false)
		if err != nil {
			return nil, err
		}
		if tx._hasPendingCommonBodyChanges(body) {
			return nil, fmt.Errorf("pending common-field changes differ from prepared transaction bodies")
		}
		bodies = append(bodies, body)
		id := _transactionIDValue(body.TransactionID)
		node := _accountIDValueKey(body.NodeAccountID)
		chunkIndex, exists := chunkIndices[id]
		if !exists {
			chunkIndex = len(chunks)
			chunkIndices[id] = chunkIndex
			chunks = append(chunks, make(map[string][sha512.Size384]byte))
		}
		chunk := chunks[chunkIndex]
		if _, duplicate := chunk[node]; duplicate {
			return nil, fmt.Errorf("duplicate prepared transaction ID and node ID")
		}
		chunk[node] = entry
		if chunkIndex == 0 {
			nodes = append(nodes, node)
		}
	}
	if tx.transactionIDs._Length() != len(chunkIndices) || tx.nodeAccountIDs._Length() != len(nodes) {
		return nil, fmt.Errorf("selected transaction IDs or node IDs differ from prepared membership")
	}
	selectedIDs := make(map[_TransactionIDValue]bool, len(chunkIndices))
	for _, value := range tx.transactionIDs.slice {
		id := _transactionIDValue(value.(TransactionID)._ToProtobuf())
		_, prepared := chunkIndices[id]
		if selectedIDs[id] || !prepared {
			return nil, fmt.Errorf("duplicate or unprepared selected transaction ID")
		}
		selectedIDs[id] = true
	}
	selectedNodes := make(map[string]bool, len(nodes))
	for _, value := range tx.nodeAccountIDs.slice {
		node := _accountIDValueKey(value.(AccountID)._ToProtobuf())
		_, prepared := chunks[0][node]
		if selectedNodes[node] || !prepared {
			return nil, fmt.Errorf("duplicate or unprepared selected node ID")
		}
		selectedNodes[node] = true
	}
	valid, err := _transactionBodiesEqual(bodies)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("failed to validate transaction bodies")
	}
	aggregate := sha512.New384()
	for _, chunk := range chunks {
		if len(chunk) != len(nodes) {
			return nil, fmt.Errorf("inconsistent prepared node sets")
		}
		for _, node := range nodes {
			entry, exists := chunk[node]
			if !exists {
				return nil, fmt.Errorf("inconsistent prepared node sets")
			}
			_, _ = aggregate.Write(entry[:]) // hash.Hash.Write never returns an error.
		}
	}
	tx.transactionIDs.locked = true
	tx.nodeAccountIDs.locked = true
	tx.bodyReviewed = true
	return aggregate.Sum(nil), nil
}

// TransactionGetTransactionBodyHash is GetTransactionBodyHash for a
// TransactionInterface, including transactions returned by TransactionFromBytes.
// It follows the same pending candidate profile in sdk-collaboration-hub PR #288.
func TransactionGetTransactionBodyHash(tx TransactionInterface) ([]byte, error) {
	return tx.getBaseTransaction().GetTransactionBodyHash()
}

// _transactionBodyHashEntry validates raw records before decoding can discard
// repeated oneof fields, and only hashes original slices of the input bytes.
func _transactionBodyHashEntry(raw []byte, child bool) ([sha512.Size384]byte, *services.TransactionBody, error) {
	var zero [sha512.Size384]byte
	body, batch, outer, err := _transactionBodyHashBody(raw, child)
	if err != nil {
		return zero, nil, err
	}
	if batch == nil {
		return sha512.Sum384(raw), body, nil
	}
	parts := sha512.New384()
	outerHash := sha512.Sum384(outer)
	_, _ = parts.Write(outerHash[:])
	children := 0
	for len(batch) > 0 {
		field, err := _bodyHashConsumeField(batch)
		if err != nil {
			return zero, nil, err
		}
		if field.number != 1 || field.wireType != protowire.BytesType {
			return zero, nil, fmt.Errorf("unsupported atomic batch wrapper field")
		}
		childBody, err := _transactionBodyHashChild(field.bytes)
		if err != nil {
			return zero, nil, err
		}
		childHash, _, err := _transactionBodyHashEntry(childBody, true)
		if err != nil {
			return zero, nil, err
		}
		_, _ = parts.Write(childHash[:])
		children++
		batch = batch[field.size:]
	}
	if children == 0 {
		return zero, nil, fmt.Errorf("empty atomic batch")
	}
	copy(zero[:], parts.Sum(nil))
	return zero, body, nil
}

func _transactionBodyHashBody(raw []byte, child bool) (*services.TransactionBody, []byte, []byte, error) {
	var body services.TransactionBody
	descriptor := body.ProtoReflect().Descriptor()
	operation := descriptor.Oneofs().ByName("data")
	operations := 0
	var batch, outer []byte
	for remaining := raw; len(remaining) > 0; {
		field, err := _bodyHashConsumeField(remaining)
		if err != nil {
			return nil, nil, nil, err
		}
		if field.number == 1 || field.number == 2 {
			if field.wireType != protowire.BytesType {
				return nil, nil, nil, fmt.Errorf("invalid transaction identifier wire type")
			}
			if err := _bodyHashValidateIdentifier(field.bytes, field.number == 1); err != nil {
				return nil, nil, nil, err
			}
		}
		known := descriptor.Fields().ByNumber(field.number)
		if known != nil && known.ContainingOneof() == operation {
			operations++
			if operations > 1 || field.wireType != protowire.BytesType {
				return nil, nil, nil, fmt.Errorf("repeated, conflicting or mistyped transaction operation")
			}
		}
		if field.number == 74 {
			if child {
				return nil, nil, nil, fmt.Errorf("nested atomic batch")
			}
			batch = field.bytes
			// Removing the entire wire record preserves every other byte and its order.
			start := len(raw) - len(remaining)
			outer = append(append([]byte(nil), raw[:start]...), remaining[field.size:]...)
		}
		remaining = remaining[field.size:]
	}
	if operations != 1 {
		return nil, nil, nil, fmt.Errorf("missing transaction operation")
	}
	if err := protobuf.Unmarshal(raw, &body); err != nil {
		return nil, nil, nil, err
	}
	if body.TransactionID == nil || body.TransactionID.AccountID == nil || body.TransactionID.TransactionValidStart == nil || body.NodeAccountID == nil {
		return nil, nil, nil, fmt.Errorf("missing prepared transaction ID or node ID")
	}
	return &body, batch, outer, nil
}

// Only TransactionID has inspected child messages: valid start (1) and account
// (2). Unknown length-delimited payloads, including aliases, remain opaque.
func _bodyHashValidateIdentifier(raw []byte, transactionID bool) error {
	for len(raw) > 0 {
		field, err := _bodyHashConsumeField(raw)
		if err != nil {
			return err
		}
		if transactionID && (field.number == 1 || field.number == 2) {
			if field.wireType != protowire.BytesType {
				return fmt.Errorf("invalid transaction ID child wire type")
			}
			if err := _bodyHashValidateIdentifier(field.bytes, false); err != nil {
				return err
			}
		}
		raw = raw[field.size:]
	}
	return nil
}

func _transactionBodyHashChild(envelope []byte) ([]byte, error) {
	var body []byte
	for remaining := envelope; len(remaining) > 0; {
		field, err := _bodyHashConsumeField(remaining)
		if err != nil {
			return nil, err
		}
		if field.number == 1 {
			if body != nil || field.wireType != protowire.BytesType || len(field.bytes) == 0 {
				return nil, fmt.Errorf("missing, empty, repeated or mistyped batch child bodyBytes")
			}
			body = field.bytes
		}
		remaining = remaining[field.size:]
	}
	if body == nil {
		return nil, fmt.Errorf("missing batch child bodyBytes")
	}
	var signed services.SignedTransaction
	if err := protobuf.Unmarshal(envelope, &signed); err != nil {
		return nil, err
	}
	return body, nil
}

type _bodyHashField struct {
	number   protowire.Number
	wireType protowire.Type
	bytes    []byte
	size     int
}

func _bodyHashConsumeField(raw []byte) (_bodyHashField, error) {
	number, wireType, tagSize := protowire.ConsumeTag(raw)
	if tagSize < 0 {
		return _bodyHashField{}, protowire.ParseError(tagSize)
	}
	if tagSize > 5 {
		return _bodyHashField{}, fmt.Errorf("unsupported protobuf tag varint longer than five bytes")
	}
	if !number.IsValid() || (wireType != protowire.VarintType && wireType != protowire.BytesType && wireType != protowire.Fixed32Type && wireType != protowire.Fixed64Type) {
		return _bodyHashField{}, fmt.Errorf("invalid or unsupported protobuf wire field")
	}
	valueSize := protowire.ConsumeFieldValue(number, wireType, raw[tagSize:])
	if valueSize < 0 {
		return _bodyHashField{}, protowire.ParseError(valueSize)
	}
	field := _bodyHashField{number: number, wireType: wireType, size: tagSize + valueSize}
	if wireType == protowire.BytesType {
		field.bytes, _ = protowire.ConsumeBytes(raw[tagSize:])
		if valueSize-len(field.bytes) > 5 {
			return _bodyHashField{}, fmt.Errorf("unsupported protobuf length varint longer than five bytes")
		}
	}
	return field, nil
}
