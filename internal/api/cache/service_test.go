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
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockService) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockService) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	args := m.Called(key)
	return args.String(0), time.Time{}, args.Error(1)
}

func (m *MockService) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	args := m.Called(ctx)
	return nil, nil, args.Error(1)
}
