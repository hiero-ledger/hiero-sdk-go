//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitMirrorNodeRestBaseURLRejectsNilClient(t *testing.T) {
	t.Parallel()

	_, err := mirrorNodeRestBaseURL(nil)
	require.ErrorIs(t, err, errNoClientProvided)
}

func TestUnitMirrorNodeRestBaseURLRejectsEmptyMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{})

	_, err = mirrorNodeRestBaseURL(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}

func TestUnitMirrorNodeGetJSONObjectRetries(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable)},
		{resp: HttpResponse{statusCode: http.StatusOK, body: []byte(`{"account":"0.0.7"}`)}},
	}}
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithTransport(transport).WithRetryPolicy(fastRetryPolicy(3)))

	result, err := mirrorNodeGetJSONObject(client, testMirrorPath(t, "/accounts/0.0.7"))
	require.NoError(t, err)
	assert.Equal(t, "0.0.7", result["account"])
	assert.Equal(t, 2, transport.callCount())
}

func TestUnitMirrorNodeGetJSONObjectDecodesNon200Body(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusNotFound, body: []byte(`{"_status":{"messages":[{"message":"Not found"}]}}`)}},
	}}
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithTransport(transport))

	result, err := mirrorNodeGetJSONObject(client, testMirrorPath(t, "/accounts/0.0.7"))
	require.NoError(t, err)
	assert.Contains(t, result, "_status")
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorNodeGetJSONObjectReportsExhaustedRetries(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusServiceUnavailable)}}}
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithTransport(transport).WithRetryPolicy(fastRetryPolicy(2)))

	_, err = mirrorNodeGetJSONObject(client, testMirrorPath(t, "/accounts/0.0.7"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "received non-200 response from mirror node: 503")
}
