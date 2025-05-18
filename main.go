package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Afra1m/product_api/handlers"
	"github.com/Afra1m/product_api/storage"
)

func main() {
	// Инициализация хранилища
	productStorage := storage.NewProductStorage()

	// Инициализация обработчиков
	productHandler := handlers.NewProductHandler(productStorage)

	// Создание маршрутизатора
	router := gin.Default()

	// Группа API
	api := router.Group("/api")
	{
		// Продукты
		products := api.Group("/products")
		{
			// Базовые CRUD операции
			products.GET("", productHandler.GetAllProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)

			// Фильтрация и поиск
			products.GET("/category/:category", productHandler.GetProductsByCategory)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/price-range", productHandler.GetProductsByPriceRange)
			products.GET("/in-stock", productHandler.GetProductsInStock)
			products.PUT("/:id/stock", productHandler.UpdateProductStock)

			// Категории и статистика
			products.GET("/categories", productHandler.GetAllCategories)
			products.GET("/stats", productHandler.GetProductStats)

			// Пакетные операции
			products.POST("/batch", productHandler.CreateBatchProducts)
			products.PUT("/batch", productHandler.UpdateBatchProducts)
			products.DELETE("/batch", productHandler.DeleteBatchProducts)

			// История и аналитика
			products.GET("/:id/history", productHandler.GetProductHistory)
			products.GET("/popular", productHandler.GetPopularProducts)
			products.GET("/new", productHandler.GetNewProducts)
			products.GET("/discount", productHandler.GetDiscountedProducts)
			products.PUT("/:id/discount", productHandler.UpdateProductDiscount)

			// Рекомендации
			products.GET("/similar/:id", productHandler.GetSimilarProducts)
			products.GET("/related/:id", productHandler.GetRelatedProducts)
			products.GET("/trending", productHandler.GetTrendingProducts)
			products.GET("/featured", productHandler.GetFeaturedProducts)
			products.PUT("/:id/feature", productHandler.UpdateProductFeature)

			// Импорт/экспорт
			products.GET("/export", productHandler.ExportProducts)
			products.POST("/import", productHandler.ImportProducts)

			// Валидация и проверка
			products.GET("/validate/:id", productHandler.ValidateProduct)
			products.GET("/duplicates", productHandler.GetDuplicateProducts)
			products.GET("/out-of-stock", productHandler.GetOutOfStockProducts)
			products.GET("/low-stock", productHandler.GetLowStockProducts)
		}
	}

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Не удалось запустить сервер:", err)
	}
}
