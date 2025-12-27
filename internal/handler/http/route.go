package http

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler, postHandler *PostHandler) *gin.Engine {
	router := gin.Default()

	// --- Public Routes ---
	// Route for authentication endpoints
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	// Post Routes (Public)
	router.GET("/posts", postHandler.GetAll)
	router.GET("/posts/:id", postHandler.GetByID)

	// --- Protected Routes ---
	protected := router.Group("/")
	protected.Use(h.AuthMiddleware())
	{
		protected.POST("/posts", postHandler.Create)
		protected.PUT("/posts/:id", postHandler.Update)
		protected.DELETE("/posts/:id", postHandler.Delete)
	}

	return router
}
