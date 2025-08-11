package param

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
)

type ContractCreateTransactionParams struct {
	BytecodeFileId               *string                  `json:"bytecodeFileId"`
	AdminKey                     *string                  `json:"adminKey"`
	Gas                          *string                  `json:"gas"`
	InitialBalance               *string                  `json:"initialBalance"`
	ConstructorParameters        *string                  `json:"constructorParameters"`
	AutoRenewPeriod              *string                  `json:"autoRenewPeriod"`
	AutoRenewAccountId           *string                  `json:"autoRenewAccountId"`
	Memo                         *string                  `json:"memo"`
	StakedAccountId              *string                  `json:"stakedAccountId"`
	StakedNodeId                 *json.Number             `json:"stakedNodeId"`
	DeclineStakingReward         *bool                    `json:"declineStakingReward"`
	MaxAutomaticTokenAssociation *int32                   `json:"maxAutomaticTokenAssociations"`
	Initcode                     *string                  `json:"initcode"`
	CommonTransactionParams      *CommonTransactionParams `json:"commonTransactionParams"`
}
