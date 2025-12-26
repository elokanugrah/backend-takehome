package http

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler) *gin.Engine {
	router := gin.Default()

	// --- Public Routes ---
	// Route for authentication endpoints
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	return router
}
