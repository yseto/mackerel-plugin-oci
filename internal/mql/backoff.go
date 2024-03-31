package mql

import (
	"context"

	"github.com/cenkalti/backoff/v4"
)

var b = backoff.NewExponentialBackOff()

func (h *Handler) QueryWithBackoffRetry(ctx context.Context, input QueryInput) (*QueryResult, error) {
	var result *QueryResult
	operation := func() error {
		var err error
		result, err = h.Query(ctx, input)
		return err
	}

	err := backoff.Retry(operation, b)
	if err != nil {
		return nil, err
	}
	return result, nil
}
