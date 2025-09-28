package models

import (
	"errors"
	"fmt"
	"strings"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenCategoryTitle       = 50
	MaxLenPCategoryDecription = 200
)

type Category struct {
	id          uuid.UUID
	title       string
	description string
}

var (
	ErrCategoryValidate = errors.New("model category validate error")
)

func NewCategory(id uuid.UUID, title string, description string) (*Category, error) {
	p := Category{
		id:          id,
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Category) validate() error {
	if p.title == "" || len(p.title) > MaxLenProductTitle {
		return fmt.Errorf("%w: title", ErrCategoryValidate)
	} else if len(p.description) > MaxLenProductDecription {
		return fmt.Errorf("%w: description", ErrCategoryValidate)
	}
	return nil
}

func (p *Category) ToResponse() reqresp.CategoryResponse {
	return reqresp.CategoryResponse{
		ID:          p.id.String(),
		Title:       p.title,
		Description: p.description,
	}
}

func (p *Category) GetID() uuid.UUID {
	return p.id
}

func (p *Category) GetDescription() string {
	return p.description
}
