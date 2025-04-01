package contract_helper

import (
	"encoding/hex"
	"fmt"
	"strconv"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
	"github.com/pkg/errors"
)

type ContractHelper struct {
	ContractID             hiero.ContractID
	stepResultValidators   map[int32]func(hiero.ContractFunctionResult) bool
	stepParameterSuppliers map[int32]func() *hiero.ContractFunctionParameters
	stepPayableAmounts     map[int32]*hiero.Hbar
	stepSigners            map[int32][]hiero.PrivateKey
	stepFeePayers          map[int32]*hiero.AccountID
	stepLogic              map[int32]func(address string)
}

func NewContractHelper(bytecode []byte, constructorParameters hiero.ContractFunctionParameters, client *hiero.Client) *ContractHelper {
	response, err := hiero.NewContractCreateFlow().
		SetBytecode(bytecode).
		SetGas(8000000).
		SetMaxChunks(30).
		SetConstructorParameters(&constructorParameters).
		Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := response.GetReceipt(client)
	if err != nil {
		panic(err)
	}
	if receipt.ContractID != nil {
		return &ContractHelper{
			ContractID:             *receipt.ContractID,
			stepResultValidators:   make(map[int32]func(hiero.ContractFunctionResult) bool),
			stepParameterSuppliers: make(map[int32]func() *hiero.ContractFunctionParameters),
			stepPayableAmounts:     make(map[int32]*hiero.Hbar),
			stepSigners:            make(map[int32][]hiero.PrivateKey),
			stepFeePayers:          make(map[int32]*hiero.AccountID),
			stepLogic:              make(map[int32]func(address string)),
		}
	}

	return &ContractHelper{}
}

func (ch *ContractHelper) SetResultValidatorForStep(stepIndex int32, validator func(hiero.ContractFunctionResult) bool) *ContractHelper {
	ch.stepResultValidators[stepIndex] = validator
	return ch
}

func (ch *ContractHelper) SetParameterSupplierForStep(stepIndex int32, supplier func() *hiero.ContractFunctionParameters) *ContractHelper {
	ch.stepParameterSuppliers[stepIndex] = supplier
	return ch
}

func (ch *ContractHelper) SetPayableAmountForStep(stepIndex int32, amount hiero.Hbar) *ContractHelper {
	ch.stepPayableAmounts[stepIndex] = &amount
	return ch
}

func (ch *ContractHelper) AddSignerForStep(stepIndex int32, signer hiero.PrivateKey) *ContractHelper {
	if _, ok := ch.stepSigners[stepIndex]; ok {
		ch.stepSigners[stepIndex] = append(ch.stepSigners[stepIndex], signer)
	} else {
		ch.stepSigners[stepIndex] = make([]hiero.PrivateKey, 0)
		ch.stepSigners[stepIndex] = append(ch.stepSigners[stepIndex], signer)
	}

	return ch
}

func (ch *ContractHelper) SetFeePayerForStep(stepIndex int32, account hiero.AccountID, accountKey hiero.PrivateKey) *ContractHelper {
	ch.stepFeePayers[stepIndex] = &account
	return ch.AddSignerForStep(stepIndex, accountKey)
}

func (ch *ContractHelper) SetStepLogic(stepIndex int32, specialFunction func(address string)) *ContractHelper {
	ch.stepLogic[stepIndex] = specialFunction
	return ch
}

func (ch *ContractHelper) GetResultValidator(stepIndex int32) func(hiero.ContractFunctionResult) bool {
	if _, ok := ch.stepResultValidators[stepIndex]; ok {
		return ch.stepResultValidators[stepIndex]
	}

	return func(result hiero.ContractFunctionResult) bool {
		responseStatus := hiero.Status(result.GetInt32(0))
		isValid := responseStatus == hiero.StatusSuccess
		if !isValid {
			println("Encountered invalid response status", responseStatus.String())
		}
		return isValid
	}
}

func (ch *ContractHelper) GetParameterSupplier(stepIndex int32) func() *hiero.ContractFunctionParameters {
	if _, ok := ch.stepParameterSuppliers[stepIndex]; ok {
		return ch.stepParameterSuppliers[stepIndex]
	}

	return func() *hiero.ContractFunctionParameters {
		return nil
	}
}

func (ch *ContractHelper) GetPayableAmount(stepIndex int32) *hiero.Hbar {
	return ch.stepPayableAmounts[stepIndex]
}

func (ch *ContractHelper) GetSigners(stepIndex int32) []hiero.PrivateKey {
	if _, ok := ch.stepSigners[stepIndex]; ok {
		return ch.stepSigners[stepIndex]
	}

	return []hiero.PrivateKey{}
}

func (ch *ContractHelper) ExecuteSteps(firstStep int32, lastStep int32, client *hiero.Client) (*ContractHelper, error) {
	for stepIndex := firstStep; stepIndex <= lastStep; stepIndex++ {
		println("Attempting to execuite step", stepIndex)

		transaction := hiero.NewContractExecuteTransaction().
			SetContractID(ch.ContractID).
			SetGas(10000000)

		payableAmount := ch.GetPayableAmount(stepIndex)
		if payableAmount != nil {
			transaction.SetPayableAmount(*payableAmount)
		}

		functionName := "step" + strconv.Itoa(int(stepIndex))
		parameters := ch.GetParameterSupplier(stepIndex)()
		if parameters != nil {
			transaction.SetFunction(functionName, parameters)
		} else {
			transaction.SetFunction(functionName, nil)
		}

		if feePayerAccountID, ok := ch.stepFeePayers[stepIndex]; ok {
			transaction.SetTransactionID(hiero.TransactionIDGenerate(*feePayerAccountID))
		}

		frozen, err := transaction.FreezeWith(client)
		if err != nil {
			return &ContractHelper{}, err
		}
		for _, signer := range ch.GetSigners(stepIndex) {
			frozen.Sign(signer)
		}

		response, err := frozen.Execute(client)
		if err != nil {
			return &ContractHelper{}, err
		}

		record, err := response.GetRecord(client)
		if err != nil {
			return &ContractHelper{}, err
		}

		functionResult, err := record.GetContractExecuteResult()
		if err != nil {
			return &ContractHelper{}, err
		}

		if ch.stepLogic[stepIndex] != nil {
			address := functionResult.GetAddress(1)
			if function, exists := ch.stepLogic[stepIndex]; exists && function != nil {
				function(hex.EncodeToString(address))
			}
		}

		if ch.GetResultValidator(stepIndex)(functionResult) {
			fmt.Printf("Step %d completed, and returned valid result. (TransactionId %s)", stepIndex, record.TransactionID.String())
		} else {
			return &ContractHelper{}, errors.New(fmt.Sprintf("Step %d returned invalid result", stepIndex))
		}
	}

	return ch, nil
}
