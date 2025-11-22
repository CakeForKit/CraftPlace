package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

type TokenVerifier interface {
	VerifyByToken(tokenStr string) (*tokenmaker.Payload, error)
}

func tokenFromHeader(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}
	accessToken := fields[1]
	return accessToken, nil
}

func AuthMiddleware(authServ TokenVerifier, authZ auth.AuthZ) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Auth middleware\n\n")
		accessToken, err := tokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf(" AuthMiddleware %w", err.Error())})
			return
		}
		payload, err := authServ.VerifyByToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("AuthMiddleware %w", err.Error())})
			return
		}

		ctx := c.Request.Context()
		ctx = authZ.Authorize(ctx, *payload)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
