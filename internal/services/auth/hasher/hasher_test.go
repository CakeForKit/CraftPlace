package hasher_test

import (
	"testing"

	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
)

type HasherSuite struct {
	suite.Suite
	hasherServ hasher.Hasher
}

func TestHasher(t *testing.T) {
	suite.RunSuite(t, new(HasherSuite))
}

func (s *HasherSuite) BeforeEach(t provider.T) {
	t.Epic("Password")
	t.Feature("Hasher Service")

	hasherServ, err := hasher.NewHasher()
	require.NoError(t, err)
	s.hasherServ = hasherServ
}

func (s *HasherSuite) TestHasher_HashAndCheckPassword(t provider.T) {
	t.WithNewStep("Успешное хэширование и проверка пароля", func(sCtx provider.StepCtx) {
		// Arrange
		password := "1234were### _"

		// Act
		hashedPassword, err := s.hasherServ.HashPassword(password)

		// Assert
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(hashedPassword)

		// Verify
		err = s.hasherServ.CheckPassword(password, hashedPassword)
		sCtx.Require().NoError(err)
	})
}

func (s *HasherSuite) TestHasher_WrongPassword(t provider.T) {
	t.WithNewStep("Проверка неверного пароля", func(sCtx provider.StepCtx) {
		// Arrange
		password := "1234were### _"
		wrongPassword := "0"

		hashedPassword, err := s.hasherServ.HashPassword(password)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(hashedPassword)

		// Act
		err = s.hasherServ.CheckPassword(wrongPassword, hashedPassword)

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, hasher.ErrPassword)
	})
}

func (s *HasherSuite) TestHasher_EmptyPassword(t provider.T) {
	t.WithNewStep("Хэширование пустого пароля", func(sCtx provider.StepCtx) {
		// Arrange
		password := ""

		// Act
		hashedPassword, err := s.hasherServ.HashPassword(password)

		// Assert
		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, hasher.ErrEmptyPassword)
		sCtx.Require().Empty(hashedPassword)
	})
}

func (s *HasherSuite) TestHasher_DifferentHashesForSamePassword(t provider.T) {
	t.WithNewStep("Разные хэши для одного пароля", func(sCtx provider.StepCtx) {
		// Arrange
		password := "1234"

		// Act
		hashedPassword1, err := s.hasherServ.HashPassword(password)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(hashedPassword1)

		hashedPassword2, err := s.hasherServ.HashPassword(password)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(hashedPassword2)

		// Assert
		sCtx.Assert().NotEqual(hashedPassword1, hashedPassword2,
			"Хэши для одного пароля должны быть разными из-за соли")
	})
}

func (s *HasherSuite) TestHasher_ValidHashFormat(t provider.T) {
	t.WithNewStep("Проверка формата хэша", func(sCtx provider.StepCtx) {
		// Arrange
		password := "secure_password_123"

		// Act
		hashedPassword, err := s.hasherServ.HashPassword(password)
		sCtx.Require().NoError(err)

		// Assert
		sCtx.Require().NotEmpty(hashedPassword)
		sCtx.Assert().Greater(len(hashedPassword), 10,
			"Хэш должен быть достаточно длинным")
	})
}

func (s *HasherSuite) TestHasher_MultiplePasswordChecks(t provider.T) {
	t.WithNewStep("Многократная проверка одного пароля", func(sCtx provider.StepCtx) {
		// Arrange
		password := "test_password"
		hashedPassword, err := s.hasherServ.HashPassword(password)
		sCtx.Require().NoError(err)

		// Act & Assert - проверяем несколько раз
		for i := 0; i < 3; i++ {
			err = s.hasherServ.CheckPassword(password, hashedPassword)
			sCtx.Require().NoError(err, "Проверка должна успешно проходить каждый раз")
		}
	})
}

func (s *HasherSuite) TestHasher_SpecialCharactersPassword(t provider.T) {
	t.WithNewStep("Пароль со специальными символами", func(sCtx provider.StepCtx) {
		// Arrange
		passwords := []string{
			"p@ssw0rd!",
			"пароль123",
			"🔐secure123",
			"very_long_password_with_many_characters_1234567890!@#$%^&*()",
		}

		for _, password := range passwords {
			sCtx.WithNewStep("Проверка пароля: "+password, func(stepCtx provider.StepCtx) {
				// Act
				hashedPassword, err := s.hasherServ.HashPassword(password)
				stepCtx.Require().NoError(err)
				stepCtx.Require().NotEmpty(hashedPassword)

				// Verify
				err = s.hasherServ.CheckPassword(password, hashedPassword)
				stepCtx.Require().NoError(err)
			})
		}
	})
}
