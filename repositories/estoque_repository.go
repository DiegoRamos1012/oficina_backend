package repositories

import (
	"OficinaMecanica/models"

	"gorm.io/gorm"
)

type EstoqueRepository interface {
	FindAll() ([]models.Estoque, error)
	FindByID(id uint) (*models.Estoque, error)
	Create(estoque *models.Estoque) error
	Update(estoque *models.Estoque) error
	Delete(id uint) error
	FindByCategoria(categoria string) ([]models.Estoque, error)
	FindBaixoEstoque() ([]models.Estoque, error)
}

type EstoqueRepositoryImpl struct {
	db *gorm.DB
}

func NewEstoqueRepository(db *gorm.DB) EstoqueRepository {
	return &EstoqueRepositoryImpl{db: db}
}

func (r *EstoqueRepositoryImpl) FindAll() ([]models.Estoque, error) {
	var itens []models.Estoque
	result := r.db.Find(&itens)
	return itens, result.Error
}

func (r *EstoqueRepositoryImpl) FindByID(id uint) (*models.Estoque, error) {
	var item models.Estoque
	result := r.db.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *EstoqueRepositoryImpl) Create(estoque *models.Estoque) error {
	return r.db.Create(estoque).Error
}

func (r *EstoqueRepositoryImpl) Update(estoque *models.Estoque) error {
	return r.db.Save(estoque).Error
}

func (r *EstoqueRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Estoque{}, id).Error
}

func (r *EstoqueRepositoryImpl) FindByCategoria(categoria string) ([]models.Estoque, error) {
	var itens []models.Estoque
	result := r.db.Where("categoria = ?", categoria).Find(&itens)
	return itens, result.Error
}

func (r *EstoqueRepositoryImpl) FindBaixoEstoque() ([]models.Estoque, error) {
	var itens []models.Estoque
	result := r.db.Where("quantidade < estoque_minimo").Find(&itens)
	return itens, result.Error
}
