package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MartinMurithi/storeforge/auth/internal/apperrors"
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

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			log.Printf("[RegisterUser] malformed JSON: %v", err)
			httpx.Error(c, http.StatusBadRequest, "malformed JSON")
			return
		}
	}

	user, err := handler.UserService.RegisterUser(c.Request.Context(), &req)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrFullNameRequired),
			errors.Is(err, apperrors.ErrEmailRequired),
			errors.Is(err, apperrors.ErrPhoneRequired),
			errors.Is(err, apperrors.ErrPasswordRequired),
			errors.Is(err, apperrors.ErrBusinessTypeRequired),
			errors.Is(err, apperrors.ErrBusinessNameRequired):
			httpx.Error(c, http.StatusBadRequest, "all fields are required")
		case errors.Is(err, apperrors.ErrInvalidEmailFormat):
			httpx.Error(c, http.StatusBadRequest, "invalid email format")
		case errors.Is(err, apperrors.ErrInvalidPhoneNumber):
			httpx.Error(c, http.StatusBadRequest, "invalid phone number")
		case errors.Is(err, apperrors.ErrUserAlreadyExists):
			httpx.Error(c, http.StatusConflict, "email already registered")
		default:
			httpx.Error(c, http.StatusInternalServerError, "internal server error")
		}
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
