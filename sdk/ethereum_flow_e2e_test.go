//go:build all || e2e
// +build all e2e

package hiero

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

var LARGE_SMART_CONTRACT_BYTECODE = "60806040526040518060400160405280600581526020017f68656c6c6f0000000000000000000000000000000000000000000000000000008152505f90816100479190610293565b50348015610053575f80fd5b50610362565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806100d457607f821691505b6020821081036100e7576100e6610090565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026101497fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261010e565b610153868361010e565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61019761019261018d8461016b565b610174565b61016b565b9050919050565b5f819050919050565b6101b08361017d565b6101c46101bc8261019e565b84845461011a565b825550505050565b5f90565b6101d86101cc565b6101e38184846101a7565b505050565b5b81811015610206576101fb5f826101d0565b6001810190506101e9565b5050565b601f82111561024b5761021c816100ed565b610225846100ff565b81016020851015610234578190505b610248610240856100ff565b8301826101e8565b50505b505050565b5f82821c905092915050565b5f61026b5f1984600802610250565b1980831691505092915050565b5f610283838361025c565b9150826002028217905092915050565b61029c82610059565b67ffffffffffffffff8111156102b5576102b4610063565b5b6102bf82546100bd565b6102ca82828561020a565b5f60209050601f8311600181146102fb575f84156102e9578287015190505b6102f38582610278565b86555061035a565b601f198416610309866100ed565b5f5b828110156103305784890151825560018201915060208501945060208101905061030b565b8683101561034d5784890151610349601f89168261025c565b8355505b6001600288020188555050505b505050505050565b6113678061036f5f395ff3fe608060405234801561000f575f80fd5b50600436106100f3575f3560e01c80636a4b285811610095578063ce6d41de11610064578063ce6d41de14610281578063e2842d791461029f578063f8b2cb4f146102bd578063fb1669ca146102ed576100f3565b80636a4b2858146101e75780638ada066e14610217578063b633941814610235578063cc56def314610251576100f3565b8063193170f9116100d1578063193170f9146101615780631e2aea0614610191578063368b8772146101c15780635b34b966146101dd576100f3565b80630d16c07e146100f75780630eb47c591461011557806315bcc96614610131575b5f80fd5b6100ff610309565b60405161010c919061093b565b60405180910390f35b61012f600480360381019061012a9190610ad2565b61035f565b005b61014b60048036038101906101469190610b2c565b610382565b6040516101589190610b66565b60405180910390f35b61017b60048036038101906101769190610b2c565b6103b6565b6040516101889190610b66565b60405180910390f35b6101ab60048036038101906101a69190610c43565b610437565b6040516101b89190610b66565b60405180910390f35b6101db60048036038101906101d69190610c8a565b610486565b005b6101e56104cf565b005b61020160048036038101906101fc9190610b2c565b6104e8565b60405161020e9190610d31565b60405180910390f35b61021f610589565b60405161022c9190610b66565b60405180910390f35b61024f600480360381019061024a9190610b2c565b610592565b005b61026b60048036038101906102669190610d51565b6105f2565b6040516102789190610b66565b60405180910390f35b610289610628565b6040516102969190610d31565b60405180910390f35b6102a76106b7565b6040516102b49190610e76565b60405180910390f35b6102d760048036038101906102d29190610ec0565b610742565b6040516102e49190610b66565b60405180910390f35b61030760048036038101906103029190610b2c565b610788565b005b6060600480548060200260200160405190810160405280929190818152602001828054801561035557602002820191905f5260205f20905b815481526020019060010190808311610341575b5050505050905090565b8060055f8481526020019081526020015f20908161037d91906110e5565b505050565b5f805f90505f5b838110156103ac57808261039d91906111e1565b91508080600101915050610389565b5080915050919050565b5f808210156103fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f19061125e565b60405180910390fd5b5f600190505f600190505b83811161042d578082610418919061127c565b91508080610425906112bd565b915050610405565b5080915050919050565b5f805f90505f5b835181101561047c5783818151811061045a57610459611304565b5b60200260200101518261046d91906111e1565b9150808060010191505061043e565b5080915050919050565b805f908161049491906110e5565b507f8fae638bf5c6396194a6bb16601c4035a07fa48191638ff4102f0d96f14cfefb816040516104c49190610d31565b60405180910390a150565b60015f8154809291906104e1906112bd565b9190505550565b606060055f8381526020019081526020015f20805461050690610f18565b80601f016020809104026020016040519081016040528092919081815260200182805461053290610f18565b801561057d5780601f106105545761010080835404028352916020019161057d565b820191905f5260205f20905b81548152906001019060200180831161056057829003601f168201915b50505050509050919050565b5f600154905090565b600481908060018154018082558091505060019003905f5260205f20015f90919091909150557f3564ffb2fd8f93d7b0e9d1173ffdff5ee9775d860bfe82eaca0d0dbe07c8b634816040516105e79190610b66565b60405180910390a150565b5f80600190505f5b8381101561061d57848261060e919061127c565b915080806001019150506105fa565b508091505092915050565b60605f805461063690610f18565b80601f016020809104026020016040519081016040528092919081815260200182805461066290610f18565b80156106ad5780601f10610684576101008083540402835291602001916106ad565b820191905f5260205f20905b81548152906001019060200180831161069057829003601f168201915b5050505050905090565b6060600380548060200260200160405190810160405280929190818152602001828054801561073857602002820191905f5260205f20905b815f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116106ef575b5050505050905090565b5f60025f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b8060025f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550600333908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff167f8ad64a0ac7700dd8425ab0499f107cb6e2cd1581d803c5b8c1c79dcb8190b1af826040516108709190610b66565b60405180910390a250565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f819050919050565b6108b6816108a4565b82525050565b5f6108c783836108ad565b60208301905092915050565b5f602082019050919050565b5f6108e98261087b565b6108f38185610885565b93506108fe83610895565b805f5b8381101561092e57815161091588826108bc565b9750610920836108d3565b925050600181019050610901565b5085935050505092915050565b5f6020820190508181035f83015261095381846108df565b905092915050565b5f604051905090565b5f80fd5b5f80fd5b610975816108a4565b811461097f575f80fd5b50565b5f813590506109908161096c565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6109e48261099e565b810181811067ffffffffffffffff82111715610a0357610a026109ae565b5b80604052505050565b5f610a1561095b565b9050610a2182826109db565b919050565b5f67ffffffffffffffff821115610a4057610a3f6109ae565b5b610a498261099e565b9050602081019050919050565b828183375f83830152505050565b5f610a76610a7184610a26565b610a0c565b905082815260208101848484011115610a9257610a9161099a565b5b610a9d848285610a56565b509392505050565b5f82601f830112610ab957610ab8610996565b5b8135610ac9848260208601610a64565b91505092915050565b5f8060408385031215610ae857610ae7610964565b5b5f610af585828601610982565b925050602083013567ffffffffffffffff811115610b1657610b15610968565b5b610b2285828601610aa5565b9150509250929050565b5f60208284031215610b4157610b40610964565b5b5f610b4e84828501610982565b91505092915050565b610b60816108a4565b82525050565b5f602082019050610b795f830184610b57565b92915050565b5f67ffffffffffffffff821115610b9957610b986109ae565b5b602082029050602081019050919050565b5f80fd5b5f610bc0610bbb84610b7f565b610a0c565b90508083825260208201905060208402830185811115610be357610be2610baa565b5b835b81811015610c0c5780610bf88882610982565b845260208401935050602081019050610be5565b5050509392505050565b5f82601f830112610c2a57610c29610996565b5b8135610c3a848260208601610bae565b91505092915050565b5f60208284031215610c5857610c57610964565b5b5f82013567ffffffffffffffff811115610c7557610c74610968565b5b610c8184828501610c16565b91505092915050565b5f60208284031215610c9f57610c9e610964565b5b5f82013567ffffffffffffffff811115610cbc57610cbb610968565b5b610cc884828501610aa5565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f610d0382610cd1565b610d0d8185610cdb565b9350610d1d818560208601610ceb565b610d268161099e565b840191505092915050565b5f6020820190508181035f830152610d498184610cf9565b905092915050565b5f8060408385031215610d6757610d66610964565b5b5f610d7485828601610982565b9250506020610d8585828601610982565b9150509250929050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610de182610db8565b9050919050565b610df181610dd7565b82525050565b5f610e028383610de8565b60208301905092915050565b5f602082019050919050565b5f610e2482610d8f565b610e2e8185610d99565b9350610e3983610da9565b805f5b83811015610e69578151610e508882610df7565b9750610e5b83610e0e565b925050600181019050610e3c565b5085935050505092915050565b5f6020820190508181035f830152610e8e8184610e1a565b905092915050565b610e9f81610dd7565b8114610ea9575f80fd5b50565b5f81359050610eba81610e96565b92915050565b5f60208284031215610ed557610ed4610964565b5b5f610ee284828501610eac565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610f2f57607f821691505b602082108103610f4257610f41610eeb565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302610fa47fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610f69565b610fae8683610f69565b95508019841693508086168417925050509392505050565b5f819050919050565b5f610fe9610fe4610fdf846108a4565b610fc6565b6108a4565b9050919050565b5f819050919050565b61100283610fcf565b61101661100e82610ff0565b848454610f75565b825550505050565b5f90565b61102a61101e565b611035818484610ff9565b505050565b5b818110156110585761104d5f82611022565b60018101905061103b565b5050565b601f82111561109d5761106e81610f48565b61107784610f5a565b81016020851015611086578190505b61109a61109285610f5a565b83018261103a565b50505b505050565b5f82821c905092915050565b5f6110bd5f19846008026110a2565b1980831691505092915050565b5f6110d583836110ae565b9150826002028217905092915050565b6110ee82610cd1565b67ffffffffffffffff811115611107576111066109ae565b5b6111118254610f18565b61111c82828561105c565b5f60209050601f83116001811461114d575f841561113b578287015190505b61114585826110ca565b8655506111ac565b601f19841661115b86610f48565b5f5b828110156111825784890151825560018201915060208501945060208101905061115d565b8683101561119f578489015161119b601f8916826110ae565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6111eb826108a4565b91506111f6836108a4565b925082820190508082111561120e5761120d6111b4565b5b92915050565b7f4e756d626572206d757374206265206e6f6e2d6e6567617469766500000000005f82015250565b5f611248601b83610cdb565b915061125382611214565b602082019050919050565b5f6020820190508181035f8301526112758161123c565b9050919050565b5f611286826108a4565b9150611291836108a4565b925082820261129f816108a4565b915082820484148315176112b6576112b56111b4565b5b5092915050565b5f6112c7826108a4565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036112f9576112f86111b4565b5b600182019050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffdfea2646970667358221220e5fc3a07831936d3b6e4bfebd7d653f5bed8d24ae30df912b0145f8cd28d76cd64736f6c634300081a0033"

func TestIntegrationEthereumFlowCanCreateLargeContract(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(*aliasAccountId, NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("B71B00")        // 12mil
	gasLimitBytes, err := hex.DecodeString("B71B00") // 12mil
	contractBytes, err := hex.DecodeString("00")
	value, err := hex.DecodeString("00")
	callDataBytes, err := hex.DecodeString(LARGE_SMART_CONTRACT_BYTECODE)
	require.NoError(t, err)

	objectsList := &RLPItem{}
	objectsList.AssignList()
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(chainId))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(nonce))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(maxPriorityGas))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(maxGas))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(gasLimitBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(contractBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(value))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(callDataBytes))
	objectsList.PushBack(NewRLPItem(LIST_TYPE))

	messageBytes, err := objectsList.Write()
	require.NoError(t, err)
	messageBytes = append([]byte{0x02}, messageBytes...)

	sig := ecdsaPrivateKey.Sign(messageBytes)

	v := sig[0]
	r := sig[1:33]
	s := sig[33:65]
	vInt := int(v)

	// The compact sig recovery code is the value 27 + public key recovery code + 4
	recId := vInt - 27 - 4
	recIdBytes := []byte{byte(recId)}

	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(recIdBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(r))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(s))

	messageBytes, err = objectsList.Write()
	require.NoError(t, err)
	messageBytes = append([]byte{0x02}, messageBytes...)

	response, err := NewEthereumFlow().
		SetEthereumDataBytes(messageBytes).
		SetMaxGasAllowance(HbarFrom(10, HbarUnits.Hbar)).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := response.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, record.CallResult.SignerNonce, int64(1))
}
