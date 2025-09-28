package auth

import (
	"context"
	"errors"

	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/google/uuid"
)

type authZContextKey int

const (
	AuthZContextKey authZContextKey = iota
)

var (
	ErrNotAuthZ    = errors.New("not authorized")
	ErrHasNoRights = errors.New("has no rights")
)

type AuthZ interface {
	Authorize(ctx context.Context, payload tokenmaker.Payload) context.Context
	UserIDFromContext(ctx context.Context) (uuid.UUID, error)
}

func NewAuthZ() (AuthZ, error) {
	return &authZ{}, nil
}

type authZ struct {
}

func (a *authZ) Authorize(ctx context.Context, payload tokenmaker.Payload) context.Context {
	return context.WithValue(ctx, AuthZContextKey, payload)
}

func (a *authZ) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	payload, ok := ctx.Value(AuthZContextKey).(tokenmaker.Payload)
	if !ok {
		return uuid.Nil, ErrNotAuthZ
	}
	if payload.Role != tokenmaker.UserRole {
		return uuid.Nil, ErrHasNoRights
	}
	return payload.PersonID, nil
}
