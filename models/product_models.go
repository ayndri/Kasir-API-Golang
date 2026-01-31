package models

type Product struct {
	ID         int      `json:"id" gorm:"primaryKey"`
	Name       string   `json:"name"`
	Price      int      `json:"price"`
	Stock      int      `json:"stock"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"category"`
}