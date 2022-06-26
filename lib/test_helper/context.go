package test_helper

import (
	"context"
	"testing"
)

func Context(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60)
	t.Cleanup(cancel)

	return ctx
}
