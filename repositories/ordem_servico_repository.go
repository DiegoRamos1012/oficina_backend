package repositories

import (
	"OficinaMecanica/models"
	"time"

	"gorm.io/gorm"
)

type OrdemServicoRepository interface {
	FindAll() ([]models.OrdemServico, error)
	FindByID(id uint) (*models.OrdemServico, error)
	Create(os *models.OrdemServico) error
	Update(os *models.OrdemServico) error
	Delete(id uint) error
	FindByClienteID(clienteID uint) ([]models.OrdemServico, error)
	FindByVeiculoID(veiculoID uint) ([]models.OrdemServico, error)
	FindByStatus(status string) ([]models.OrdemServico, error)
	FindByPeriodo(inicio, fim time.Time) ([]models.OrdemServico, error)
	FindByNumeroOS(numeroOS string) (*models.OrdemServico, error)
	AddItem(item *models.ItemOrdemServico) error
	RemoveItem(itemID uint) error
	UpdateItem(item *models.ItemOrdemServico) error
	FindItens(osID uint) ([]models.ItemOrdemServico, error)
}

type OrdemServicoRepositoryImpl struct {
	db *gorm.DB
}

func NewOrdemServicoRepository(db *gorm.DB) OrdemServicoRepository {
	return &OrdemServicoRepositoryImpl{db: db}
}

func (r *OrdemServicoRepositoryImpl) FindAll() ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").Find(&ordens)
	return ordens, result.Error
}

func (r *OrdemServicoRepositoryImpl) FindByID(id uint) (*models.OrdemServico, error) {
	var os models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Preload("ItensUtilizados").Preload("ItensUtilizados.Item").
		First(&os, id)
	return &os, result.Error
}

func (r *OrdemServicoRepositoryImpl) Create(os *models.OrdemServico) error {
	return r.db.Create(os).Error
}

func (r *OrdemServicoRepositoryImpl) Update(os *models.OrdemServico) error {
	return r.db.Save(os).Error
}

func (r *OrdemServicoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.OrdemServico{}, id).Error
}

func (r *OrdemServicoRepositoryImpl) FindByClienteID(clienteID uint) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Where("cliente_id = ?", clienteID).Find(&ordens)
	return ordens, result.Error
}

func (r *OrdemServicoRepositoryImpl) FindByVeiculoID(veiculoID uint) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Where("veiculo_id = ?", veiculoID).Find(&ordens)
	return ordens, result.Error
}

func (r *OrdemServicoRepositoryImpl) FindByStatus(status string) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Where("status = ?", status).Find(&ordens)
	return ordens, result.Error
}

func (r *OrdemServicoRepositoryImpl) FindByPeriodo(inicio, fim time.Time) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Where("data_entrada BETWEEN ? AND ?", inicio, fim).Find(&ordens)
	return ordens, result.Error
}

func (r *OrdemServicoRepositoryImpl) FindByNumeroOS(numeroOS string) (*models.OrdemServico, error) {
	var os models.OrdemServico
	result := r.db.Preload("Veiculo").Preload("Cliente").Preload("Funcionario").
		Preload("ItensUtilizados").Preload("ItensUtilizados.Item").
		Where("numero_os = ?", numeroOS).First(&os)
	return &os, result.Error
}

func (r *OrdemServicoRepositoryImpl) AddItem(item *models.ItemOrdemServico) error {
	return r.db.Create(item).Error
}

func (r *OrdemServicoRepositoryImpl) RemoveItem(itemID uint) error {
	return r.db.Delete(&models.ItemOrdemServico{}, itemID).Error
}

func (r *OrdemServicoRepositoryImpl) UpdateItem(item *models.ItemOrdemServico) error {
	return r.db.Save(item).Error
}

func (r *OrdemServicoRepositoryImpl) FindItens(osID uint) ([]models.ItemOrdemServico, error) {
	var itens []models.ItemOrdemServico
	result := r.db.Preload("Item").Where("ordem_servico_id = ?", osID).Find(&itens)
	return itens, result.Error
}
