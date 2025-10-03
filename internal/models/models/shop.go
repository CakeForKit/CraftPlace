package models

import (
	"errors"
	"fmt"
	"strings"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenShopTitle      = 50
	MaxLenShopDecription = 200
)

type Shop struct {
	id          uuid.UUID
	title       string
	description string
	userID      uuid.UUID
}

var (
	ErrShopValidate = errors.New("model Shop validate error")
)

func NewShop(id uuid.UUID, title string, description string, userID uuid.UUID) (*Shop, error) {
	s := Shop{
		id:          id,
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
		userID:      userID,
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Shop) validate() error {
	if s.title == "" || len(s.title) > MaxLenShopTitle {
		return fmt.Errorf("%w: title", ErrShopValidate)
	} else if len(s.description) > MaxLenShopDecription {
		return fmt.Errorf("%w: description", ErrShopValidate)
	} else if s.userID == uuid.Nil {
		return fmt.Errorf("%w: userID", ErrShopValidate)
	}
	return nil
}

func (p *Shop) ToResponse() reqresp.ShopResponse {
	return reqresp.ShopResponse{
		ShopID:      p.id.String(),
		Title:       p.title,
		Description: p.description,
		UserID:      p.userID,
	}
}

func (s *Shop) GetID() uuid.UUID {
	return s.id
}

func (s *Shop) GetTitle() string {
	return s.title
}

func (s *Shop) GetDescription() string {
	return s.description
}

func (s *Shop) GetUserID() uuid.UUID {
	return s.userID
}
