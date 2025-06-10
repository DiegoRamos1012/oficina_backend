package repositories

import (
	"OficinaMecanica/models"

	"gorm.io/gorm"
)

type FuncionarioRepository interface {
	FindAll() ([]models.Funcionario, error)
	FindByID(id uint) (*models.Funcionario, error)
	Create(funcionario *models.Funcionario) error
	Update(funcionario *models.Funcionario) error
	Delete(id uint) error
	FindByCPF(cpf string) (*models.Funcionario, error)
	FindByCargo(cargo string) ([]models.Funcionario, error)
}

type FuncionarioRepositoryImpl struct {
	db *gorm.DB
}

func NewFuncionarioRepository(db *gorm.DB) FuncionarioRepository {
	return &FuncionarioRepositoryImpl{db: db}
}

func (r *FuncionarioRepositoryImpl) FindAll() ([]models.Funcionario, error) {
	var funcionarios []models.Funcionario
	result := r.db.Find(&funcionarios)
	return funcionarios, result.Error
}

func (r *FuncionarioRepositoryImpl) FindByID(id uint) (*models.Funcionario, error) {
	var funcionario models.Funcionario
	result := r.db.First(&funcionario, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &funcionario, nil
}

func (r *FuncionarioRepositoryImpl) Create(funcionario *models.Funcionario) error {
	return r.db.Create(funcionario).Error
}

func (r *FuncionarioRepositoryImpl) Update(funcionario *models.Funcionario) error {
	return r.db.Save(funcionario).Error
}

func (r *FuncionarioRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Funcionario{}, id).Error
}

func (r *FuncionarioRepositoryImpl) FindByCPF(cpf string) (*models.Funcionario, error) {
	var funcionario models.Funcionario
	result := r.db.Where("cpf = ?", cpf).First(&funcionario)
	if result.Error != nil {
		return nil, result.Error
	}
	return &funcionario, nil
}

func (r *FuncionarioRepositoryImpl) FindByCargo(cargo string) ([]models.Funcionario, error) {
	var funcionarios []models.Funcionario
	result := r.db.Where("cargo = ?", cargo).Find(&funcionarios)
	return funcionarios, result.Error
}
