package http

import (
	"net/http"
	"strconv"

	"github.com/elokanugrah/backend-takehome/internal/dto"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	usecase *usecase.CommentUseCase
}

func NewCommentHandler(uc *usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{usecase: uc}
}

func (h *CommentHandler) Create(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(claims.UserID)

	comment, err := h.usecase.Create(c.Request.Context(), postID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	comments, err := h.usecase.GetByPostID(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []dto.CommentResponse
	for _, c := range comments {
		res = append(res, dto.CommentResponse{
			ID:         c.ID,
			PostID:     c.PostID,
			AuthorName: c.AuthorName,
			Content:    c.Content,
			CreatedAt:  c.CreatedAt,
		})
	}

	if res == nil {
		res = []dto.CommentResponse{}
	}

	c.JSON(http.StatusOK, res)
}
