package postrep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type PostRep interface {
	GetByFilter(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error)
	Add(ctx context.Context, m *models.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrPostNotFound    = errors.New("the Post was not found")
	ErrFailedToAddUser = errors.New("failed to add the Post")
	ErrUpdatePost      = errors.New("failed to update the Post in the repository")
)
