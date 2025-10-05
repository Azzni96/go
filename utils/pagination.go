package utils

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Skip     int `json:"skip"`
}

// NewPaginationParams creates new pagination parameters with validation
func NewPaginationParams(page, pageSize int) PaginationParams {
	// Set default values
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	// Limit maximum page size
	if pageSize > 100 {
		pageSize = 100
	}

	skip := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Skip:     skip,
	}
}

// PaginationResponse holds pagination response data
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginationResponse creates a pagination response
func NewPaginationResponse(data interface{}, page, pageSize int, totalItems int64) PaginationResponse {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
