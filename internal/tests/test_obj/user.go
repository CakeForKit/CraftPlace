package testobj

import (
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type UserMother interface {
	// DefaultUser(userID uuid.UUID) models.User
	UserWithPswdHash(userID uuid.UUID, hashedPassword string) models.User
	DefaultUserP(userID uuid.UUID) *models.User
	UserWithLoginP(userID uuid.UUID, login string) *models.User
}

func NewUserMother() UserMother {
	return &userMother{}
}

type userMother struct{}

func (um *userMother) UserWithPswdHash(userID uuid.UUID, hashedPassword string) models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login"+uuid.New().String(),
		hashedPassword,
	)
	return user
}

func (um *userMother) DefaultUserP(userID uuid.UUID) *models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login"+uuid.New().String(),
		"hashed-password",
	)
	return &user
}

func (um *userMother) UserWithLoginP(userID uuid.UUID, login string) *models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		login,
		"hashed-password",
	)
	return &user
}

/*
func (um *userMother) DefaultUser(userID uuid.UUID) models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login-"+uuid.New().String(),
		"hashed-password",
	)
	return user
}

*/
