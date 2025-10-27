package models

import "time"

type SensorData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    string    `json:"status"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}
