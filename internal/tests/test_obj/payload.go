package testobj

import (
	"time"

	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/google/uuid"
)

type PayloadMother interface {
	UserPayload(userID uuid.UUID) tokenmaker.Payload
}

func NewPayloadMother() PayloadMother {
	return &payloadMother{}
}

type payloadMother struct{}

func (pm *payloadMother) UserPayload(userID uuid.UUID) tokenmaker.Payload {
	return tokenmaker.Payload{
		PersonID:  userID,
		Role:      tokenmaker.UserRole,
		ExpiredAt: time.Now().Add(time.Hour),
	}
}
