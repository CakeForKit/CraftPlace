package auth

import (
	"context"

	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthZ реализует AuthZ интерфейс для тестирования
type MockAuthZ struct {
	mock.Mock
}

func (m *MockAuthZ) Authorize(ctx context.Context, payload tokenmaker.Payload) context.Context {
	args := m.Called(ctx, payload)
	return args.Get(0).(context.Context)
}

func (m *MockAuthZ) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
