package product

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
}

type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

func ParsePagination(c *gin.Context) (Pagination, error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "15")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return Pagination{}, errors.New("invalid page")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		return Pagination{}, errors.New("invalid limit")
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}, nil
}
