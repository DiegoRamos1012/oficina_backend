package repositories

import (
	"OficinaMecanica/models"

	"gorm.io/gorm"
)

type ClienteRepositoryGorm interface {
	FindAll() ([]models.Cliente, error)
	FindByID(id uint) (*models.Cliente, error)
	Create(cliente *models.Cliente) error
	Update(cliente *models.Cliente) error
	Delete(id uint) error
	FindWithVeiculos(id uint) (*models.Cliente, error)
}

type ClienteRepositoryGormImpl struct {
	db *gorm.DB
}

func NewClienteRepositoryGorm(db *gorm.DB) ClienteRepositoryGorm {
	return &ClienteRepositoryGormImpl{db: db}
}

func (r *ClienteRepositoryGormImpl) FindAll() ([]models.Cliente, error) {
	var clientes []models.Cliente
	result := r.db.Find(&clientes)
	return clientes, result.Error
}

func (r *ClienteRepositoryGormImpl) FindByID(id uint) (*models.Cliente, error) {
	var cliente models.Cliente
	result := r.db.First(&cliente, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cliente, nil
}

func (r *ClienteRepositoryGormImpl) Create(cliente *models.Cliente) error {
	return r.db.Create(cliente).Error
}

func (r *ClienteRepositoryGormImpl) Update(cliente *models.Cliente) error {
	return r.db.Save(cliente).Error
}

func (r *ClienteRepositoryGormImpl) Delete(id uint) error {
	return r.db.Delete(&models.Cliente{}, id).Error
}

func (r *ClienteRepositoryGormImpl) FindWithVeiculos(id uint) (*models.Cliente, error) {
	var cliente models.Cliente
	result := r.db.Preload("Veiculos").First(&cliente, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cliente, nil
}
