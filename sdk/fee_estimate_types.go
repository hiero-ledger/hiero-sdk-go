package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"

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
	Name       string `json:"name"`       // The unique name of this extra fee as defined in the fee schedule
	Included   uint32 `json:"included"`   // The count of this "extra" that is included for free
	Count      uint32 `json:"count"`      // The actual count of items received
	Charged    uint32 `json:"charged"`    // The charged count of items as calculated by max(0, count - included)
	FeePerUnit uint64 `json:"feePerUnit"` // The fee price per unit in tinycents
	Subtotal   uint64 `json:"subtotal"`   // The subtotal in tinycents for this extra fee
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
	Base   uint64     `json:"base"`             // The base fee price, in tinycents
	Extras []FeeExtra `json:"extras,omitempty"` // The extra fees that apply for this fee component
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
	Multiplier uint32 `json:"multiplier"` // Multiplied by the node fee to determine the total network fee
	Subtotal   uint64 `json:"subtotal"`   // The subtotal in tinycents for the network fee component
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
	Mode       FeeEstimateMode `json:"mode"`            // The mode that was used to calculate the fees
	NetworkFee NetworkFee      `json:"network"`         // The network fee component
	NodeFee    FeeEstimate     `json:"node"`            // The node fee component
	ServiceFee FeeEstimate     `json:"service"`         // The service fee component
	Notes      []string        `json:"notes,omitempty"` // An array of strings for any caveats
	Total      uint64          `json:"total"`           // The sum of the network, node, and service subtotals in tinycents
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

// feeEstimateModeFromString converts string mode to SDK FeeEstimateMode
func feeEstimateModeFromString(mode string) FeeEstimateMode {
	switch mode {
	case "STATE":
		return FeeEstimateModeState
	case "INTRINSIC":
		return FeeEstimateModeIntrinsic
	default:
		return FeeEstimateModeState
	}
}

// feeEstimateResponseFromREST converts REST API JSON response to SDK FeeEstimateResponse
// The mode field comes as a string in JSON, so we need to handle it specially
func feeEstimateResponseFromREST(data []byte) (FeeEstimateResponse, error) {
	// Temporary struct to handle mode as string during unmarshaling
	type tempResponse struct {
		Mode       string      `json:"mode"`
		NetworkFee NetworkFee  `json:"network"`
		NodeFee    FeeEstimate `json:"node"`
		ServiceFee FeeEstimate `json:"service"`
		Notes      []string    `json:"notes,omitempty"`
		Total      uint64      `json:"total"`
	}

	var temp tempResponse
	if err := json.Unmarshal(data, &temp); err != nil {
		return FeeEstimateResponse{}, err
	}

	notes := make([]string, len(temp.Notes))
	copy(notes, temp.Notes)

	return FeeEstimateResponse{
		Mode:       feeEstimateModeFromString(temp.Mode),
		NetworkFee: temp.NetworkFee,
		NodeFee:    temp.NodeFee,
		ServiceFee: temp.ServiceFee,
		Notes:      notes,
		Total:      temp.Total,
	}, nil
}
