package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxLenPostDecription = 200
)

type Post struct {
	id              uuid.UUID
	description     string
	timePublication time.Time
}

var (
	ErrPostValidate = errors.New("model post validate error")
)

func NewPost(id uuid.UUID, description string, timePublication time.Time) (*Post, error) {
	p := Post{
		id:              id,
		description:     strings.TrimSpace(description),
		timePublication: timePublication,
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Post) validate() error {
	if p.timePublication.IsZero() {
		return fmt.Errorf("%w: timePublication", ErrPostValidate)
	} else if len(p.description) > MaxLenPostDecription {
		return fmt.Errorf("%w: description", ErrPostValidate)
	}
	return nil
}

func (p *Post) GetID() uuid.UUID {
	return p.id
}

func (p *Post) GetDescription() string {
	return p.description
}

func (p *Post) GetTimePublication() time.Time {
	return p.timePublication
}
