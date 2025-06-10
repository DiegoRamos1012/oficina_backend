package models

import (
	"time"

	"gorm.io/gorm"
)

type Estoque struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Nome          string         `json:"nome" gorm:"not null;size:100;index" binding:"required"`
	Codigo        string         `json:"codigo" gorm:"size:100;uniqueIndex"`
	Descricao     string         `json:"descricao" gorm:"type:text"`
	Categoria     string         `json:"categoria" gorm:"size:50;index"`
	Quantidade    int            `json:"quantidade" gorm:"default:0;not null"`
	EstoqueMinimo int            `json:"estoque_minimo" gorm:"default:5"`
	PrecoUnitario float64        `json:"preco_unitario" gorm:"type:decimal(10,2);not null;default:0.00"`
	PrecoVenda    float64        `json:"preco_venda" gorm:"type:decimal(10,2);not null;default:0.00"`
	Fornecedor    string         `json:"fornecedor" gorm:"size:100;index"`
	Status        string         `json:"status" gorm:"size:20;default:'disponível';index"`
	Observacoes   string         `json:"observacoes" gorm:"type:text"`
	CriadoEm      time.Time      `json:"criado_em" gorm:"autoCreateTime"`
	AtualizadoEm  time.Time      `json:"atualizado_em" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Estoque) TableName() string {
	return "estoque"
}

// BeforeCreate executa ações antes de criar um novo registro
func (e *Estoque) BeforeCreate(tx *gorm.DB) error {
	// Você poderia adicionar regras de negócio aqui, como:
	// - Validar se o preço de venda é maior que o preço unitário
	// - Gerar um código único automático se não for fornecido
	// - Normalizar dados como categoria (converter para minúsculas, etc.)

	return nil
}

// BeforeUpdate executa ações antes de atualizar um registro existente
func (e *Estoque) BeforeUpdate(tx *gorm.DB) error {
	// Lógica de validação para atualizações
	return nil
}

// CalcularLucro retorna o lucro estimado por unidade
func (e *Estoque) CalcularLucro() float64 {
	return e.PrecoVenda - e.PrecoUnitario
}

// CalcularValorTotal retorna o valor total do item em estoque
func (e *Estoque) CalcularValorTotal() float64 {
	return float64(e.Quantidade) * e.PrecoVenda
}

// PrecisaReposicao verifica se o estoque está abaixo do mínimo
func (e *Estoque) PrecisaReposicao() bool {
	return e.Quantidade < e.EstoqueMinimo
}
