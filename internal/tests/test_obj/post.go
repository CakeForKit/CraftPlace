package testobj

import (
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type PostMother interface {
	PostP() *models.Post
}

func NewPostMother() PostMother {
	return &postMother{}
}

type postMother struct{}

func (um *postMother) PostP() *models.Post {
	post, _ := models.NewPost(
		uuid.New(),
		"test-desription",
		time.Now().UTC(),
		uuid.New(),
	)
	return post
}
