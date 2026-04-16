package shared

import (
	"strconv"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/gin-gonic/gin"
)

type PaginationDTO struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type PaginationMetaDTO struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int32 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func ParsePagination(c *gin.Context) (PaginationDTO, error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "15")

	// Use ParseInt to specify 32-bit to match Proto
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		return PaginationDTO{}, apperrors.ErrInvalidPageNumber
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit < 1 || limit > 100 {
		return PaginationDTO{}, apperrors.ErrInvalidLimitNumber
	}

	return PaginationDTO{
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}
