package services

import (
	"errors"
	"time"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type UsuarioService interface {
	BuscarTodos() ([]models.Usuario, error)
	BuscarPorID(id uint) (*models.Usuario, error)
	BuscarPorEmail(email string) (*models.Usuario, error)
	Criar(usuario *models.Usuario) (*models.Usuario, error)
	Atualizar(usuario *models.Usuario) (*models.Usuario, error)
	Deletar(id uint) error
	AtualizarUltimoLogin(id uint) error
	ValidarCredenciais(email, senha string) (*models.Usuario, error)
}

type UsuarioServiceImpl struct {
	usuarioRepo repositories.UsuarioRepository
}

func NewUsuarioService(usuarioRepo repositories.UsuarioRepository) UsuarioService {
	return &UsuarioServiceImpl{
		usuarioRepo: usuarioRepo,
	}
}

func (s *UsuarioServiceImpl) BuscarTodos() ([]models.Usuario, error) {
	return s.usuarioRepo.FindAll()
}

func (s *UsuarioServiceImpl) BuscarPorID(id uint) (*models.Usuario, error) {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	return usuario, nil
}

func (s *UsuarioServiceImpl) BuscarPorEmail(email string) (*models.Usuario, error) {
	if email == "" {
		return nil, errors.New("email não pode ser vazio")
	}

	usuario, err := s.usuarioRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	return usuario, nil
}

func (s *UsuarioServiceImpl) Criar(usuario *models.Usuario) (*models.Usuario, error) {
	// Validações básicas
	if usuario.Nome == "" {
		return nil, errors.New("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return nil, errors.New("email do usuário é obrigatório")
	}

	if usuario.Senha == "" {
		return nil, errors.New("senha do usuário é obrigatória")
	}

	// Verificar se já existe usuário com o mesmo email
	existente, err := s.usuarioRepo.FindByEmail(usuario.Email)
	if err == nil && existente != nil {
		return nil, errors.New("já existe um usuário com este email")
	}

	// Senha será criptografada pelo hook BeforeCreate do modelo

	err = s.usuarioRepo.Create(usuario)
	if err != nil {
		return nil, errors.New("erro ao criar usuário")
	}

	return usuario, nil
}

func (s *UsuarioServiceImpl) Atualizar(usuario *models.Usuario) (*models.Usuario, error) {
	// Verificar se o usuário existe
	_, err := s.usuarioRepo.FindByID(usuario.ID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Validações básicas
	if usuario.Nome == "" {
		return nil, errors.New("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return nil, errors.New("email do usuário é obrigatório")
	}

	// Verificar se já existe outro usuário com o mesmo email
	existente, err := s.usuarioRepo.FindByEmail(usuario.Email)
	if err == nil && existente != nil && existente.ID != usuario.ID {
		return nil, errors.New("já existe outro usuário com este email")
	}

	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return nil, errors.New("erro ao atualizar usuário")
	}

	return usuario, nil
}

func (s *UsuarioServiceImpl) Deletar(id uint) error {
	// Verificar se o usuário existe
	_, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	err = s.usuarioRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar usuário")
	}

	return nil
}

func (s *UsuarioServiceImpl) AtualizarUltimoLogin(id uint) error {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Atualizar último login
	now := time.Now()
	usuario.UltimoLogin = &now

	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return errors.New("erro ao atualizar último login")
	}

	return nil
}

func (s *UsuarioServiceImpl) ValidarCredenciais(email, senha string) (*models.Usuario, error) {
	if email == "" || senha == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	usuario, err := s.usuarioRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// A verificação da senha será feita no controller usando bcrypt

	return usuario, nil
}
