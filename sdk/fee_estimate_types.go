package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
)

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

// toProto converts SDK FeeEstimateMode to proto EstimateMode
func (m FeeEstimateMode) toProto() mirror.EstimateMode {
	switch m {
	case FeeEstimateModeState:
		return mirror.EstimateMode_STATE
	case FeeEstimateModeIntrinsic:
		return mirror.EstimateMode_INTRINSIC
	default:
		return mirror.EstimateMode_STATE
	}
}

// feeEstimateModeFromProto converts proto EstimateMode to SDK FeeEstimateMode
func feeEstimateModeFromProto(m mirror.EstimateMode) FeeEstimateMode {
	switch m {
	case mirror.EstimateMode_STATE:
		return FeeEstimateModeState
	case mirror.EstimateMode_INTRINSIC:
		return FeeEstimateModeIntrinsic
	default:
		return FeeEstimateModeState
	}
}

// FeeExtra represents an extra fee charged for the transaction
type FeeExtra struct {
	Name       string // The unique name of this extra fee as defined in the fee schedule
	Included   uint32 // The count of this "extra" that is included for free
	Count      uint32 // The actual count of items received
	Charged    uint32 // The charged count of items as calculated by max(0, count - included)
	FeePerUnit uint64 // The fee price per unit in tinycents
	Subtotal   uint64 // The subtotal in tinycents for this extra fee
}

// feeExtraFromProto converts proto FeeExtra to SDK FeeExtra
func feeExtraFromProto(pb *mirror.FeeExtra) FeeExtra {
	if pb == nil {
		return FeeExtra{}
	}
	return FeeExtra{
		Name:       pb.GetName(),
		Included:   pb.GetIncluded(),
		Count:      pb.GetCount(),
		Charged:    pb.GetCharged(),
		FeePerUnit: pb.GetFeePerUnit(),
		Subtotal:   pb.GetSubtotal(),
	}
}

// FeeEstimate represents the fee estimate for a component
type FeeEstimate struct {
	Base   uint64     // The base fee price, in tinycents
	Extras []FeeExtra // The extra fees that apply for this fee component
}

// Subtotal returns the total subtotal for this fee estimate (base + sum of all extras)
func (fe *FeeEstimate) Subtotal() uint64 {
	total := fe.Base
	for _, extra := range fe.Extras {
		total += extra.Subtotal
	}
	return total
}

// feeEstimateFromProto converts proto FeeEstimate to SDK FeeEstimate
func feeEstimateFromProto(pb *mirror.FeeEstimate) FeeEstimate {
	if pb == nil {
		return FeeEstimate{}
	}
	extras := make([]FeeExtra, 0, len(pb.GetExtras()))
	for _, extra := range pb.GetExtras() {
		extras = append(extras, feeExtraFromProto(extra))
	}
	return FeeEstimate{
		Base:   pb.GetBase(),
		Extras: extras,
	}
}

// NetworkFee represents the network fee component
type NetworkFee struct {
	Multiplier uint32 // Multiplied by the node fee to determine the total network fee
	Subtotal   uint64 // The subtotal in tinycents for the network fee component
}

// networkFeeFromProto converts proto NetworkFee to SDK NetworkFee
func networkFeeFromProto(pb *mirror.NetworkFee) NetworkFee {
	if pb == nil {
		return NetworkFee{}
	}
	return NetworkFee{
		Multiplier: pb.GetMultiplier(),
		Subtotal:   pb.GetSubtotal(),
	}
}

// FeeEstimateResponse represents the response containing the estimated transaction fees
type FeeEstimateResponse struct {
	Mode       FeeEstimateMode // The mode that was used to calculate the fees
	NetworkFee NetworkFee      // The network fee component
	NodeFee    FeeEstimate     // The node fee component
	ServiceFee FeeEstimate     // The service fee component
	Notes      []string        // An array of strings for any caveats
	Total      uint64          // The sum of the network, node, and service subtotals in tinycents
}

// feeEstimateResponseFromProto converts proto FeeEstimateResponse to SDK FeeEstimateResponse
func feeEstimateResponseFromProto(pb *mirror.FeeEstimateResponse) FeeEstimateResponse {
	if pb == nil {
		return FeeEstimateResponse{}
	}
	notes := make([]string, len(pb.GetNotes()))
	copy(notes, pb.GetNotes())
	return FeeEstimateResponse{
		Mode:       feeEstimateModeFromProto(pb.GetMode()),
		NetworkFee: networkFeeFromProto(pb.GetNetwork()),
		NodeFee:    feeEstimateFromProto(pb.GetNode()),
		ServiceFee: feeEstimateFromProto(pb.GetService()),
		Notes:      notes,
		Total:      pb.GetTotal(),
	}
}
