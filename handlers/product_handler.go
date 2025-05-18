package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Afra1m/product_api/models"
	"github.com/Afra1m/product_api/storage"
)

// ProductHandler представляет собой обработчик для продуктов
type ProductHandler struct {
	storage *storage.ProductStorage
}

// NewProductHandler создает новый обработчик продуктов
func NewProductHandler(storage *storage.ProductStorage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

// GetAllProducts возвращает список всех продуктов
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products := h.storage.GetAll()
	c.JSON(http.StatusOK, products)
}

// GetProductByID возвращает продукт по ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// CreateProduct создает новый продукт
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	product := models.Product{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		Stock:       input.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
		Discount:    input.Discount,
		Featured:    input.Featured,
		Tags:        input.Tags,
		SKU:         input.SKU,
		Barcode:     input.Barcode,
		Weight:      input.Weight,
		Dimensions:  input.Dimensions,
		Status:      input.Status,
	}

	if err := h.storage.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct обновляет существующий продукт
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingProduct, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	existingProduct.Name = input.Name
	existingProduct.Description = input.Description
	existingProduct.Price = input.Price
	existingProduct.Category = input.Category
	existingProduct.Stock = input.Stock
	existingProduct.Discount = input.Discount
	existingProduct.Featured = input.Featured
	existingProduct.Tags = input.Tags
	existingProduct.SKU = input.SKU
	existingProduct.Barcode = input.Barcode
	existingProduct.Weight = input.Weight
	existingProduct.Dimensions = input.Dimensions
	existingProduct.Status = input.Status
	existingProduct.UpdatedAt = time.Now()

	if err := h.storage.Update(id, existingProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingProduct)
}

// DeleteProduct удаляет продукт
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.storage.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetProductsByCategory возвращает продукты по категории
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	products := h.storage.GetByCategory(category)
	c.JSON(http.StatusOK, products)
}

// SearchProducts ищет продукты
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр поиска 'q' обязателен"})
		return
	}

	products := h.storage.GetAll()
	var results []models.Product
	for _, product := range products {
		if contains(product.Name, query) || contains(product.Description, query) {
			results = append(results, product)
		}
	}
	c.JSON(http.StatusOK, results)
}

// GetProductsByPriceRange возвращает продукты в указанном диапазоне цен
func (h *ProductHandler) GetProductsByPriceRange(c *gin.Context) {
	minStr := c.Query("min")
	maxStr := c.Query("max")

	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат минимальной цены"})
		return
	}

	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат максимальной цены"})
		return
	}

	products := h.storage.GetByPriceRange(min, max)
	c.JSON(http.StatusOK, products)
}

// GetProductsInStock возвращает продукты в наличии
func (h *ProductHandler) GetProductsInStock(c *gin.Context) {
	products := h.storage.GetInStock()
	c.JSON(http.StatusOK, products)
}

// UpdateProductStock обновляет количество товара
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Stock int `json:"stock" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.UpdateStock(id, input.Stock); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetAllCategories возвращает список всех категорий
func (h *ProductHandler) GetAllCategories(c *gin.Context) {
	categories := h.storage.GetAllCategories()
	c.JSON(http.StatusOK, categories)
}

// GetProductStats возвращает статистику по продуктам
func (h *ProductHandler) GetProductStats(c *gin.Context) {
	stats := h.storage.GetStats()
	c.JSON(http.StatusOK, stats)
}

// CreateBatchProducts создает несколько продуктов
func (h *ProductHandler) CreateBatchProducts(c *gin.Context) {
	var input []models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products := make([]models.Product, len(input))
	now := time.Now()

	for i, in := range input {
		products[i] = models.Product{
			ID:          uuid.New().String(),
			Name:        in.Name,
			Description: in.Description,
			Price:       in.Price,
			Category:    in.Category,
			Stock:       in.Stock,
			CreatedAt:   now,
			UpdatedAt:   now,
			Discount:    in.Discount,
			Featured:    in.Featured,
			Tags:        in.Tags,
			SKU:         in.SKU,
			Barcode:     in.Barcode,
			Weight:      in.Weight,
			Dimensions:  in.Dimensions,
			Status:      in.Status,
		}
	}

	if err := h.storage.CreateBatch(products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, products)
}

// UpdateBatchProducts обновляет несколько продуктов
func (h *ProductHandler) UpdateBatchProducts(c *gin.Context) {
	var input models.BatchProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]models.Product)
	for _, id := range input.IDs {
		product, err := h.storage.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "продукт с ID " + id + " не найден"})
			return
		}

		product.Name = input.Update.Name
		product.Description = input.Update.Description
		product.Price = input.Update.Price
		product.Category = input.Update.Category
		product.Stock = input.Update.Stock
		product.Discount = input.Update.Discount
		product.Featured = input.Update.Featured
		product.Tags = input.Update.Tags
		product.SKU = input.Update.SKU
		product.Barcode = input.Update.Barcode
		product.Weight = input.Update.Weight
		product.Dimensions = input.Update.Dimensions
		product.Status = input.Update.Status
		product.UpdatedAt = time.Now()

		updates[id] = product
	}

	if err := h.storage.UpdateBatch(updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteBatchProducts удаляет несколько продуктов
func (h *ProductHandler) DeleteBatchProducts(c *gin.Context) {
	var input struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.DeleteBatch(input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProductHistory возвращает историю изменений продукта
func (h *ProductHandler) GetProductHistory(c *gin.Context) {
	id := c.Param("id")
	product, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product.History)
}

// GetPopularProducts возвращает популярные продукты
func (h *ProductHandler) GetPopularProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат лимита"})
		return
	}

	products := h.storage.GetPopular(limit)
	c.JSON(http.StatusOK, products)
}

// GetNewProducts возвращает новые продукты
func (h *ProductHandler) GetNewProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат лимита"})
		return
	}

	products := h.storage.GetNew(limit)
	c.JSON(http.StatusOK, products)
}

// GetDiscountedProducts возвращает продукты со скидкой
func (h *ProductHandler) GetDiscountedProducts(c *gin.Context) {
	products := h.storage.GetDiscounted()
	c.JSON(http.StatusOK, products)
}

// UpdateProductDiscount обновляет скидку продукта
func (h *ProductHandler) UpdateProductDiscount(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Discount float64 `json:"discount" binding:"required,gte=0,lte=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.UpdateDiscount(id, input.Discount); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetSimilarProducts возвращает похожие продукты
func (h *ProductHandler) GetSimilarProducts(c *gin.Context) {
	id := c.Param("id")
	product, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	products := h.storage.GetAll()
	var similar []models.Product
	for _, p := range products {
		if p.ID != id && p.Category == product.Category {
			similar = append(similar, p)
		}
	}

	c.JSON(http.StatusOK, similar)
}

// GetRelatedProducts возвращает связанные продукты
func (h *ProductHandler) GetRelatedProducts(c *gin.Context) {
	id := c.Param("id")
	product, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	products := h.storage.GetAll()
	var related []models.Product
	for _, p := range products {
		if p.ID != id && hasCommonTags(p.Tags, product.Tags) {
			related = append(related, p)
		}
	}

	c.JSON(http.StatusOK, related)
}

// GetTrendingProducts возвращает трендовые продукты
func (h *ProductHandler) GetTrendingProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат лимита"})
		return
	}

	products := h.storage.GetPopular(limit)
	c.JSON(http.StatusOK, products)
}

// GetFeaturedProducts возвращает рекомендуемые продукты
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	products := h.storage.GetFeatured()
	c.JSON(http.StatusOK, products)
}

// UpdateProductFeature обновляет статус рекомендации продукта
func (h *ProductHandler) UpdateProductFeature(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Featured bool `json:"featured"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.UpdateFeature(id, input.Featured); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ExportProducts экспортирует продукты
func (h *ProductHandler) ExportProducts(c *gin.Context) {
	products := h.storage.GetAll()
	c.JSON(http.StatusOK, products)
}

// ImportProducts импортирует продукты
func (h *ProductHandler) ImportProducts(c *gin.Context) {
	var input []models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products := make([]models.Product, len(input))
	now := time.Now()

	for i, in := range input {
		products[i] = models.Product{
			ID:          uuid.New().String(),
			Name:        in.Name,
			Description: in.Description,
			Price:       in.Price,
			Category:    in.Category,
			Stock:       in.Stock,
			CreatedAt:   now,
			UpdatedAt:   now,
			Discount:    in.Discount,
			Featured:    in.Featured,
			Tags:        in.Tags,
			SKU:         in.SKU,
			Barcode:     in.Barcode,
			Weight:      in.Weight,
			Dimensions:  in.Dimensions,
			Status:      in.Status,
		}
	}

	if err := h.storage.CreateBatch(products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, products)
}

// ValidateProduct проверяет валидность продукта
func (h *ProductHandler) ValidateProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.storage.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	validation := make(map[string]bool)
	validation["name"] = product.Name != ""
	validation["price"] = product.Price > 0
	validation["category"] = product.Category != ""
	validation["stock"] = product.Stock >= 0
	validation["sku"] = product.SKU != ""
	validation["barcode"] = product.Barcode != ""

	c.JSON(http.StatusOK, validation)
}

// GetDuplicateProducts возвращает дубликаты продуктов
func (h *ProductHandler) GetDuplicateProducts(c *gin.Context) {
	products := h.storage.GetAll()
	duplicates := make(map[string][]models.Product)

	for _, p1 := range products {
		for _, p2 := range products {
			if p1.ID != p2.ID && p1.Name == p2.Name {
				duplicates[p1.Name] = append(duplicates[p1.Name], p1, p2)
			}
		}
	}

	c.JSON(http.StatusOK, duplicates)
}

// GetOutOfStockProducts возвращает продукты, которых нет в наличии
func (h *ProductHandler) GetOutOfStockProducts(c *gin.Context) {
	products := h.storage.GetOutOfStock()
	c.JSON(http.StatusOK, products)
}

// GetLowStockProducts возвращает продукты с низким запасом
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "10")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат порога"})
		return
	}

	products := h.storage.GetLowStock(threshold)
	c.JSON(http.StatusOK, products)
}

// Вспомогательные функции

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && s[0:len(substr)] == substr
}

func hasCommonTags(tags1, tags2 []string) bool {
	for _, t1 := range tags1 {
		for _, t2 := range tags2 {
			if t1 == t2 {
				return true
			}
		}
	}
	return false
}
