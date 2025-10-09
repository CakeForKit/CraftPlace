package models

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenProductTitle      = 50
	MaxLenProductDecription = 200
)

type Product struct {
	id          uuid.UUID
	title       string
	description string
	cost        uint64
	shopID      uuid.UUID
	categoryIDs uuid.UUIDs
	updateTime  time.Time
}

var (
	ErrProductValidate = errors.New("model Product validate error")
)

func NewProduct(
	id uuid.UUID, title string,
	description string, cost uint64,
	shopID uuid.UUID, categoryIDs uuid.UUIDs,
	updateTime time.Time,
) (*Product, error) {
	p := Product{
		id:          id,
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
		cost:        cost,
		shopID:      shopID,
		categoryIDs: categoryIDs,
		updateTime:  updateTime,
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Product) validate() error {
	if p.title == "" || len(p.title) > MaxLenProductTitle {
		return fmt.Errorf("%w: title", ErrProductValidate)
	} else if len(p.description) > MaxLenProductDecription {
		return fmt.Errorf("%w: description", ErrProductValidate)
	} else if p.shopID == uuid.Nil {
		return fmt.Errorf("%w: shopID", ErrProductValidate)
	}
	return nil
}

func (p *Product) ToResponse() reqresp.ProductResponse {
	return reqresp.ProductResponse{
		ID:          p.id.String(),
		Title:       p.title,
		Description: p.description,
		Cost:        p.cost,
		ShopID:      p.shopID,
		CategoryIDs: p.categoryIDs,
		UpdateTime:  p.updateTime,
	}
}

func (p *Product) AddCategoryIDs(cids uuid.UUIDs) error {
	for _, oldID := range p.categoryIDs {
		if slices.Contains(cids, oldID) {
			return fmt.Errorf("Event.AddArtworks %w", ErrCategoryValidate)
		}
	}
	p.categoryIDs = append(p.categoryIDs, cids...)
	return nil
}

func (p *Product) GetID() uuid.UUID {
	return p.id
}

func (p *Product) GetTitle() string {
	return p.title
}

func (p *Product) GetDescription() string {
	return p.description
}

func (p *Product) GetCost() uint64 {
	return p.cost
}

func (p *Product) GetShopID() uuid.UUID {
	return p.shopID
}

func (p *Product) GetUpdateTime() time.Time {
	return p.updateTime
}
