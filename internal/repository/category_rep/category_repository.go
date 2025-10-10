package categoryrep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type CategoryRep interface {
	GetByFilter(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error)
	Add(ctx context.Context, m *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrCategoryNotFound = errors.New("the Category was not found")
	ErrFailedToAddUser  = errors.New("failed to add the Category")
	ErrUpdateCategory   = errors.New("failed to update the Category in the repository")
)
