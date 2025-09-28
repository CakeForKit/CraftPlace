package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Category struct {
	id          uuid.UUID
	description string
}

var (
	ErrCategoryValidate = errors.New("model category validate error")
)

func NewCategory(id uuid.UUID, description string) (*Category, error) {
	p := Category{
		id:          id,
		description: strings.TrimSpace(description),
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Category) validate() error {
	return nil
}

func (p *Category) GetID() uuid.UUID {
	return p.id
}

func (p *Category) GetDescription() string {
	return p.description
}
