package models

import (
	"time"
)

// Product представляет собой модель продукта
type Product struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Category    string           `json:"category"`
	Stock       int              `json:"stock"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Discount    float64          `json:"discount"`
	Featured    bool             `json:"featured"`
	Popularity  int              `json:"popularity"`
	Views       int              `json:"views"`
	Tags        []string         `json:"tags"`
	SKU         string           `json:"sku"`
	Barcode     string           `json:"barcode"`
	Weight      float64          `json:"weight"`
	Dimensions  string           `json:"dimensions"`
	Status      string           `json:"status"`
	History     []ProductHistory `json:"history,omitempty"`
}

// ProductInput представляет собой структуру для создания/обновления продукта
type ProductInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Category    string   `json:"category" binding:"required"`
	Stock       int      `json:"stock" binding:"required,gte=0"`
	Discount    float64  `json:"discount" binding:"gte=0,lte=100"`
	Featured    bool     `json:"featured"`
	Tags        []string `json:"tags"`
	SKU         string   `json:"sku"`
	Barcode     string   `json:"barcode"`
	Weight      float64  `json:"weight"`
	Dimensions  string   `json:"dimensions"`
	Status      string   `json:"status"`
}

// ProductHistory представляет историю изменений продукта
type ProductHistory struct {
	Field     string      `json:"field"`
	OldValue  interface{} `json:"old_value"`
	NewValue  interface{} `json:"new_value"`
	Timestamp time.Time   `json:"timestamp"`
}

// ProductStats представляет статистику по продуктам
type ProductStats struct {
	TotalProducts       int     `json:"total_products"`
	TotalCategories     int     `json:"total_categories"`
	AveragePrice        float64 `json:"average_price"`
	TotalStock          int     `json:"total_stock"`
	OutOfStockCount     int     `json:"out_of_stock_count"`
	LowStockCount       int     `json:"low_stock_count"`
	DiscountedCount     int     `json:"discounted_count"`
	FeaturedCount       int     `json:"featured_count"`
	MostPopularCategory string  `json:"most_popular_category"`
}

// BatchProductInput представляет структуру для пакетных операций
type BatchProductInput struct {
	IDs    []string     `json:"ids" binding:"required"`
	Update ProductInput `json:"update" binding:"required"`
}
