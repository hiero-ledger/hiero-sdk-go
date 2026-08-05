//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestUnitExecuteRequestCancelsContext(t *testing.T) {
	t.Parallel()

	var requestContext context.Context
	method := _Method{
		query: func(ctx context.Context, _ *services.Query, _ ...grpc.CallOption) (*services.Response, error) {
			requestContext = ctx
			return &services.Response{}, nil
		},
	}

	_, _, err := _ExecuteRequest(method, time.Minute, &services.Query{})
	require.NoError(t, err)
	require.ErrorIs(t, requestContext.Err(), context.Canceled)
}
