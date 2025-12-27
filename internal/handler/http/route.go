package http

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler, postHandler *PostHandler, commentHandler *CommentHandler) *gin.Engine {
	router := gin.Default()

	// --- Public Routes ---
	// Authentication Routes (Public)
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	// Post Routes (Public)
	router.GET("/posts", postHandler.GetAll)
	router.GET("/posts/:id", postHandler.GetByID)

	// Comment Routes (Public)
	router.GET("/posts/:id/comments", commentHandler.GetByPostID)

	// --- Protected Routes ---
	protected := router.Group("/")
	protected.Use(h.AuthMiddleware())
	{
		// Post Routes (Protected)
		protected.POST("/posts", postHandler.Create)
		protected.PUT("/posts/:id", postHandler.Update)
		protected.DELETE("/posts/:id", postHandler.Delete)

		// Comment Routes (Protected)
		protected.POST("/posts/:id/comments", commentHandler.Create)
	}

	return router
}
