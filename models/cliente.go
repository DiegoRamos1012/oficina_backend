package models

import (
	"gorm.io/gorm"
)

// Cliente representa a tabela clientes no banco de dados
type Cliente struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Nome           string         `json:"nome" gorm:"not null" binding:"required"`
	Email          *string        `json:"email"`
	Telefone       *string        `json:"telefone"`
	Endereco       string         `json:"endereco"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`                                 // Suporte a soft delete
	Veiculos       []Veiculo      `json:"veiculos,omitempty" gorm:"foreignKey:ClienteID"` // Relacionamento um para muitos
}

// TableName especifica o nome da tabela a ser usada
func (Cliente) TableName() string {
	return "clientes"
}

// BeforeCreate é um hook que é executado antes de criar um registro
func (c *Cliente) BeforeCreate(tx *gorm.DB) error {
	// Lógica personalizada antes da criação, se necessário
	return nil
}

// ClienteVeiculosDTO é um DTO para retornar cliente com seus veículos
// Embora com GORM você possa usar diretamente a relação, este DTO pode ser útil para casos específicos
type ClienteVeiculosDTO struct {
	Cliente  Cliente   `json:"cliente"`
	Veiculos []Veiculo `json:"veiculos"`
}
