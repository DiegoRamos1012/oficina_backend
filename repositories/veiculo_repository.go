package repositories

import (
	"OficinaMecanica/models"

	"gorm.io/gorm"
)

type VeiculoRepository interface {
	FindAll() ([]models.Veiculo, error)
	FindByID(id uint) (*models.Veiculo, error)
	Create(veiculo *models.Veiculo) error
	Update(veiculo *models.Veiculo) error
	Delete(id uint) error
	FindByPlaca(placa string) (*models.Veiculo, error)
	FindByClienteID(clienteID uint) ([]models.Veiculo, error)
}

type VeiculoRepositoryImpl struct {
	db *gorm.DB
}

func NewVeiculoRepository(db *gorm.DB) VeiculoRepository {
	return &VeiculoRepositoryImpl{db: db}
}

func (r *VeiculoRepositoryImpl) FindAll() ([]models.Veiculo, error) {
	var veiculos []models.Veiculo
	result := r.db.Find(&veiculos)
	return veiculos, result.Error
}

func (r *VeiculoRepositoryImpl) FindByID(id uint) (*models.Veiculo, error) {
	var veiculo models.Veiculo
	result := r.db.First(&veiculo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &veiculo, nil
}

func (r *VeiculoRepositoryImpl) Create(veiculo *models.Veiculo) error {
	return r.db.Create(veiculo).Error
}

func (r *VeiculoRepositoryImpl) Update(veiculo *models.Veiculo) error {
	return r.db.Save(veiculo).Error
}

func (r *VeiculoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Veiculo{}, id).Error
}

func (r *VeiculoRepositoryImpl) FindByPlaca(placa string) (*models.Veiculo, error) {
	var veiculo models.Veiculo
	result := r.db.Where("placa = ?", placa).First(&veiculo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &veiculo, nil
}

func (r *VeiculoRepositoryImpl) FindByClienteID(clienteID uint) ([]models.Veiculo, error) {
	var veiculos []models.Veiculo
	result := r.db.Where("cliente_id = ?", clienteID).Find(&veiculos)
	return veiculos, result.Error
}
