package postservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	postrep "github.com/CakeForKit/CraftPlace.git/internal/repository/post_rep"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	"github.com/google/uuid"
)

type PostServ interface {
	// GetPostsByFilter(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error)
	Add(ctx context.Context, addReq AddPostData) (*models.Post, error)
	Delete(ctx context.Context, postID uuid.UUID) error
}

var (
	ErrPostServ = errors.New("PostServ")
)

type AddPostData struct {
	Description string    `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	ShopID      uuid.UUID `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type postServ struct {
	postRep postrep.PostRep
	authz   auth.AuthZ
	shopRep shoprep.ShopRep
}

func NewPostServ(postRep postrep.PostRep, authz auth.AuthZ, shopRep shoprep.ShopRep) PostServ {
	return &postServ{
		postRep: postRep,
		authz:   authz,
		shopRep: shopRep,
	}
}

func (s *postServ) checkUserRights(ctx context.Context, shopID uuid.UUID) error {
	baseErr := fmt.Errorf("checkUserRights")
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	shop, err := s.shopRep.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if shop.GetUserID() != userID {
		return fmt.Errorf("%w: %w", baseErr, auth.ErrHasNoRights)
	}
	return nil
}

func (s *postServ) Add(ctx context.Context, addReq AddPostData) (*models.Post, error) {
	baseErr := fmt.Errorf("%w Add", ErrPostServ)
	err := s.checkUserRights(ctx, addReq.ShopID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	newPost, err := models.NewPost(
		uuid.New(),
		addReq.Description,
		time.Now().UTC(),
		addReq.ShopID,
	)
	err = s.postRep.Add(ctx, newPost)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return newPost, nil
}

func (s *postServ) Delete(ctx context.Context, postID uuid.UUID) error {
	baseErr := fmt.Errorf("%w Add", ErrPostServ)
	post, err := s.postRep.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.checkUserRights(ctx, post.GetShopID())
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.postRep.Delete(ctx, postID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}

// func (s *postServ) GetPostsByFilter(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error) {
// 	baseErr := fmt.Errorf("%w GetPostsByFilter", ErrPostServ)
// 	posts, err := s.postRep.GetByFilter(ctx, filterOps)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %w", baseErr, err)
// 	}
// 	return posts, nil
// }
