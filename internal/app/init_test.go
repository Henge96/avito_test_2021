package app_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"packs/internal/app"
	"packs/lib/test_helper"
	"testing"
)

type mocks struct {
	repo     *MockRepo
	exchange *MockExchangeClient
}

func start(t *testing.T) (context.Context, *app.Core, *mocks, *require.Assertions) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := NewMockRepo(ctrl)
	mockExchange := NewMockExchangeClient(ctrl)

	module := app.New(mockRepo, mockExchange)
	mocks := &mocks{
		repo:     mockRepo,
		exchange: mockExchange,
	}

	return test_helper.Context(t), module, mocks, require.New(t)
}
