package models

import (
	"time"
	"gorm.io/gorm"
)

type Veiculo struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Marca         string         `json:"marca" gorm:"size:50;index"`
	Modelo        string         `json:"modelo" gorm:"size:100;index"`
	Placa         string         `json:"placa" gorm:"not null;unique;size:10;index" binding:"required"`
	Cor           string         `json:"cor" gorm:"size:30"`
	AnoModelo     string         `json:"anoModelo" gorm:"column:ano_modelo;size:10"`
	ClienteID     uint           `json:"clienteId" gorm:"not null;index"`
	OrdemServico  string         `json:"ordemServico" gorm:"column:ordem_servico;size:30;not null" binding:"required"`
	CreatedAt     time.Time      `json:"criadoEm" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"atualizadoEm" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Veiculo) TableName() string {
	return "veiculos"
}

func (v *Veiculo) BeforeCreate(tx *gorm.DB) error {
	return nil
}
