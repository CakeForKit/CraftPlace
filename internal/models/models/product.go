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

type AddProductData struct {
	Title       string      `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"100"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

type UpdateProductData struct {
	Title       string      `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"200"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

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

func (p *Product) Update(updateReq *UpdateProductData) error {
	copyP := *p
	if updateReq.Title != "" {
		copyP.title = updateReq.Title
	}
	if updateReq.Description != "" {
		copyP.description = updateReq.Description
	}
	copyP.cost = updateReq.Cost
	copyP.updateTime = time.Now().UTC()
	if err := copyP.validate(); err != nil {
		return err
	}
	*p = copyP
	return nil
}

func (p *Product) AddCategoryIDs(cids uuid.UUIDs) error {
	for _, oldID := range p.categoryIDs {
		if slices.Contains(cids, oldID) {
			return fmt.Errorf("Event.AddCategoryIDs %w", ErrCategoryValidate)
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
