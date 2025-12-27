package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/elokanugrah/backend-takehome/internal/usecase"
	"github.com/gin-gonic/gin"
)

const (
	// claimsKey is the key used to store the user claims in the Gin context.
	claimsKey = "claims"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the header is in "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := h.authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		// Store the claims in the context for later handlers
		c.Set(claimsKey, claims)

		// Continue to the next handler
		c.Next()
	}
}

// getClaimsFromContext is a helper to retrieve claims from the Gin context.
func getClaimsFromContext(c *gin.Context) (*usecase.AuthClaims, error) {
	val, exists := c.Get(claimsKey)
	if !exists {
		return nil, errors.New("claims not found in context")
	}

	claims, ok := val.(*usecase.AuthClaims)
	if !ok {
		return nil, errors.New("could not parse claims from context")
	}

	return claims, nil
}
