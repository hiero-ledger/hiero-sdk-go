package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
)

// AddressBookQuery query an address book for its list of nodes
type AddressBookQuery struct {
	attempt     uint64
	maxAttempts uint64
	fileID      *FileID
	limit       int32
}

// Query the mirror node for the address book.
func NewAddressBookQuery() *AddressBookQuery {
	return &AddressBookQuery{
		fileID: nil,
		limit:  0,
	}
}

// SetFileID set the ID of the address book file on the network. Can be either 0.0.101 or 0.0.102.
func (q *AddressBookQuery) SetFileID(id FileID) *AddressBookQuery {
	q.fileID = &id
	return q
}

func (q *AddressBookQuery) GetFileID() FileID {
	if q.fileID == nil {
		return FileID{}
	}

	return *q.fileID
}

// SetLimit
// Set the maximum number of node addresses to receive before stopping.
// If not set or set to zero it will return all node addresses in the database.
func (q *AddressBookQuery) SetLimit(limit int32) *AddressBookQuery {
	q.limit = limit
	return q
}

func (q *AddressBookQuery) GetLimit() int32 {
	return q.limit
}

func (q *AddressBookQuery) SetMaxAttempts(maxAttempts uint64) *AddressBookQuery {
	q.maxAttempts = maxAttempts
	return q
}

func (q *AddressBookQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

func (q *AddressBookQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.fileID != nil {
		if err := q.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *AddressBookQuery) build() *mirror.AddressBookQuery {
	body := &mirror.AddressBookQuery{
		Limit: q.limit,
	}
	if q.fileID != nil {
		body.FileId = q.fileID._ToProtobuf()
	}

	return body
}

// RecvStream is a generic interface for any gRPC client stream that has a Recv() method.
// This allows processProtoMessageStream to work with any stream type, not just specific ones.
type RecvStream[T any] interface {
	Recv() (T, error)
}

// Execute executes the Query with the provided client
func (q *AddressBookQuery) Execute(client *Client) (NodeAddressBook, error) {
	err := q.validateNetworkOnIDs(client)
	if err != nil {
		return NodeAddressBook{}, err
	}

	pbBody := q.build()

	mirrorNode, err := client.mirrorNetwork._GetNextMirrorNode()
	if err != nil {
		return NodeAddressBook{}, err
	}

	channel, err := mirrorNode._GetNetworkServiceClient()
	if err != nil {
		return NodeAddressBook{}, err
	}

	ctx := client.networkUpdateContext

	stream, err := channel.GetNodes(ctx, pbBody)
	if err != nil {
		return NodeAddressBook{}, err
	}

	resultStream := processProtoMessageStream(ctx, stream, q.attempt, q.maxAttempts)

	results := make([]NodeAddress, 0)

	for result := range resultStream {
		if result.err != nil {
			return NodeAddressBook{}, result.err
		}

		results = append(results, _NodeAddressFromProtobuf(result.data))
	}

	return NodeAddressBook{
		NodeAddresses: results,
	}, nil
}
