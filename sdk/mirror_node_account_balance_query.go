package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// MirrorNodeAccountBalanceQuery retrieves an account's or contract's hbar balance from the mirror
// node REST API, replacing AccountBalanceQuery. It is free and requires no operator. The mirror node
// trails the network by a few seconds, so results are not read-after-write consistent.
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

// Execute executes the Query with the provided client
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

	endpoint, err := q.resolveEndpoint(client)
	if err != nil {
		return MirrorNodeAccountBalance{}, err
	}

	body, err := fetchAccountBalances(client, endpoint, q.resolveAttempts(client))
	if err != nil {
		return MirrorNodeAccountBalance{}, err
	}

	return parseAccountBalances(body)
}

func (q *MirrorNodeAccountBalanceQuery) resolveEndpoint(client *Client) (string, error) {
	mirrorUrl, err := mirrorNodeRestBaseURL(client)
	if err != nil {
		return "", err
	}

	return q.buildURL(mirrorUrl), nil
}

// resolveAttempts picks the retry budget: query setting first,
// client default second, the SDK default as the final fallback.
func (q *MirrorNodeAccountBalanceQuery) resolveAttempts(client *Client) uint64 {
	if q.maxAttempts > 0 {
		return q.maxAttempts
	}
	if clientMax := client.GetMaxAttempts(); clientMax > 0 {
		return uint64(clientMax)
	}
	return maxAttempts
}

// buildURL escapes the filter: an alias is opaque bytes and must not alter the query string.
func (q *MirrorNodeAccountBalanceQuery) buildURL(mirrorBaseURL string) string {
	params := url.Values{}
	params.Set("account.id", q.accountID._MirrorNodePathID())

	return fmt.Sprintf("%s/balances?%s", mirrorBaseURL, params.Encode())
}

func (q *MirrorNodeAccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

// fetchAccountBalances GETs the balances endpoint through the shared retry core.
func fetchAccountBalances(client *Client, endpoint string, attempts uint64) ([]byte, error) {
	resp, err := mirrorNodeGetWithRetry(client, endpoint, attempts, mirrorNodeDefaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return mirrorNodeReadBody(resp)
}

// parseAccountBalances reads an empty list as a zero balance: the endpoint returns one for an
// unknown account rather than a 404, so zero is not proof the account exists.
func parseAccountBalances(body []byte) (MirrorNodeAccountBalance, error) {
	var raw accountBalancesResponseJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return MirrorNodeAccountBalance{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(raw.Balances) == 0 {
		return MirrorNodeAccountBalance{Hbars: HbarFromTinybar(0)}, nil
	}

	return MirrorNodeAccountBalance{Hbars: HbarFromTinybar(raw.Balances[0].Balance)}, nil
}

// An exact account.id filter matches one account, so there is never a next page to follow.
type accountBalancesResponseJSON struct {
	Balances []accountBalanceJSON `json:"balances"`
}

type accountBalanceJSON struct {
	Account string `json:"account"`
	Balance int64  `json:"balance"`
}
