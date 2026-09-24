package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	// The only URL schemes a mirror node REST call can be made over.
	mirrorNodeSchemeHTTP  = "http"
	mirrorNodeSchemeHTTPS = "https"
)

// mirrorNodeRestBaseURL returns the client's mirror node REST API base URL, erroring when no
// mirror network is configured.
func mirrorNodeRestBaseURL(client *Client) (string, error) {
	if client == nil {
		return "", errNoClientProvided
	}
	// The nil check is not redundant: GetMirrorNetwork dereferences the field.
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return "", errors.New("mirror node is not set")
	}

	return client.GetMirrorRestApiBaseUrl()
}

// mirrorNodeGetJSONObject GETs path and decodes the body as a JSON object. A non-2xx body is decoded too, unless retries ran out.
func mirrorNodeGetJSONObject(client *Client, path mirrorNodeRestPath) (map[string]interface{}, error) {
	restClient, err := client.mirrorRestClient(path, client.mirrorHttpPolicy())
	if err != nil {
		return nil, err
	}

	resp, err := restClient.get(path, CancellationNone())
	if err != nil {
		if errors.Is(err, errMirrorHttpRetriesExhausted) {
			return nil, mirrorNodeStatusError(resp)
		}
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// linksJSON is the pagination block on every paginated mirror node response.
type linksJSON struct {
	Next *string `json:"next"`
}
