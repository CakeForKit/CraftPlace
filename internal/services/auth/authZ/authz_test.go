package auth_test

import (
	"context"
	"testing"

	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	testobj "github.com/CakeForKit/CraftPlace.git/internal/tests/test_obj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert" // исправлено: было stretchr -> stretchr
)

type AuthZServiceSuite struct {
	suite.Suite
}

func TestAuthZService(t *testing.T) {
	suite.RunSuite(t, new(AuthZServiceSuite))
}

func (s *AuthZServiceSuite) TestAuthZ_Authz_GetUserID(t provider.T) {
	payloadMother := testobj.NewPayloadMother()
	authzServ, err := auth.NewAuthZ()
	t.Require().NoError(err, "Failed to create authzServ")

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		userID := uuid.New()
		payload := payloadMother.UserPayload(userID)

		ctx = authzServ.Authorize(ctx, payload)

		resUserID, err := authzServ.UserIDFromContext(ctx)
		sCtx.Require().NoError(err)
		assert.Equal(t, userID, resUserID)
	})
	t.WithNewStep("not found userID", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		_, err := authzServ.UserIDFromContext(ctx)

		assert.ErrorIs(t, err, auth.ErrNotAuthZ)
	})
}
