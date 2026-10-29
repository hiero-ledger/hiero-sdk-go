package param

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"testing"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestFillOutTransactionHighVolume(t *testing.T) {
	trueValue := true
	falseValue := false

	tests := []struct {
		name       string
		highVolume *bool
		expected   bool
	}{
		{name: "true sets the flag", highVolume: &trueValue, expected: true},
		{name: "false leaves the flag unset", highVolume: &falseValue, expected: false},
		{name: "missing leaves the flag unset", highVolume: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := hiero.NewAccountCreateTransaction()
			common := &CommonTransactionParams{HighVolume: tt.highVolume}

			if err := common.FillOutTransaction(tx, nil); err != nil {
				t.Fatalf("FillOutTransaction returned an error: %v", err)
			}
			if got := tx.GetHighVolume(); got != tt.expected {
				t.Fatalf("GetHighVolume() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCommonTransactionParamsHighVolumeJSON(t *testing.T) {
	var present CommonTransactionParams
	if err := json.Unmarshal([]byte(`{"highVolume": true}`), &present); err != nil {
		t.Fatal(err)
	}
	if present.HighVolume == nil || !*present.HighVolume {
		t.Fatalf("expected HighVolume to be a pointer to true, got %v", present.HighVolume)
	}

	var missing CommonTransactionParams
	if err := json.Unmarshal([]byte(`{}`), &missing); err != nil {
		t.Fatal(err)
	}
	if missing.HighVolume != nil {
		t.Fatalf("expected HighVolume to be nil when absent, got %v", *missing.HighVolume)
	}
}
