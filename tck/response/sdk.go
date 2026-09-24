package response

// SPDX-License-Identifier: Apache-2.0

type SetupResponse struct {
	Message string
	Status  string
}

type PingResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
