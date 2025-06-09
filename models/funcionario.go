package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Funcionario struct {
	ID              uint            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Nome            string          `json:"nome" gorm:"not null;size:100;index" binding:"required"`
	Telefone        string          `json:"telefone" gorm:"not null;size:20" binding:"required"`
	TelefoneReserva *string         `json:"telefoneReserva" gorm:"size:20"`
	CPF             string          `json:"cpf" gorm:"not null;unique;size:14;index" binding:"required"`
	Endereco        string          `json:"endereco" gorm:"size:255"`
	DataNascimento  time.Time       `json:"dataNascimento" gorm:"column:data_nascimento"`
	DataAdmissao    time.Time       `json:"dataAdmissao" gorm:"column:data_admissao"`
	Salario         decimal.Decimal `json:"salario" gorm:"type:decimal(10,2)"`
	Observacoes     string          `json:"observacoes" gorm:"type:text"`
	Cargo           string          `json:"cargo" gorm:"size:50;index"`

	// Chave estrangeira para Usuário
	UsuarioID uint    `json:"usuarioId" gorm:"uniqueIndex"` // Um funcionário tem exatamente um usuário
	Usuario   Usuario `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`

	CreatedAt time.Time       `json:"criadoEm" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"atualizadoEm" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt  `json:"-" gorm:"index"` // Adicionado: suporte a soft delete
}

func (Funcionario) TableName() string {
	return "funcionarios"
}

func (f *Funcionario) BeforeCreate(tx *gorm.DB) error {
	// Lógica adicional antes de criar um funcionário
	return nil
}
