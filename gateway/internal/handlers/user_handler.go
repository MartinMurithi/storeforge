package handlers

import (
	"fmt"
	"net/http"

	userv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/gateway/internal/util"
	"github.com/MartinMurithi/storeforge/pkg/errconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserClient userv1.UserServiceClient
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Get ID from context (Injected by AuthMiddleware)
	userID, err := request.GetUserId(c)

	fmt.Println("error getting user: %w", err)

	if err != nil {
		fmt.Println("error getting user ID: %w", err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}

	// Call User Service via gRPC Client
	res, err := h.UserClient.GetCurrentUser(c.Request.Context(), &userv1.GetCurrentUserRequest{
		Id: userID,
	})

	fmt.Println("current user: %w", res)

	if err != nil {
		fmt.Println("error getting user: %w", err)
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := mapper.MapUserProtoToDTO(res.User)
	response.JSON(c, http.StatusOK, resp)
}

func (h *UserHandler) FetchAll(c *gin.Context) {
	pagination, err := shared.ParsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAGINATION", err.Error())
		return
	}

	// Access the fields from the struct for the gRPC call
	res, err := h.UserClient.GetAllUsers(c.Request.Context(), &userv1.GetAllUsersRequest{
		Page:  int32(pagination.Page),
		Limit: int32(pagination.Limit),
	})

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	userList := mapper.MapUserProtosToDTOs(res.Users)

	resp := dto.GetAllUsersResponseDTO{
		Users: userList,
		Meta: shared.PaginationMetaDTO{
			Page:       res.Meta.Page,
			Limit:      res.Meta.Limit,
			Total:      int64(res.Meta.Total),
			TotalPages: res.Meta.TotalPages,
			HasNext:    res.Meta.HasNext,
			HasPrev:    res.Meta.HasPrev,
		},
	}

	// Map a slice of Protos to a slice of DTOs
	response.JSON(c, http.StatusOK, resp)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	var reqDTO dto.UpdateUserRequestDTO

	if !util.BindAndValidateJSON(c, &reqDTO) {
		return
	}

	// Get ID from context (Injected by AuthMiddleware)
	userID, err := request.GetUserId(c)

	if err != nil {
		fmt.Println("error getting user ID: %w", err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}

	res, err := h.UserClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id:    userID,
		Email: reqDTO.Email,
		Phone: reqDTO.Phone,
	})

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusOK, mapper.MapUserProtoToDTO(res.User))
}

func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID, _ := request.GetUserId(c)

	_, err := h.UserClient.DeleteUser(c.Request.Context(), &userv1.DeleteUserRequest{
		Id: userID,
	})

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusNoContent, nil)
}
