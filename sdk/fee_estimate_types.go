package hiero

// SPDX-License-Identifier: Apache-2.0

// FeeEstimateMode represents the mode of fee estimation
type FeeEstimateMode int32

const (
	// FeeEstimateModeIntrinsic estimates from the payload alone, ignoring state-dependent costs (default)
	FeeEstimateModeIntrinsic FeeEstimateMode = iota
	// FeeEstimateModeState uses the mirror node's latest known state for fee estimation
	FeeEstimateModeState
)

// String returns the string representation of FeeEstimateMode
func (m FeeEstimateMode) String() string {
	switch m {
	case FeeEstimateModeIntrinsic:
		return "INTRINSIC"
	case FeeEstimateModeState:
		return "STATE"
	default:
		return "UNKNOWN"
	}
}

// FeeExtra represents an extra fee charged for the transaction
type FeeExtra struct {
	Name       string `json:"name"`         // The unique name of this extra fee as defined in the fee schedule
	Included   uint64 `json:"included"`     // The count of this "extra" that is included for free
	Count      uint64 `json:"count"`        // The actual count of items received
	Charged    uint64 `json:"charged"`      // The charged count of items as calculated by max(0, count - included)
	FeePerUnit uint64 `json:"fee_per_unit"` // The fee price per unit in tinycents
	Subtotal   uint64 `json:"subtotal"`     // The subtotal in tinycents for this extra fee
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
	// The high-volume pricing multiplier per HIP-1313. A value of 1 indicates no
	// high-volume pricing. A value greater than 1 applies when the transaction's
	// highVolume flag is true and throttle utilization is non-zero.
	HighVolumeMultiplier uint64 `json:"high_volume_multiplier"`
	// NetworkFee covers gossip, consensus, signature verification, fee payment, and storage.
	NetworkFee NetworkFee `json:"network"`
	// NodeFee is paid to the submitting node for pre-checking the transaction (required).
	NodeFee FeeEstimate `json:"node"`
	// ServiceFee covers execution, Merkle state, and blockchain storage costs.
	ServiceFee FeeEstimate `json:"service"`
	// Total is the sum of network, node, and service subtotals in tinycents.
	Total uint64 `json:"total"`
}
