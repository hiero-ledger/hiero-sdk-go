package hiero

// SPDX-License-Identifier: Apache-2.0

// FeeEstimateMode represents the mode of fee estimation
type FeeEstimateMode int32

const (
	// FeeEstimateModeState uses the latest known state for fee estimation
	FeeEstimateModeState FeeEstimateMode = iota
	// FeeEstimateModeIntrinsic estimates from the payload alone, ignoring state-dependent costs
	FeeEstimateModeIntrinsic
)

// String returns the string representation of FeeEstimateMode
func (m FeeEstimateMode) String() string {
	switch m {
	case FeeEstimateModeState:
		return "STATE"
	case FeeEstimateModeIntrinsic:
		return "INTRINSIC"
	default:
		return "UNKNOWN"
	}
}

// FeeExtra represents an extra fee charged for the transaction
type FeeExtra struct {
	Name       string `json:"name"`       // The unique name of this extra fee as defined in the fee schedule
	Included   uint32 `json:"included"`   // The count of this "extra" that is included for free
	Count      uint32 `json:"count"`      // The actual count of items received
	Charged    uint32 `json:"charged"`    // The charged count of items as calculated by max(0, count - included)
	FeePerUnit uint64 `json:"feePerUnit"` // The fee price per unit in tinycents
	Subtotal   uint64 `json:"subtotal"`   // The subtotal in tinycents for this extra fee
}

// FeeEstimate represents the fee estimate for a component
type FeeEstimate struct {
	Base   uint64     `json:"base"`             // The base fee price, in tinycents (required)
	Extras []FeeExtra `json:"extras,omitempty"` // The extra fees that apply for this fee component (optional, can be empty array)
}

// Subtotal returns the total subtotal for this fee estimate (base + sum of all extras)
func (fe *FeeEstimate) Subtotal() uint64 {
	total := fe.Base
	for _, extra := range fe.Extras {
		total += extra.Subtotal
	}
	return total
}

// NetworkFee represents the network fee component which covers the cost of gossip, consensus,
// signature verifications, fee payment, and storage.
type NetworkFee struct {
	Multiplier uint32 `json:"multiplier"` // Multiplied by the node fee to determine the total network fee
	Subtotal   uint64 `json:"subtotal"`   // The subtotal in tinycents for the network fee component
}

// FeeEstimateResponse represents the response containing the estimated transaction fees
type FeeEstimateResponse struct {
	Mode       FeeEstimateMode `json:"mode"`    // The mode that was used to calculate the fees
	NetworkFee NetworkFee      `json:"network"` // The network fee component
	NodeFee    FeeEstimate     `json:"node"`    // The node fee component which is to be paid to the node that submitted the transaction to the network. This fee exists to compensate the node for the work it performed to pre-check the transaction before submitting it, and incentivizes the node to accept new transactions from users (required)
	ServiceFee FeeEstimate     `json:"service"` // The service fee component which covers execution costs, state saved in the Merkle tree, and additional costs to the blockchain storage
	Notes      []string        `json:"notes"`   // An array of strings for any caveats
	Total      uint64          `json:"total"`   // The sum of the network, node, and service subtotals in tinycents
}
