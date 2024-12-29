package cache

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) EvictAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockService) Evict(ctx context.Context, key string) (value interface{}, err error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (m *MockService) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockService) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Get(1).(time.Time), args.Error(3)
}

func (m *MockService) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	args := m.Called(ctx)

	if keysResult, ok := args.Get(0).([]string); ok {
		keys = keysResult
	} else {
		keys = nil
	}

	if valuesResult, ok := args.Get(1).([]interface{}); ok {
		values = valuesResult
	} else {
		values = nil
	}

	err = args.Error(2)
	return keys, values, err
}
