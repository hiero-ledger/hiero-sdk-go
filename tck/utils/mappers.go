package utils

import "fmt"

// SPDX-License-Identifier: Apache-2.0

// Helper function to map serial numbers to strings
func MapSerialNumbersToString(serials []int64) []string {
	serialStrings := make([]string, len(serials))
	for i, serial := range serials {
		serialStrings[i] = fmt.Sprintf("%d", serial)
	}
	return serialStrings
}
