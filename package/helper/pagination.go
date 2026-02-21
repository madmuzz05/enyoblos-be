package helper

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type PaginationRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Sort     string `query:"sort"`
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}

// ParsePaginationFromQuery mengparse pagination parameter dari query string
func ParsePaginationFromQuery(c fiber.Ctx) PaginationRequest {
	page := c.Query("page", "1")
	pageSize := c.Query("page_size", "10")
	sort := c.Query("sort", "ASC")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// Validasi
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}

	return PaginationRequest{
		Page:     pageInt,
		PageSize: pageSizeInt,
		Sort:     sort,
	}
}

// GetOffset menghitung offset untuk query database
func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// CalculateTotalPages menghitung total halaman
func CalculateTotalPages(totalRecords int64, pageSize int) int {
	return int(math.Ceil(float64(totalRecords) / float64(pageSize)))
}

// BuildPaginationResponse membuild response pagination
func BuildPaginationResponse(currentPage, pageSize int, totalRecords int64) PaginationResponse {
	totalPages := CalculateTotalPages(totalRecords, pageSize)
	return PaginationResponse{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
	}
}

// SendPaginatedResponse mengirim response pagination dengan format standar
func SendPaginatedResponse(c fiber.Ctx, statusCode int, message string, currentPage, pageSize int, totalRecords int64, data interface{}) error {
	pagination := BuildPaginationResponse(currentPage, pageSize, totalRecords)
	return c.Status(statusCode).JSON(fiber.Map{
		"status":     "success",
		"message":    message,
		"pagination": pagination,
		"data":       data,
	})
}
