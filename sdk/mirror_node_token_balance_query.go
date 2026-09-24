package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// The largest page the mirror node returns.
const mirrorNodeTokenBalancePageLimit = "100"

// MirrorNodeTokenBalanceQuery retrieves one page of an account's token balances from the mirror
// node REST API. It is free and requires no operator.
//
// The mirror node trails the network by a few seconds, so results are not read-after-write
// consistent: a freshly created account fails with StatusInvalidAccountID until it is ingested.
type MirrorNodeTokenBalanceQuery struct {
	accountID *AccountID
	tokenID   *TokenID
	nextPage  string
}

// NewMirrorNodeTokenBalanceQuery creates a new MirrorNodeTokenBalanceQuery
func NewMirrorNodeTokenBalanceQuery() *MirrorNodeTokenBalanceQuery {
	return &MirrorNodeTokenBalanceQuery{}
}

// SetAccountID sets the account to query. Accepts shard.realm.num, an EVM address or a public key
// alias. Required unless a next page is set.
func (q *MirrorNodeTokenBalanceQuery) SetAccountID(accountID AccountID) *MirrorNodeTokenBalanceQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the account to query, or the zero AccountID if not set.
func (q *MirrorNodeTokenBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}
	return *q.accountID
}

// SetTokenID limits the query to one token. An unassociated token returns an empty page, not an error.
func (q *MirrorNodeTokenBalanceQuery) SetTokenID(tokenID TokenID) *MirrorNodeTokenBalanceQuery {
	q.tokenID = &tokenID
	return q
}

// GetTokenID returns the token to query, or the zero TokenID if not set.
func (q *MirrorNodeTokenBalanceQuery) GetTokenID() TokenID {
	if q.tokenID == nil {
		return TokenID{}
	}
	return *q.tokenID
}

// SetNextPage sets the MirrorNodeTokenBalancePage.Next cursor of a previous page. The cursor
// already names the account, so it cannot be combined with SetAccountID or SetTokenID.
func (q *MirrorNodeTokenBalanceQuery) SetNextPage(next string) *MirrorNodeTokenBalanceQuery {
	q.nextPage = next
	return q
}

// GetNextPage returns the cursor set with SetNextPage, or an empty string if not set.
func (q *MirrorNodeTokenBalanceQuery) GetNextPage() string {
	return q.nextPage
}

// Execute executes the query with the provided client
func (q *MirrorNodeTokenBalanceQuery) Execute(client *Client) (MirrorNodeTokenBalancePage, error) {
	if client == nil {
		return MirrorNodeTokenBalancePage{}, errNoClientProvided
	}

	path, err := q.buildPath(client)
	if err != nil {
		return MirrorNodeTokenBalancePage{}, err
	}

	restClient, err := client.mirrorRestClient(path, client.mirrorHttpPolicy())
	if err != nil {
		return MirrorNodeTokenBalancePage{}, err
	}

	resp, err := restClient.get(path, CancellationNone())
	switch {
	case errors.Is(err, errMirrorHttpRetriesExhausted):
		return MirrorNodeTokenBalancePage{}, mirrorNodeStatusError(resp)
	case err != nil:
		return MirrorNodeTokenBalancePage{}, fmt.Errorf("failed to send request: %w", err)
	case resp.statusCode == http.StatusNotFound:
		// This endpoint reports an unknown account as 404.
		return MirrorNodeTokenBalancePage{}, ErrHederaPreCheckStatus{Status: StatusInvalidAccountID}
	case resp.statusCode != http.StatusOK:
		return MirrorNodeTokenBalancePage{}, mirrorNodeStatusError(resp)
	}

	page, err := parseTokenRelationships(resp.body)
	if err != nil {
		return MirrorNodeTokenBalancePage{}, err
	}
	// A token.id filter's next link repeats the request, so following it would loop forever.
	if q.tokenID != nil {
		page.Next = ""
	}

	return page, nil
}

// buildPath validates the query before any request is made and returns the path to fetch.
func (q *MirrorNodeTokenBalanceQuery) buildPath(client *Client) (mirrorNodeRestPath, error) {
	if q.nextPage != "" {
		if q.accountID != nil || q.tokenID != nil {
			return "", errMirrorNodeTokenBalanceQueryNextPageConflict
		}
		return tokenBalancePagePath(q.nextPage)
	}
	if q.accountID == nil {
		return "", errMirrorNodeTokenBalanceQueryNoAccountID
	}
	if err := q.validateNetworkOnIDs(client); err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("limit", mirrorNodeTokenBalancePageLimit)
	if q.tokenID != nil {
		params.Set("token.id", q.tokenID.String())
	} else {
		params.Set("order", "asc")
	}

	return newMirrorNodeRestPath("/accounts/" + q.accountID._MirrorNodePathID() + "/tokens?" + params.Encode())
}

// tokenBalancePagePath accepts only a cursor for this endpoint, so SetNextPage cannot fetch other paths.
func tokenBalancePagePath(next string) (mirrorNodeRestPath, error) {
	path, err := nextPagePath(next)
	if err != nil {
		return "", fmt.Errorf("invalid next page %q: %w", next, err)
	}

	endpoint, _, _ := strings.Cut(path.String(), "?")
	if !strings.HasPrefix(endpoint, "/accounts/") || !strings.HasSuffix(endpoint, "/tokens") {
		return "", fmt.Errorf("invalid next page %q: not a token balance page", next)
	}

	return path, nil
}

func (q *MirrorNodeTokenBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if !client.autoValidateChecksums {
		return nil
	}

	// Only an ID with no account number is sent as an alias, and an alias carries no checksum.
	isAlias := q.accountID.Account == 0 && (q.accountID.AliasKey != nil || q.accountID.AliasEvmAddress != nil)
	if !isAlias {
		if err := q.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}
	if q.tokenID != nil {
		return q.tokenID.ValidateChecksum(client)
	}

	return nil
}

func parseTokenRelationships(body []byte) (MirrorNodeTokenBalancePage, error) {
	var raw tokenRelationshipsResponseJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return MirrorNodeTokenBalancePage{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if raw.Tokens == nil {
		return MirrorNodeTokenBalancePage{}, errors.New("mirror node response has no tokens array")
	}

	page := MirrorNodeTokenBalancePage{Tokens: make([]MirrorNodeTokenBalance, 0, len(raw.Tokens))}
	for _, rel := range raw.Tokens {
		tokenID, err := TokenIDFromString(rel.TokenID)
		if err != nil {
			return MirrorNodeTokenBalancePage{}, fmt.Errorf("mirror node returned an invalid token ID %q: %w", rel.TokenID, err)
		}
		page.Tokens = append(page.Tokens, MirrorNodeTokenBalance{
			TokenID:              tokenID,
			Balance:              rel.Balance,
			Decimals:             rel.Decimals,
			AutomaticAssociation: rel.AutomaticAssociation,
			CreatedTimestamp:     rel.CreatedTimestamp,
			FreezeStatus:         rel.FreezeStatus,
			KycStatus:            rel.KycStatus,
		})
	}
	if raw.Links.Next != nil {
		page.Next = *raw.Links.Next
	}

	return page, nil
}

type tokenRelationshipsResponseJSON struct {
	Tokens []tokenRelationshipJSON `json:"tokens"`
	Links  linksJSON               `json:"links"`
}

type tokenRelationshipJSON struct {
	TokenID              string `json:"token_id"`
	Balance              uint64 `json:"balance"`
	Decimals             uint64 `json:"decimals"`
	AutomaticAssociation *bool  `json:"automatic_association"`
	CreatedTimestamp     string `json:"created_timestamp"`
	FreezeStatus         string `json:"freeze_status"`
	KycStatus            string `json:"kyc_status"`
}
