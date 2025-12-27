package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/elokanugrah/backend-takehome/internal/dto"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	usecase *usecase.PostUseCase
}

func NewPostHandler(uc *usecase.PostUseCase) *PostHandler {
	return &PostHandler{usecase: uc}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
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

	post, err := h.usecase.Create(c.Request.Context(), req.Title, req.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author: dto.AuthorResponse{
			ID:    post.AuthorID,
			Name:  post.AuthorName,
			Email: post.AuthorEmail,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

func (h *PostHandler) GetAll(c *gin.Context) {
	posts, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []dto.PostResponse
	for _, p := range posts {
		res = append(res, dto.PostResponse{
			ID:      p.ID,
			Title:   p.Title,
			Content: p.Content,
			Author: dto.AuthorResponse{
				ID:    p.AuthorID,
				Name:  p.AuthorName,
				Email: p.AuthorEmail,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	if res == nil {
		res = []dto.PostResponse{}
	}

	c.JSON(http.StatusOK, res)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author: dto.AuthorResponse{
			ID:    post.AuthorID,
			Name:  post.AuthorName,
			Email: post.AuthorEmail,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePostRequest
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

	post, err := h.usecase.Update(c.Request.Context(), id, req.Title, req.Content, userID)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if err.Error() == "forbidden: only the author can update this post" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author: dto.AuthorResponse{
			ID:    post.AuthorID,
			Name:  post.AuthorName,
			Email: post.AuthorEmail,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(claims.UserID)

	err = h.usecase.Delete(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if err.Error() == "forbidden: only the author can delete this post" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}
