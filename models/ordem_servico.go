package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrdemServico representa uma ordem de serviço na oficina mecânica
type OrdemServico struct {
	ID                 uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	VeiculoID          uint           `json:"veiculoId" gorm:"not null;index" binding:"required"`
	Veiculo            Veiculo        `json:"veiculo,omitempty" gorm:"foreignKey:VeiculoID"`
	ClienteID          uint           `json:"clienteId" gorm:"not null;index" binding:"required"`
	Cliente            Cliente        `json:"cliente,omitempty" gorm:"foreignKey:ClienteID"`
	FuncionarioID      uint           `json:"funcionarioId" gorm:"index"`
	Funcionario        Funcionario    `json:"funcionario,omitempty" gorm:"foreignKey:FuncionarioID"`
	NumeroOS           string         `json:"numeroOS" gorm:"size:20;unique;index"`
	DataEntrada        time.Time      `json:"dataEntrada" gorm:"not null"`
	DataPrevisao       time.Time      `json:"dataPrevisao"`
	DataConclusao      *time.Time     `json:"dataConclusao"`
	Status             string         `json:"status" gorm:"not null;default:'aberta';size:20;index"` // Aberta, EmAndamento, Concluida, Cancelada
	Descricao          string         `json:"descricao" gorm:"type:text" binding:"required"`
	Diagnostico        string         `json:"diagnostico" gorm:"type:text"`
	ValorPecas         float64        `json:"valorPecas" gorm:"type:decimal(10,2);default:0"`
	ValorServico       float64        `json:"valorServico" gorm:"type:decimal(10,2);default:0"`
	ValorDesconto      float64        `json:"valorDesconto" gorm:"type:decimal(10,2);default:0"`
	ValorTotal         float64        `json:"valorTotal" gorm:"type:decimal(10,2);default:0"`
	FormaPagamento     string         `json:"formaPagamento" gorm:"size:50"`
	Observacoes        string         `json:"observacoes" gorm:"type:text"`
	ServicosRealizados string         `json:"servicosRealizados" gorm:"type:text"`
	CreatedAt          time.Time      `json:"criadoEm" gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `json:"atualizadoEm" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamento com itens utilizados
	ItensUtilizados []ItemOrdemServico `json:"itensUtilizados,omitempty" gorm:"foreignKey:OrdemServicoID"`
}

// ItemOrdemServico representa um item de estoque utilizado em uma ordem de serviço
type ItemOrdemServico struct {
	ID             uint         `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	OrdemServicoID uint         `json:"ordemServicoId" gorm:"not null;index" binding:"required"`
	OrdemServico   OrdemServico `json:"-" gorm:"foreignKey:OrdemServicoID"`
	EstoqueID      uint         `json:"estoqueId" gorm:"not null;index" binding:"required"`
	Item           Estoque      `json:"item,omitempty" gorm:"foreignKey:EstoqueID"`
	Quantidade     int          `json:"quantidade" gorm:"not null;default:1" binding:"required,min=1"`
	ValorUnitario  float64      `json:"valorUnitario" gorm:"type:decimal(10,2);not null" binding:"required"`
	ValorTotal     float64      `json:"valorTotal" gorm:"type:decimal(10,2);not null"`
	CreatedAt      time.Time    `json:"criadoEm" gorm:"autoCreateTime"`
	UpdatedAt      time.Time    `json:"atualizadoEm" gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela para OrdemServico
func (OrdemServico) TableName() string {
	return "ordens_servico"
}

// TableName especifica o nome da tabela para ItemOrdemServico
func (ItemOrdemServico) TableName() string {
	return "itens_ordem_servico"
}

// BeforeCreate gera um número de OS automático
func (os *OrdemServico) BeforeCreate(tx *gorm.DB) error {
	if os.NumeroOS == "" {
		// Gera um número de OS baseado na data e um contador
		ano, mes, dia := time.Now().Date()
		var contador int64
		tx.Model(&OrdemServico{}).Count(&contador)
		os.NumeroOS = fmt.Sprintf("OS%d%02d%02d-%04d", ano, int(mes), dia, contador+1)
	}

	// Inicializa a data de entrada se não for fornecida
	if os.DataEntrada.IsZero() {
		os.DataEntrada = time.Now()
	}

	// Define o status padrão
	if os.Status == "" {
		os.Status = "aberta"
	}

	return nil
}

// BeforeSave calcula o valor total
func (os *OrdemServico) BeforeSave(tx *gorm.DB) error {
	os.ValorTotal = os.ValorPecas + os.ValorServico - os.ValorDesconto
	return nil
}

// BeforeSave calcula o valor total do item
func (item *ItemOrdemServico) BeforeSave(tx *gorm.DB) error {
	item.ValorTotal = float64(item.Quantidade) * item.ValorUnitario
	return nil
}

// OrdemServicoDTO é um DTO para retornar uma ordem de serviço completa com seus relacionamentos
type OrdemServicoDTO struct {
	OrdemServico    OrdemServico       `json:"ordemServico"`
	Cliente         Cliente            `json:"cliente"`
	Veiculo         Veiculo            `json:"veiculo"`
	Funcionario     Funcionario        `json:"funcionario,omitempty"`
	ItensUtilizados []ItemOrdemServico `json:"itensUtilizados"`
}
