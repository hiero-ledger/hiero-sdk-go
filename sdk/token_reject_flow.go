package hiero

// SPDX-License-Identifier: Apache-2.0

type TokenRejectFlow struct {
	TokenRejectTransaction     *TokenRejectTransaction
	TokenDissociateTransaction *TokenDissociateTransaction
	ownerID                    *AccountID
	tokenIDs                   []TokenID
	nftIDs                     []NftID
	freezeWithClient           *Client
	signPrivateKey             *PrivateKey
	signPublicKey              *PublicKey
	transactionSigner          *TransactionSigner
}

func NewTokenRejectFlow() *TokenRejectFlow {
	tx := TokenRejectFlow{
		TokenRejectTransaction:     NewTokenRejectTransaction(),
		TokenDissociateTransaction: NewTokenDissociateTransaction(),
	}
	return &tx
}

// SetOwnerID Sets the account which owns the tokens to be rejected
func (tx *TokenRejectFlow) SetOwnerID(ownerID AccountID) *TokenRejectFlow {
	tx.ownerID = &ownerID
	return tx
}

// GetOwnerID Gets the account which owns the tokens to be rejected
func (tx *TokenRejectFlow) GetOwnerID() AccountID {
	if tx.ownerID == nil {
		return AccountID{}
	}
	return *tx.ownerID
}

// SetTokenIDs Sets the tokens to be rejected
func (tx *TokenRejectFlow) SetTokenIDs(ids ...TokenID) *TokenRejectFlow {
	tx.tokenIDs = make([]TokenID, len(ids))
	copy(tx.tokenIDs, ids)

	return tx
}

// AddTokenID Adds a token to be rejected
func (tx *TokenRejectFlow) AddTokenID(id TokenID) *TokenRejectFlow {
	tx.tokenIDs = append(tx.tokenIDs, id)
	return tx
}

// GetTokenIDs Gets the tokens to be rejected
func (tx *TokenRejectFlow) GetTokenIDs() []TokenID {
	return tx.tokenIDs
}

// SetNftIDs Sets the NFTs to be rejected
func (tx *TokenRejectFlow) SetNftIDs(ids ...NftID) *TokenRejectFlow {
	tx.nftIDs = make([]NftID, len(ids))
	copy(tx.nftIDs, ids)

	return tx
}

// AddNftID Adds an NFT to be rejected
func (tx *TokenRejectFlow) AddNftID(id NftID) *TokenRejectFlow {
	tx.nftIDs = append(tx.nftIDs, id)
	return tx
}

// GetNftIDs Gets the NFTs to be rejected
func (tx *TokenRejectFlow) GetNftIDs() []NftID {
	return tx.nftIDs
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenRejectFlow) Sign(privateKey PrivateKey) *TokenRejectFlow {
	tx.signPrivateKey = &privateKey
	return tx
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenRejectFlow) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRejectFlow {
	tx.signPublicKey = &publicKey
	tx.transactionSigner = &signer
	return tx
}

func (tx *TokenRejectFlow) FreezeWith(client *Client) (*TokenRejectFlow, error) {
	tx.freezeWithClient = client
	return tx, nil
}

func (tx *TokenRejectFlow) fillTokenDissociate() error {
	if tx.ownerID != nil {
		tx.TokenDissociateTransaction.SetAccountID(*tx.ownerID)
	}

	tokenIDs := make([]TokenID, 0)
	if tx.tokenIDs != nil {
		tokenIDs = append(tokenIDs, tx.tokenIDs...)
	}

	if tx.nftIDs != nil {
		seenTokenIDs := make(map[TokenID]struct{})
		for _, nftID := range tx.nftIDs {
			if _, exists := seenTokenIDs[nftID.TokenID]; !exists {
				seenTokenIDs[nftID.TokenID] = struct{}{}
				tokenIDs = append(tokenIDs, nftID.TokenID)
			}
		}
	}

	if len(tokenIDs) != 0 {
		tx.TokenDissociateTransaction.SetTokenIDs(tokenIDs...)
	}

	if tx.freezeWithClient != nil {
		_, err := tx.TokenDissociateTransaction.FreezeWith(tx.freezeWithClient)
		if err != nil {
			return err
		}
	}

	if tx.signPrivateKey != nil {
		tx.TokenDissociateTransaction = tx.TokenDissociateTransaction.Sign(*tx.signPrivateKey)
	}

	if tx.signPublicKey != nil && tx.transactionSigner != nil {
		tx.TokenDissociateTransaction = tx.TokenDissociateTransaction.SignWith(*tx.signPublicKey, *tx.transactionSigner)
	}
	return nil
}

func (tx *TokenRejectFlow) fillTokenReject() error {
	if tx.ownerID != nil {
		tx.TokenRejectTransaction.SetOwnerID(*tx.ownerID)
	}

	if tx.tokenIDs != nil {
		tx.TokenRejectTransaction.SetTokenIDs(tx.tokenIDs...)
	}

	if tx.nftIDs != nil {
		tx.TokenRejectTransaction.SetNftIDs(tx.nftIDs...)
	}

	if tx.freezeWithClient != nil {
		_, err := tx.TokenRejectTransaction.FreezeWith(tx.freezeWithClient)
		if err != nil {
			return err
		}
	}

	if tx.signPrivateKey != nil {
		tx.TokenRejectTransaction = tx.TokenRejectTransaction.Sign(*tx.signPrivateKey)
	}

	if tx.signPublicKey != nil && tx.transactionSigner != nil {
		tx.TokenRejectTransaction = tx.TokenRejectTransaction.SignWith(*tx.signPublicKey, *tx.transactionSigner)
	}

	return nil
}

func (tx *TokenRejectFlow) Execute(client *Client) (TransactionResponse, error) {
	err := tx.fillTokenReject()
	if err != nil {
		return TransactionResponse{}, err
	}

	err = tx.fillTokenDissociate()
	if err != nil {
		return TransactionResponse{}, err
	}

	tokenRejectResponse, err := tx.TokenRejectTransaction.Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = tokenRejectResponse.GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	tokenDissociateResponse, err := tx.TokenDissociateTransaction.Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = tokenDissociateResponse.GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return tokenRejectResponse, nil
}
