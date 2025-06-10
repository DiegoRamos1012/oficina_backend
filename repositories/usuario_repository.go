package repositories

import (
	"OficinaMecanica/models"

	"gorm.io/gorm"
)

// UsuarioRepository define a interface para operações de repositório do usuário
type UsuarioRepository interface {
	FindAll() ([]models.Usuario, error)
	FindByID(id uint) (*models.Usuario, error)
	FindByEmail(email string) (*models.Usuario, error)
	Create(usuario *models.Usuario) error
	Update(usuario *models.Usuario) error
	Delete(id uint) error
}

// UsuarioRepositoryImpl implementa a interface UsuarioRepository
type UsuarioRepositoryImpl struct {
	db *gorm.DB
}

// NewUsuarioRepository cria uma nova instância de UsuarioRepository
func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &UsuarioRepositoryImpl{db: db}
}

// FindAll busca todos os usuários
func (r *UsuarioRepositoryImpl) FindAll() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	result := r.db.Find(&usuarios)
	return usuarios, result.Error
}

// FindByID busca um usuário pelo ID
func (r *UsuarioRepositoryImpl) FindByID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	result := r.db.First(&usuario, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &usuario, nil
}

// FindByEmail busca um usuário pelo email
func (r *UsuarioRepositoryImpl) FindByEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	result := r.db.Where("email = ?", email).First(&usuario)
	if result.Error != nil {
		return nil, result.Error
	}
	return &usuario, nil
}

// Create cria um novo usuário
func (r *UsuarioRepositoryImpl) Create(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

// Update atualiza um usuário existente
func (r *UsuarioRepositoryImpl) Update(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

// Delete remove um usuário pelo ID (soft delete)
func (r *UsuarioRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}
