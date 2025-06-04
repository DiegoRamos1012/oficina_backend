package models

type Usuario struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Nome string `json:"nome" gorm:"not null" binding:"required"`
	Email string `json:"email" binding:"required"`
	Senha string `json:"senha" binding:"required"`
}
