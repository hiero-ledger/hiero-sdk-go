package utils

import (
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func SetTokenSupplyParams(transaction *hiero.TokenCreateTransaction, params param.CreateTokenParams) error {
	if params.MaxSupply != nil {
		maxSupply, err := strconv.ParseInt(*params.MaxSupply, 10, 64)
		if err != nil {
			return err
		}
		transaction.SetMaxSupply(maxSupply)
	}
	if params.InitialSupply != nil {
		initialSupply, err := strconv.ParseInt(*params.InitialSupply, 10, 64)
		if err != nil {
			return err
		}
		transaction.SetInitialSupply(uint64(initialSupply))
	}
	return nil
}

func SetTokenTypes(transaction *hiero.TokenCreateTransaction, params param.CreateTokenParams) error {
	if params.TokenType != nil {
		switch *params.TokenType {
		case "ft":
			transaction.SetTokenType(hiero.TokenTypeFungibleCommon)
		case "nft":
			transaction.SetTokenType(hiero.TokenTypeNonFungibleUnique)
		default:
			return response.InvalidParams.WithData("Invalid token type")
		}
	}

	if params.SupplyType != nil {
		switch *params.SupplyType {
		case "finite":
			transaction.SetSupplyType(hiero.TokenSupplyTypeFinite)
		case "infinite":
			transaction.SetSupplyType(hiero.TokenSupplyTypeInfinite)
		default:
			return response.InvalidParams.WithData("Invalid supply type")
		}
	}
	return nil
}
