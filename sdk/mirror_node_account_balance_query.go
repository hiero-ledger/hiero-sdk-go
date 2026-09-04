package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// MirrorNodeAccountBalanceQuery retrieves an account's or contract's hbar balance from the mirror
// node REST API, replacing AccountBalanceQuery. It is free and requires no operator.
//
// The mirror node trails the network by a few seconds, so results are not read-after-write
// consistent: a freshly created account fails with StatusInvalidAccountID until it is ingested.
type MirrorNodeAccountBalanceQuery struct {
	accountID   *AccountID
	maxAttempts uint64
}

// NewMirrorNodeAccountBalanceQuery creates a new MirrorNodeAccountBalanceQuery
func NewMirrorNodeAccountBalanceQuery() *MirrorNodeAccountBalanceQuery {
	return &MirrorNodeAccountBalanceQuery{}
}

// SetAccountID sets the account to query. Accepts shard.realm.num, an EVM address, a public key
// alias, or a contract addressed as AccountID{Shard, Realm, Account: contractID.Contract}.
func (q *MirrorNodeAccountBalanceQuery) SetAccountID(accountID AccountID) *MirrorNodeAccountBalanceQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the account to query, or the zero AccountID if not set.
func (q *MirrorNodeAccountBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}
	return *q.accountID
}

// SetMaxAttempts sets the total number of attempts (initial try + retries).
// Zero (the default) defers to the client.
func (q *MirrorNodeAccountBalanceQuery) SetMaxAttempts(maxAttempts uint64) *MirrorNodeAccountBalanceQuery {
	q.maxAttempts = maxAttempts
	return q
}

// GetMaxAttempts returns the configured retry budget, or 0 when it is left to the client.
func (q *MirrorNodeAccountBalanceQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

// Execute executes the query with the provided client
func (q *MirrorNodeAccountBalanceQuery) Execute(client *Client) (MirrorNodeAccountBalance, error) {
	if client == nil {
		return MirrorNodeAccountBalance{}, errNoClientProvided
	}

	if q.accountID == nil {
		return MirrorNodeAccountBalance{}, errMirrorNodeAccountBalanceQueryNoAccountID
	}

	if err := q.validateNetworkOnIDs(client); err != nil {
		return MirrorNodeAccountBalance{}, err
	}

	body, err := q.fetch(context.Background(), client)
	if err != nil {
		return MirrorNodeAccountBalance{}, err
	}

	return parseAccountBalances(body)
}

// fetch runs the request through the mirror HTTP layer: the Client owns the transport, this
// query contributes only its own attempt budget.
func (q *MirrorNodeAccountBalanceQuery) fetch(ctx context.Context, client *Client) ([]byte, error) {
	restClient, err := client.mirrorRestClient(client.mirrorHttpOptionsForQuery(q.maxAttempts))
	if err != nil {
		return nil, err
	}

	path, err := q.buildPath()
	if err != nil {
		return nil, err
	}

	resp, err := restClient.get(ctx, path)
	if err != nil {
		// A retryable status that outlived the budget still carries its response, so report it
		// as the non-200 it is rather than as a transport failure.
		if errors.Is(err, errMirrorHttpRetriesExhausted) {
			return nil, mirrorNodeStatusError(resp)
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if resp.statusCode != http.StatusOK {
		return nil, mirrorNodeStatusError(resp)
	}

	return resp.body, nil
}

// resolveAttempts picks the retry budget: query setting first,
// client default second, the mirror node default as the final fallback.
func (q *MirrorNodeAccountBalanceQuery) resolveAttempts(client *Client) uint64 {
	return uint64(client.mirrorHttpOptionsForQuery(q.maxAttempts).maxAttempts)
}

func (q *MirrorNodeAccountBalanceQuery) buildPath() (mirrorRestPath, error) {
	params := url.Values{}
	params.Set("account.id", q.accountID._MirrorNodePathID())

	return newMirrorRestPath("/balances?" + params.Encode())
}

func (q *MirrorNodeAccountBalanceQuery) buildURL(mirrorBaseURL string) string {
	path, err := q.buildPath()
	if err != nil {
		return ""
	}

	return resolveMirrorPath(mirrorBaseURL, path)
}

func (q *MirrorNodeAccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums || q.accountID == nil {
		return nil
	}

	// Only an ID with no account number is sent as an alias (see _MirrorNodePathID), and an alias
	// carries no checksum to validate. Anything sent by number is validated as usual.
	if q.accountID.Account == 0 && (q.accountID.AliasKey != nil || q.accountID.AliasEvmAddress != nil) {
		return nil
	}

	return q.accountID.ValidateChecksum(client)
}

// parseAccountBalances maps an empty balances list onto StatusInvalidAccountID, the status
// AccountBalanceQuery returned for an unknown account. The endpoint reports one as 200 with no rows.
func parseAccountBalances(body []byte) (MirrorNodeAccountBalance, error) {
	var raw accountBalancesResponseJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return MirrorNodeAccountBalance{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if raw.Balances == nil {
		return MirrorNodeAccountBalance{}, errors.New("mirror node response has no balances array")
	}

	// An existing account with no hbar returns "balance": 0, so an empty list only means unknown.
	if len(raw.Balances) == 0 {
		return MirrorNodeAccountBalance{}, ErrHederaPreCheckStatus{Status: StatusInvalidAccountID}
	}

	return MirrorNodeAccountBalance{Hbars: HbarFromTinybar(raw.Balances[0].Balance)}, nil
}

// An exact account.id filter matches one account, so there is never a next page to follow.
type accountBalancesResponseJSON struct {
	Balances []accountBalanceJSON `json:"balances"`
}

type accountBalanceJSON struct {
	Balance int64 `json:"balance"`
}
