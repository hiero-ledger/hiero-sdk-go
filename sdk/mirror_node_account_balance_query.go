package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Caps the /tokens pagination walk so a misbehaving mirror node cannot loop forever.
const mirrorNodeAccountBalanceMaxPages = 100

// The mirror node's maximum page size; the default of 25 would cap the walk at 2,500 tokens.
const mirrorNodeAccountBalanceTokensPageSize = 100

// MirrorNodeAccountBalanceQuery retrieves the hbar and token balances of an account or contract
// from the mirror node REST API. It replaces the deprecated AccountBalanceQuery.
//
// Hbars come from GET /api/v1/accounts/{id}; token balances and decimals from the paginated
// GET /api/v1/accounts/{id}/tokens (the accounts response's embedded token list is truncated).
type MirrorNodeAccountBalanceQuery struct {
	accountID   *AccountID
	contractID  *ContractID
	timeout     time.Duration
	maxAttempts uint64
}

// mirrorNodeAccountBalanceResponse is the hbar-balance subset of GET /api/v1/accounts/{id}.
type mirrorNodeAccountBalanceResponse struct {
	Balance struct {
		Balance int64 `json:"balance"`
	} `json:"balance"`
}

// mirrorNodeAccountTokensResponse is one page of GET /api/v1/accounts/{id}/tokens.
type mirrorNodeAccountTokensResponse struct {
	Tokens []struct {
		TokenID  string `json:"token_id"`
		Balance  uint64 `json:"balance"`
		Decimals uint64 `json:"decimals"`
	} `json:"tokens"`
	Links *linksJSON `json:"links"`
}

// NewMirrorNodeAccountBalanceQuery creates a query for an account's or contract's balance.
func NewMirrorNodeAccountBalanceQuery() *MirrorNodeAccountBalanceQuery {
	return &MirrorNodeAccountBalanceQuery{
		timeout:     mirrorNodeDefaultTimeout,
		maxAttempts: mirrorNodeDefaultMaxAttempts,
	}
}

// SetAccountID sets the AccountID to query. Account and contract are mutually exclusive, so this
// clears any contract ID already set.
func (q *MirrorNodeAccountBalanceQuery) SetAccountID(accountID AccountID) *MirrorNodeAccountBalanceQuery {
	q.accountID = &accountID
	q.contractID = nil
	return q
}

// GetAccountID returns the AccountID to query.
func (q *MirrorNodeAccountBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}
	return *q.accountID
}

// SetContractID sets the ContractID to query. Account and contract are mutually exclusive, so this
// clears any account ID already set.
func (q *MirrorNodeAccountBalanceQuery) SetContractID(contractID ContractID) *MirrorNodeAccountBalanceQuery {
	q.contractID = &contractID
	q.accountID = nil
	return q
}

// GetContractID returns the ContractID to query.
func (q *MirrorNodeAccountBalanceQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}
	return *q.contractID
}

// SetTimeout sets the per-request timeout for the mirror node REST call. A timeout of 0 disables it.
func (q *MirrorNodeAccountBalanceQuery) SetTimeout(timeout time.Duration) *MirrorNodeAccountBalanceQuery {
	q.timeout = timeout
	return q
}

// GetTimeout returns the per-request timeout for the mirror node REST call.
func (q *MirrorNodeAccountBalanceQuery) GetTimeout() time.Duration {
	return q.timeout
}

// SetMaxAttempts sets how many times a transient (transport or 5xx/429) failure is retried.
func (q *MirrorNodeAccountBalanceQuery) SetMaxAttempts(maxAttempts uint64) *MirrorNodeAccountBalanceQuery {
	q.maxAttempts = maxAttempts
	return q
}

// GetMaxAttempts returns how many times a transient failure is retried.
func (q *MirrorNodeAccountBalanceQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

// Execute returns the configured account's or contract's hbar balance, token balances and decimals.
func (q *MirrorNodeAccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}

	if err := q.validateNetworkOnIDs(client); err != nil {
		return AccountBalance{}, err
	}

	var idStr string
	switch {
	case q.accountID != nil:
		idStr = q.accountID._MirrorNodePathID()
	case q.contractID != nil:
		idStr = q.contractID._MirrorNodePathID()
	default:
		return AccountBalance{}, errors.New("either an account ID or a contract ID must be set")
	}

	baseURL, err := mirrorNodeRestBaseURL(client)
	if err != nil {
		return AccountBalance{}, err
	}

	hbars, err := q.fetchHbarBalance(client, baseURL, idStr)
	if err != nil {
		return AccountBalance{}, err
	}

	balances, decimals, tokens, err := q.fetchTokenBalances(client, baseURL, idStr)
	if err != nil {
		return AccountBalance{}, err
	}

	return _AccountBalanceFromMirrorNode(hbars, balances, decimals, tokens), nil
}

// validateNetworkOnIDs checks the configured ID's checksum against the client's network when
// auto-validation is enabled, as the consensus node queries do.
func (q *MirrorNodeAccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.accountID != nil {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

// fetchHbarBalance retrieves the hbar balance from GET /api/v1/accounts/{id}.
func (q *MirrorNodeAccountBalanceQuery) fetchHbarBalance(client *Client, baseURL, idStr string) (Hbar, error) {
	body, err := q.getPage(client, fmt.Sprintf("%s/accounts/%s", baseURL, idStr))
	if err != nil {
		return Hbar{}, err
	}

	var parsed mirrorNodeAccountBalanceResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return Hbar{}, fmt.Errorf("failed to decode mirror node response: %w", err)
	}

	return HbarFromTinybar(parsed.Balance.Balance), nil
}

// fetchTokenBalances accumulates token balances and decimals across every /tokens page, keyed both
// by token ID string (for the public maps) and by TokenID (for the deprecated Token map).
func (q *MirrorNodeAccountBalanceQuery) fetchTokenBalances(client *Client, baseURL, idStr string) (map[string]uint64, map[string]uint64, map[TokenID]uint64, error) {
	balances := make(map[string]uint64)
	decimals := make(map[string]uint64)
	tokens := make(map[TokenID]uint64)

	err := mirrorNodeWalkPages(
		// The limit only needs to be set on the first page; the mirror node echoes it in links.next.
		fmt.Sprintf("%s/accounts/%s/tokens?limit=%d", baseURL, idStr, mirrorNodeAccountBalanceTokensPageSize),
		mirrorNodeAccountBalanceMaxPages,
		func(pageURL string) ([]byte, error) {
			return q.getPage(client, pageURL)
		},
		func(body []byte) (*string, error) {
			var parsed mirrorNodeAccountTokensResponse
			if err := json.Unmarshal(body, &parsed); err != nil {
				return nil, fmt.Errorf("failed to decode mirror node response: %w", err)
			}

			for _, token := range parsed.Tokens {
				// An unparseable token ID means the balance would be silently wrong, so fail
				// rather than report a partial balance as complete.
				tokenID, err := TokenIDFromString(token.TokenID)
				if err != nil {
					return nil, fmt.Errorf("mirror node returned an invalid token ID %q: %w", token.TokenID, err)
				}

				key := tokenID.String()
				balances[key] = token.Balance
				decimals[key] = token.Decimals
				tokens[tokenID] = token.Balance
			}

			if parsed.Links == nil {
				return nil, nil
			}
			return parsed.Links.Next, nil
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return balances, decimals, tokens, nil
}

// getPage issues a retrying GET and returns the body, turning failures into errors.
func (q *MirrorNodeAccountBalanceQuery) getPage(client *Client, url string) ([]byte, error) {
	resp, err := mirrorNodeGetWithRetry(client, url, q.maxAttempts, q.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to mirror node: %w", err)
	}

	return mirrorNodeReadBody(resp)
}
