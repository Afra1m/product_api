package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/Afra1m/product_api/models"
)

// ProductStorage представляет собой хранилище продуктов
type ProductStorage struct {
	products map[string]models.Product
	mu       sync.RWMutex
}

// NewProductStorage создает новое хранилище продуктов
func NewProductStorage() *ProductStorage {
	return &ProductStorage{
		products: make(map[string]models.Product),
	}
}

// GetAll возвращает все продукты
func (s *ProductStorage) GetAll() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]models.Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

// GetByID возвращает продукт по ID
func (s *ProductStorage) GetByID(id string) (models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return models.Product{}, errors.New("продукт не найден")
	}
	return product, nil
}

// Create создает новый продукт
func (s *ProductStorage) Create(product models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[product.ID]; exists {
		return errors.New("продукт с таким ID уже существует")
	}

	s.products[product.ID] = product
	return nil
}

// Update обновляет существующий продукт
func (s *ProductStorage) Update(id string, product models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[id]; !exists {
		return errors.New("продукт не найден")
	}

	// Добавляем запись в историю
	oldProduct := s.products[id]
	if oldProduct.Name != product.Name {
		historyEntry := models.ProductHistory{
			Field:     "name",
			OldValue:  oldProduct.Name,
			NewValue:  product.Name,
			Timestamp: time.Now(),
		}
		product.History = append(oldProduct.History, historyEntry)
	}
	if oldProduct.Description != product.Description {
		historyEntry := models.ProductHistory{
			Field:     "description",
			OldValue:  oldProduct.Description,
			NewValue:  product.Description,
			Timestamp: time.Now(),
		}
		product.History = append(oldProduct.History, historyEntry)
	}
	if oldProduct.Price != product.Price {
		historyEntry := models.ProductHistory{
			Field:     "price",
			OldValue:  oldProduct.Price,
			NewValue:  product.Price,
			Timestamp: time.Now(),
		}
		product.History = append(oldProduct.History, historyEntry)
	}
	if oldProduct.Category != product.Category {
		historyEntry := models.ProductHistory{
			Field:     "category",
			OldValue:  oldProduct.Category,
			NewValue:  product.Category,
			Timestamp: time.Now(),
		}
		product.History = append(oldProduct.History, historyEntry)
	}
	if oldProduct.Stock != product.Stock {
		historyEntry := models.ProductHistory{
			Field:     "stock",
			OldValue:  oldProduct.Stock,
			NewValue:  product.Stock,
			Timestamp: time.Now(),
		}
		product.History = append(oldProduct.History, historyEntry)
	}

	s.products[id] = product
	return nil
}

// Delete удаляет продукт
func (s *ProductStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[id]; !exists {
		return errors.New("продукт не найден")
	}

	delete(s.products, id)
	return nil
}

// GetByCategory возвращает продукты по категории
func (s *ProductStorage) GetByCategory(category string) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Category == category {
			products = append(products, product)
		}
	}
	return products
}

// GetByPriceRange возвращает продукты в указанном диапазоне цен
func (s *ProductStorage) GetByPriceRange(min, max float64) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Price >= min && product.Price <= max {
			products = append(products, product)
		}
	}
	return products
}

// GetInStock возвращает продукты в наличии
func (s *ProductStorage) GetInStock() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Stock > 0 {
			products = append(products, product)
		}
	}
	return products
}

// UpdateStock обновляет количество товара
func (s *ProductStorage) UpdateStock(id string, stock int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[id]
	if !exists {
		return errors.New("продукт не найден")
	}

	oldStock := product.Stock
	product.Stock = stock
	product.UpdatedAt = time.Now()

	// Добавляем запись в историю
	historyEntry := models.ProductHistory{
		Field:     "stock",
		OldValue:  oldStock,
		NewValue:  stock,
		Timestamp: time.Now(),
	}
	product.History = append(product.History, historyEntry)

	s.products[id] = product
	return nil
}

// GetAllCategories возвращает список всех категорий
func (s *ProductStorage) GetAllCategories() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categories := make(map[string]struct{})
	for _, product := range s.products {
		categories[product.Category] = struct{}{}
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// GetStats возвращает статистику по продуктам
func (s *ProductStorage) GetStats() models.ProductStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := models.ProductStats{
		TotalProducts: len(s.products),
	}

	categories := make(map[string]struct{})
	var totalPrice float64
	var totalStock int
	var outOfStockCount int
	var lowStockCount int

	for _, product := range s.products {
		categories[product.Category] = struct{}{}
		totalPrice += product.Price
		totalStock += product.Stock
		if product.Stock == 0 {
			outOfStockCount++
		}
		if product.Stock < 10 {
			lowStockCount++
		}
	}

	stats.TotalCategories = len(categories)
	if len(s.products) > 0 {
		stats.AveragePrice = totalPrice / float64(len(s.products))
	}
	stats.TotalStock = totalStock
	stats.OutOfStockCount = outOfStockCount
	stats.LowStockCount = lowStockCount

	return stats
}

// CreateBatch создает несколько продуктов
func (s *ProductStorage) CreateBatch(products []models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, product := range products {
		if _, exists := s.products[product.ID]; exists {
			return errors.New("продукт с ID " + product.ID + " уже существует")
		}
		s.products[product.ID] = product
	}
	return nil
}

// UpdateBatch обновляет несколько продуктов
func (s *ProductStorage) UpdateBatch(updates map[string]models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, product := range updates {
		if _, exists := s.products[id]; !exists {
			return errors.New("продукт с ID " + id + " не найден")
		}
		s.products[id] = product
	}
	return nil
}

// DeleteBatch удаляет несколько продуктов
func (s *ProductStorage) DeleteBatch(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ids {
		if _, exists := s.products[id]; !exists {
			return errors.New("продукт с ID " + id + " не найден")
		}
		delete(s.products, id)
	}
	return nil
}

// GetPopular возвращает популярные продукты
func (s *ProductStorage) GetPopular(limit int) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]models.Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}

	// Сортировка по популярности
	for i := 0; i < len(products)-1; i++ {
		for j := i + 1; j < len(products); j++ {
			if products[i].Popularity < products[j].Popularity {
				products[i], products[j] = products[j], products[i]
			}
		}
	}

	if limit > len(products) {
		limit = len(products)
	}
	return products[:limit]
}

// GetNew возвращает новые продукты
func (s *ProductStorage) GetNew(limit int) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]models.Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}

	// Сортировка по дате создания
	for i := 0; i < len(products)-1; i++ {
		for j := i + 1; j < len(products); j++ {
			if products[i].CreatedAt.Before(products[j].CreatedAt) {
				products[i], products[j] = products[j], products[i]
			}
		}
	}

	if limit > len(products) {
		limit = len(products)
	}
	return products[:limit]
}

// GetDiscounted возвращает продукты со скидкой
func (s *ProductStorage) GetDiscounted() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Discount > 0 {
			products = append(products, product)
		}
	}
	return products
}

// UpdateDiscount обновляет скидку продукта
func (s *ProductStorage) UpdateDiscount(id string, discount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[id]
	if !exists {
		return errors.New("продукт не найден")
	}

	oldDiscount := product.Discount
	product.Discount = discount
	product.UpdatedAt = time.Now()

	// Добавляем запись в историю
	historyEntry := models.ProductHistory{
		Field:     "discount",
		OldValue:  oldDiscount,
		NewValue:  discount,
		Timestamp: time.Now(),
	}
	product.History = append(product.History, historyEntry)

	s.products[id] = product
	return nil
}

// GetFeatured возвращает рекомендуемые продукты
func (s *ProductStorage) GetFeatured() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Featured {
			products = append(products, product)
		}
	}
	return products
}

// UpdateFeature обновляет статус рекомендации продукта
func (s *ProductStorage) UpdateFeature(id string, featured bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[id]
	if !exists {
		return errors.New("продукт не найден")
	}

	oldFeatured := product.Featured
	product.Featured = featured
	product.UpdatedAt = time.Now()

	// Добавляем запись в историю
	historyEntry := models.ProductHistory{
		Field:     "featured",
		OldValue:  oldFeatured,
		NewValue:  featured,
		Timestamp: time.Now(),
	}
	product.History = append(product.History, historyEntry)

	s.products[id] = product
	return nil
}

// GetOutOfStock возвращает продукты, которых нет в наличии
func (s *ProductStorage) GetOutOfStock() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Stock == 0 {
			products = append(products, product)
		}
	}
	return products
}

// GetLowStock возвращает продукты с низким запасом
func (s *ProductStorage) GetLowStock(threshold int) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	for _, product := range s.products {
		if product.Stock > 0 && product.Stock <= threshold {
			products = append(products, product)
		}
	}
	return products
}
