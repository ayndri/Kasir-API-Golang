package models

type Category struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products,omitempty"`
}