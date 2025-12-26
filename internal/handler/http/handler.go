package http

import (
	"github.com/elokanugrah/backend-takehome/internal/usecase"
)

type Handler struct {
	userUseCase *usecase.UserUseCase
	authService usecase.AuthService
}

func NewHandler(uuc *usecase.UserUseCase, as usecase.AuthService) *Handler {
	return &Handler{
		userUseCase: uuc,
		authService: as,
	}
}
