package param

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
)

type ContractCreateTransactionParams struct {
	BaseParams
	BytecodeFileId               *string      `json:"bytecodeFileId"`
	AdminKey                     *string      `json:"adminKey"`
	Gas                          *string      `json:"gas"`
	InitialBalance               *string      `json:"initialBalance"`
	ConstructorParameters        *string      `json:"constructorParameters"`
	AutoRenewPeriod              *string      `json:"autoRenewPeriod"`
	AutoRenewAccountId           *string      `json:"autoRenewAccountId"`
	Memo                         *string      `json:"memo"`
	StakedAccountId              *string      `json:"stakedAccountId"`
	StakedNodeId                 *json.Number `json:"stakedNodeId"`
	DeclineStakingReward         *bool        `json:"declineStakingReward"`
	MaxAutomaticTokenAssociation *int32       `json:"maxAutomaticTokenAssociations"`
	Initcode                     *string      `json:"initcode"`
}

type ContractUpdateTransactionParams struct {
	BaseParams
	ContractId                   *string      `json:"contractId"`
	AdminKey                     *string      `json:"adminKey"`
	AutoRenewPeriod              *string      `json:"autoRenewPeriod"`
	ExpirationTime               *string      `json:"expirationTime"`
	Memo                         *string      `json:"memo"`
	AutoRenewAccountId           *string      `json:"autoRenewAccountId"`
	MaxAutomaticTokenAssociation *int32       `json:"maxAutomaticTokenAssociations"`
	StakedAccountId              *string      `json:"stakedAccountId"`
	StakedNodeId                 *json.Number `json:"stakedNodeId"`
	DeclineStakingReward         *bool        `json:"declineStakingReward"`
}

type ContractDeleteTransactionParams struct {
	BaseParams
	ContractId         *string `json:"contractId"`
	TransferContractId *string `json:"transferContractId"`
	TransferAccountId  *string `json:"transferAccountId"`
	PermanentRemoval   *bool   `json:"permanentRemoval"`
}

type ContractExecuteTransactionParams struct {
	BaseParams
	ContractId         *string `json:"contractId"`
	Gas                *string `json:"gas"`
	PayableAmount      *string `json:"amount"`
	FunctionParameters *string `json:"functionParameters"`
}
