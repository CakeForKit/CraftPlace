package userrep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type UserRep interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Add(ctx context.Context, m *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error)
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrUserNotFound       = errors.New("the User was not found")
	ErrFailedToAddUser    = errors.New("failed to add the User")
	ErrDuplicateLoginUser = errors.New("a user with this login already exists")
	ErrUpdateUser         = errors.New("failed to update the User in the repository")
	ErrReadOnly           = errors.New("it is read only connection")
)
