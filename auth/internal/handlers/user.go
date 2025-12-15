package handlers

import (
	"net/http"

	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/handlers/httpx"
	"github.com/MartinMurithi/storeforge/auth/internal/mapper"
	"github.com/MartinMurithi/storeforge/auth/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (handler *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequestDTO

	err := c.ShouldBindJSON(&req)

	// Invalid request
	if err != nil {
		httpx.ValidationError(c)
	}

	input := &services.RegisterUserInput{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password,
		BusinessType: req.BusinessType,
		BusinessName: req.BusinessName,
	}

	user, err := handler.UserService.RegisterUser(c.Request.Context(), input)

	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response := mapper.ToRegisterUserResponse(user)

	httpx.JSON(c, http.StatusCreated, response)

}

func (handler *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequestDTO

	err := c.ShouldBindJSON(&req)

	// Invalid request
	if err != nil {
		httpx.ValidationError(c)
	}

	input := &services.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := handler.UserService.LoginUser(input, c.Request.Context())

	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response := mapper.ToLoginUserResponse(token, user)

	httpx.JSON(c, http.StatusCreated, response)

}
