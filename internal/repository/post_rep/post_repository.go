package postrep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type PostRep interface {
	GetByFilter(ctx context.Context, filterOps *models.PostFilter) ([]*models.Post, error)
	GetByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	Add(ctx context.Context, m *models.Post) error
	Delete(ctx context.Context, postID uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrPostNotFound    = errors.New("the Post was not found")
	ErrFailedToAddUser = errors.New("failed to add the Post")
	ErrUpdatePost      = errors.New("failed to update the Post in the repository")
)
