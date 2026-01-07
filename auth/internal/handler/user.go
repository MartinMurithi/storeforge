package handler

import (
	"encoding/json"
	"errors"
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
			httpx.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "malformed JSON")
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
			httpx.Error(c, http.StatusBadRequest, "ALL_FIELDS_REQUIRED", "all fields are required")
		case errors.Is(err, apperrors.ErrInvalidEmailFormat):
			httpx.Error(c, http.StatusBadRequest, "INVALID_EMAIL_FORMAT", "invalid email format")
		case errors.Is(err, apperrors.ErrInvalidPhoneNumber):
			httpx.Error(c, http.StatusBadRequest, "INVALID_PHONE_NUMBER", "invalid phone number")
		case errors.Is(err, apperrors.ErrUserAlreadyExists):
			httpx.Error(c, http.StatusConflict, "USER_ALREADY_EXISTS", "email already registered")
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	response := mapper.ToRegisterUserResponse(user)

	httpx.JSON(c, http.StatusCreated, response)

}

func (handler *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequestDTO

	err := c.ShouldBindJSON(&req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			log.Printf("[LoginUser] malformed JSON: %v", err)
			httpx.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "malformed JSON")
			return
		}
	}

	user, token, err := handler.UserService.LoginUser(c.Request.Context(), &req)

	if err != nil {
		switch {
		case
			errors.Is(err, apperrors.ErrEmailRequired),
			errors.Is(err, apperrors.ErrPasswordRequired):
			httpx.Error(c, http.StatusBadRequest, "EMAIL_AND_PASSWORD_REQUIRED", "both email and password are required")
		case errors.Is(err, apperrors.ErrInvalidEmailFormat):
			httpx.Error(c, http.StatusBadRequest, "INVALID_EMAIL_FORMAT", "invalid email format")
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			httpx.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	response := mapper.ToLoginUserResponse(token, user)

	httpx.JSON(c, http.StatusOK, response)

}

func (handler *UserHandler) FetchAllUsers(c *gin.Context) {

	p, err := dto.ParsePagination(c)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidPageNumber):
			httpx.Error(c, http.StatusBadRequest, "INVALID_PAGE_NUMBER", "invalid page number")
		case errors.Is(err, apperrors.ErrInvalidLimitNumber):
			httpx.Error(c, http.StatusBadRequest, "INVALID_LIMIT_NUMBER", "invalid limit number")
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
	}

	users, meta, err := handler.UserService.FetchAllUsers(c.Request.Context(), p)

	if err != nil {
		log.Printf("[FetchAllUsers] failed to fetch users: %v", err)
		httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response := mapper.ToFetchAllUsersResponse(users, meta)

	httpx.JSON(c, http.StatusOK, response)

}

func (handler *UserHandler) GetCurrentUser(c *gin.Context) {

	id, err := dto.GetUserId(c)

	log.Printf("[HANDLER]: user id %v", id.Valid)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrIdIsRequired):
			httpx.Error(c, http.StatusNotFound, "ID_NOT_FOUND", "id not found in context")
		case errors.Is(err, apperrors.ErrInvalidUserIdFormat):
			httpx.Error(c, http.StatusBadRequest, "INVALID_USER_ID_FORMAT", "invalid user id format")
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
	}

	user, err := handler.UserService.GetCurrentUserById(c.Request.Context(), id)

	if err != nil {
		log.Printf("[FetchUser] failed to fetch user: %v", err)
		httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response := mapper.ToFetchUserResponse(user)

	httpx.JSON(c, http.StatusOK, response)

}

func (h *UserHandler) PatchMe(c *gin.Context) {
	var req dto.PatchUserRequestDTO

	err := c.ShouldBindJSON(&req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			log.Printf("[PatchMe] malformed JSON: %v", err)
			httpx.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "malformed JSON")
			return
		}
	}

	id, err := dto.GetUserId(c)

	log.Printf("[HANDLER]: user id %v", id.Valid)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrIdIsRequired):
			httpx.Error(c, http.StatusNotFound, "ID_NOT_FOUND", "id not found")
		case errors.Is(err, apperrors.ErrInvalidUserIdFormat):
			httpx.Error(c, http.StatusBadRequest, "INVALID_USER_ID_FORMAT", "invalid user id format")
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	input := mapper.ToPatchUserRequest(id, &req)

	updatedUser, err := h.UserService.UpdateCurrentUser(c.Request.Context(), input)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			httpx.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}
		log.Printf("[PatchMe] failed to update user: %v", err)
		httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response := mapper.ToFetchUserResponse(updatedUser)

	httpx.JSON(c, http.StatusOK, response)
}
