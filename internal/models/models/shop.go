package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenShopTitle       = 255
	MaxLenShopDescription = 500
)

type Shop struct {
	id          uuid.UUID
	title       string
	description string
	userID      uuid.UUID
	updateTime  time.Time
}

var (
	ErrShopValidate = errors.New("model Shop validate error")
)

func NewShop(id uuid.UUID, title string, description string, userID uuid.UUID, updateTime time.Time) (*Shop, error) {
	s := Shop{
		id:          id,
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
		userID:      userID,
		updateTime:  updateTime,
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Shop) validate() error {
	if s.title == "" || len(s.title) > MaxLenShopTitle {
		return fmt.Errorf("%w: title", ErrShopValidate)
	} else if len(s.description) > MaxLenShopDescription {
		return fmt.Errorf("%w: description", ErrShopValidate)
	} else if s.userID == uuid.Nil {
		return fmt.Errorf("%w: userID", ErrShopValidate)
	}
	return nil
}

func (s *Shop) ToResponse() reqresp.ShopResponse {
	return reqresp.ShopResponse{
		ShopID:      s.id.String(),
		Title:       s.title,
		Description: s.description,
		UserID:      s.userID,
		UpdateTime:  s.updateTime,
	}
}

func (s *Shop) Update(updateReq *reqresp.UpdateShopRequest) error {
	copyS := *s
	if updateReq.Title != "" {
		copyS.title = updateReq.Title
	}
	if updateReq.Description != "" {
		copyS.description = updateReq.Description
	}
	if err := copyS.validate(); err != nil {
		return err
	}
	*s = copyS
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

func (s *Shop) GetUserID() uuid.UUID {
	return s.userID
}

func (s *Shop) GetUpdateTime() time.Time {
	return s.updateTime
}
