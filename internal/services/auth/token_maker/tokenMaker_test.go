package tokenmaker_test

import (
	"testing"
	"time"

	token "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TokenMakerSuite struct {
	suite.Suite
	maker token.TokenMaker
}

func TestTokenMaker(t *testing.T) {
	suite.RunSuite(t, new(TokenMakerSuite))
}

func (s *TokenMakerSuite) BeforeEach(t provider.T) {
	t.Tag("Token Maker")

	maker, err := token.NewTokenMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	s.maker = maker
}

func (s *TokenMakerSuite) TestJWTMaker_CreateAndVerifyToken(t provider.T) {
	t.WithNewStep("Успешное создание и верификация токена", func(sCtx provider.StepCtx) {
		// Arrange
		userID := uuid.New()
		role := token.UserRole
		duration := time.Minute

		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		// Act
		tokenStr, err := s.maker.CreateToken(userID, role, duration)

		// Assert
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(tokenStr)

		// Verify
		payload, err := s.maker.VerifyToken(tokenStr, role)
		sCtx.Require().NoError(err)
		sCtx.Require().NotNil(payload)

		assert.Equal(t, userID, payload.GetPersonID())
		assert.Equal(t, role, payload.GetRole())
		assert.WithinDuration(t, issuedAt, payload.GetExpiredAt().Add(-duration), time.Second)
		assert.WithinDuration(t, expiredAt, payload.GetExpiredAt(), time.Second)
	})
}

func (s *TokenMakerSuite) TestJWTMaker_ExpiredToken(t provider.T) {
	t.WithNewStep("Токен с истекшим сроком действия", func(sCtx provider.StepCtx) {
		// Arrange
		userID := uuid.New()
		role := token.UserRole
		negativeDuration := -time.Minute

		// Act
		tokenStr, err := s.maker.CreateToken(userID, role, negativeDuration)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(tokenStr)

		// Assert
		payload, err := s.maker.VerifyToken(tokenStr, role)
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, token.ErrExpiredToken)
		sCtx.Require().Nil(payload)
	})
}

func (s *TokenMakerSuite) TestJWTMaker_InvalidTokenAlgNone(t provider.T) {
	t.WithNewStep("Токен с неверным алгоритмом подписи", func(sCtx provider.StepCtx) {
		// Arrange
		userID := uuid.New()
		role := token.UserRole

		payload, err := token.NewPayload(userID, role, time.Minute)
		sCtx.Require().NoError(err)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
		tokenStr, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		sCtx.Require().NoError(err)

		// Act
		resultPayload, err := s.maker.VerifyToken(tokenStr, role)

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, token.ErrInvalidToken)
		sCtx.Require().Nil(resultPayload)
	})
}

func (s *TokenMakerSuite) TestJWTMaker_IncorrectRole(t provider.T) {
	t.WithNewStep("Токен с неверной ролью", func(sCtx provider.StepCtx) {
		// Arrange
		userID := uuid.New()
		correctRole := token.UserRole
		incorrectRole := token.RoleAuth("admin_role")
		duration := time.Minute

		// Act
		tokenStr, err := s.maker.CreateToken(userID, correctRole, duration)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(tokenStr)

		// Assert
		payload, err := s.maker.VerifyToken(tokenStr, incorrectRole)
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, token.ErrIncorrectRole)
		sCtx.Require().Nil(payload)
	})
}

func (s *TokenMakerSuite) TestJWTMaker_InvalidSecretKey(t provider.T) {
	t.WithNewStep("Создание maker с неверным секретным ключом", func(sCtx provider.StepCtx) {
		// Arrange & Act
		maker, err := token.NewTokenMaker("short_key")

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().Nil(maker)
		assert.Contains(t, err.Error(), "invalid key size")
	})
}

func (s *TokenMakerSuite) TestJWTMaker_ValidWithDifferentRoles(t provider.T) {
	t.WithNewStep("Успешная верификация с правильной ролью", func(sCtx provider.StepCtx) {
		// Arrange
		userID := uuid.New()
		role := token.UserRole
		duration := time.Minute

		// Act
		tokenStr, err := s.maker.CreateToken(userID, role, duration)
		sCtx.Require().NoError(err)

		// Assert
		payload, err := s.maker.VerifyToken(tokenStr, role)
		sCtx.Require().NoError(err)
		sCtx.Require().NotNil(payload)
		assert.Equal(t, userID, payload.GetPersonID())
		assert.Equal(t, role, payload.GetRole())
	})
}

func (s *TokenMakerSuite) TestJWTMaker_EmptyToken(t provider.T) {
	t.WithNewStep("Верификация пустого токена", func(sCtx provider.StepCtx) {
		// Arrange
		emptyToken := ""
		role := token.UserRole

		// Act
		payload, err := s.maker.VerifyToken(emptyToken, role)

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, token.ErrInvalidToken)
		sCtx.Require().Nil(payload)
	})
}

func (s *TokenMakerSuite) TestJWTMaker_MalformedToken(t provider.T) {
	t.WithNewStep("Верификация поврежденного токена", func(sCtx provider.StepCtx) {
		// Arrange
		malformedToken := "malformed.jwt.token"
		role := token.UserRole

		// Act
		payload, err := s.maker.VerifyToken(malformedToken, role)

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, token.ErrInvalidToken)
		sCtx.Require().Nil(payload)
	})
}
