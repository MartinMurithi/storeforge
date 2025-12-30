package handler

import (
	"fmt"
	"net/http"

	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/handler/httpx"
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

	fmt.Println("user ..... ", user)

	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response := mapper.ToRegisterUserResponse(user)

	httpx.JSON(c, http.StatusCreated, response)

}

func (handler *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequestDTO

	fmt.Println("user login 1..... ")

	err := c.ShouldBindJSON(&req)
	
fmt.Println("user login 2..... ")
	// Invalid request
	if err != nil {
		httpx.ValidationError(c)
	}

	input := &services.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	fmt.Println("user login 3..... ")

	user, token, err := handler.UserService.LoginUser(input, c.Request.Context())

	fmt.Println("user login 4..... ", err)

	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}
	fmt.Println("user login 5..... ")

	response := mapper.ToLoginUserResponse(token, user)
	fmt.Println("user login 6..... ")

	httpx.JSON(c, http.StatusCreated, response)

}
