package model

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Price    float64 `json:"price" binding:"gte=0"`
	Stock    int     `json:"stock" binding:"gte=0"`
}

type ProductQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
	Category string `form:"category"`
	Stock    string `form:"stock"`
	Sort     string `form:"sort"`
	Order    string `form:"order"`
}

type PaginatedProducts struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
}
