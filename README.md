# Product API

API для управления продуктами, разработанное на Go с использованием Gin framework.

## Требования

- Go 1.16 или выше
- Gin framework
- Postman (для тестирования)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Afra1m/product_api.git
cd product_api
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

Сервер будет запущен на `http://localhost:8080`.

## API Endpoints

### Базовые CRUD операции

- `GET /api/products` - Получить список всех продуктов
- `GET /api/products/:id` - Получить продукт по ID
- `POST /api/products` - Создать новый продукт
- `PUT /api/products/:id` - Обновить существующий продукт
- `DELETE /api/products/:id` - Удалить продукт

### Фильтрация и поиск

- `GET /api/products/category/:category` - Получить продукты по категории
- `GET /api/products/search?q=query` - Поиск продуктов
- `GET /api/products/price-range?min=X&max=Y` - Получить продукты в указанном диапазоне цен
- `GET /api/products/in-stock` - Получить продукты в наличии
- `PUT /api/products/:id/stock` - Обновить количество товара

### Категории и статистика

- `GET /api/products/categories` - Получить список всех категорий
- `GET /api/products/stats` - Получить статистику по продуктам

### Пакетные операции

- `POST /api/products/batch` - Создать несколько продуктов
- `PUT /api/products/batch` - Обновить несколько продуктов
- `DELETE /api/products/batch` - Удалить несколько продуктов

### История и аналитика

- `GET /api/products/:id/history` - Получить историю изменений продукта
- `GET /api/products/popular?limit=N` - Получить популярные продукты
- `GET /api/products/new?limit=N` - Получить новые продукты
- `GET /api/products/discount` - Получить продукты со скидкой
- `PUT /api/products/:id/discount` - Обновить скидку продукта

### Рекомендации

- `GET /api/products/similar/:id` - Получить похожие продукты
- `GET /api/products/related/:id` - Получить связанные продукты
- `GET /api/products/trending?limit=N` - Получить трендовые продукты
- `GET /api/products/featured` - Получить рекомендуемые продукты
- `PUT /api/products/:id/feature` - Обновить статус рекомендации

### Импорт/экспорт

- `GET /api/products/export` - Экспорт продуктов
- `POST /api/products/import` - Импорт продуктов

### Валидация и проверка

- `GET /api/products/validate/:id` - Проверить валидность продукта
- `GET /api/products/duplicates` - Получить дубликаты продуктов
- `GET /api/products/out-of-stock` - Получить продукты, которых нет в наличии
- `GET /api/products/low-stock?threshold=N` - Получить продукты с низким запасом

## Модели данных

### Product

```go
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Category    string    `json:"category"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Discount    float64   `json:"discount"`
    Featured    bool      `json:"featured"`
    Popularity  int       `json:"popularity"`
    Views       int       `json:"views"`
    Tags        []string  `json:"tags"`
    SKU         string    `json:"sku"`
    Barcode     string    `json:"barcode"`
    Weight      float64   `json:"weight"`
    Dimensions  string    `json:"dimensions"`
    Status      string    `json:"status"`
    History     []ProductHistory `json:"history"`
}
```

### ProductInput

```go
type ProductInput struct {
    Name        string    `json:"name" binding:"required"`
    Description string    `json:"description"`
    Price       float64   `json:"price" binding:"required,gt=0"`
    Category    string    `json:"category" binding:"required"`
    Stock       int       `json:"stock" binding:"required,gte=0"`
    Discount    float64   `json:"discount" binding:"gte=0,lte=100"`
    Featured    bool      `json:"featured"`
    Tags        []string  `json:"tags"`
    SKU         string    `json:"sku"`
    Barcode     string    `json:"barcode"`
    Weight      float64   `json:"weight"`
    Dimensions  string    `json:"dimensions"`
    Status      string    `json:"status"`
}
```

### ProductHistory

```go
type ProductHistory struct {
    Field     string      `json:"field"`
    OldValue  interface{} `json:"old_value"`
    NewValue  interface{} `json:"new_value"`
    Timestamp time.Time   `json:"timestamp"`
}
```

### ProductStats

```go
type ProductStats struct {
    TotalProducts     int     `json:"total_products"`
    TotalCategories   int     `json:"total_categories"`
    AveragePrice      float64 `json:"average_price"`
    TotalStock        int     `json:"total_stock"`
    OutOfStockCount   int     `json:"out_of_stock_count"`
    LowStockCount     int     `json:"low_stock_count"`
    DiscountedCount   int     `json:"discounted_count"`
    FeaturedCount     int     `json:"featured_count"`
    MostPopularCategory string `json:"most_popular_category"`
}
```

## Тестирование

Для тестирования API используется Postman. Коллекция тестов находится в файле `postman_collection.json`.

### Настройка переменных окружения в Postman

1. Создайте новое окружение в Postman
2. Добавьте следующие переменные:
   - `base_url`: `http://localhost:8080`
   - `product_id`: ID продукта (будет установлен после создания продукта)

### Запуск тестов

1. Импортируйте коллекцию `postman_collection.json` в Postman
2. Выберите созданное окружение
3. Запустите коллекцию

## Структура проекта

```
product_api/
├── main.go              # Точка входа приложения
├── models/
│   └── product.go       # Модели данных
├── storage/
│   └── product_storage.go # Хранилище данных
├── handlers/
│   └── product_handler.go # Обработчики HTTP-запросов
├── postman_collection.json # Коллекция тестов Postman
└── README.md            # Документация
```

