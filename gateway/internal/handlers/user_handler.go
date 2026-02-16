package handlers

import (
	"net/http"

	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/client"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/pkg/errconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserClient *client.UserClient
}

func (h *UserHandler) GetMe(c *gin.Context) {
	// Get ID from context (Injected by AuthMiddleware)
	userID, err := request.GetUserId(c)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}

	// Call User Service via gRPC Client
	res, err := h.UserClient.GetCurrentUser(c.Request.Context(), &userv1.GetCurrentUserRequest{
		Id: userID,
	})

	if err != nil {
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
		code, slug, msg := errconv.(err)
		response.Error(c, code, slug, msg)
		return
	}

	// 2. Access the fields from the struct for the gRPC call
	res, err := h.UserClient.GetAllUsers(c.Request.Context(), &userv1.GetAllUsersRequest{
		// Access via pagination.Page and pagination.Limit
		Page:  int32(pagination.Page),
		Limit: int32(pagination.Limit),
	})

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	// Map a slice of Protos to a slice of DTOs
	response.JSON(c, http.StatusOK, mapper.MapUserProtoListToDTO(res.Users, res.Meta))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	var reqDTO dto.UpdateUserRequestDTO
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "Invalid request body")
		return
	}

	userID, _ := dto.GetUserId(c)

	res, err := h.UserClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id:           userID,
		FullName:     reqDTO.FullName,
		BusinessName: reqDTO.BusinessName,
		// ... other fields
	})

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusOK, mapper.MapUserProtoToDTO(res.User))
}

func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID, _ := dto.GetUserId(c)

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
