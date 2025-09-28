package models

import (
	"errors"
	"fmt"
	"strings"

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
}

var (
	ErrShopValidate = errors.New("model Shop validate error")
)

func NewShop(id uuid.UUID, title string, description string) (*Shop, error) {
	s := Shop{
		id:          id,
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
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
	}
	return nil
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
