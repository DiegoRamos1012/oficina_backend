package services

import (
	"errors"
	"time"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

// UsuarioService define a interface para operações relacionadas a usuários
// Esta interface permite que possamos substituir a implementação real por mocks em testes
type UsuarioService interface {
	BuscarTodos() ([]models.Usuario, error)                          // Retorna todos os usuários cadastrados
	BuscarPorID(id uint) (*models.Usuario, error)                    // Busca um usuário pelo ID
	BuscarPorEmail(email string) (*models.Usuario, error)            // Busca um usuário pelo email
	Criar(usuario *models.Usuario) (*models.Usuario, error)          // Cria um novo usuário
	Atualizar(usuario *models.Usuario) (*models.Usuario, error)      // Atualiza os dados de um usuário
	Deletar(id uint) error                                           // Remove um usuário (soft delete)
	AtualizarUltimoLogin(id uint) error                              // Atualiza o timestamp de último login
	ValidarCredenciais(email, senha string) (*models.Usuario, error) // Verifica se as credenciais são válidas
	AlterarSenha(id uint, senhaAtual, novaSenha string) error        // Altera a senha do usuário
	AlterarStatus(id uint, ativo bool) error                         // Ativa ou desativa um usuário
	AtualizarAvatar(id uint, avatarPath string) error                // Atualiza o avatar do usuário
}

// UsuarioServiceImpl implementa a interface UsuarioService
// Contém a instância do repositório que será usada para operações de persistência
type UsuarioServiceImpl struct {
	usuarioRepo repositories.UsuarioRepository // Repositório de usuários injetado
}

// NewUsuarioService cria uma nova instância do serviço de usuários
// Implementa o padrão de injeção de dependência
func NewUsuarioService(usuarioRepo repositories.UsuarioRepository) UsuarioService {
	return &UsuarioServiceImpl{
		usuarioRepo: usuarioRepo,
	}
}

// BuscarTodos retorna todos os usuários cadastrados no sistema
func (s *UsuarioServiceImpl) BuscarTodos() ([]models.Usuario, error) {
	return s.usuarioRepo.FindAll()
}

// BuscarPorID busca um usuário pelo seu ID
// Retorna erro se o usuário não for encontrado
func (s *UsuarioServiceImpl) BuscarPorID(id uint) (*models.Usuario, error) {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	return usuario, nil
}

// BuscarPorEmail busca um usuário pelo seu email
// Retorna erro se o email for vazio ou se o usuário não for encontrado
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

// Criar registra um novo usuário no sistema
// Realiza validações dos dados antes de persistir
func (s *UsuarioServiceImpl) Criar(usuario *models.Usuario) (*models.Usuario, error) {
	// Validações básicas dos campos obrigatórios
	if usuario.Nome == "" {
		return nil, errors.New("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return nil, errors.New("email do usuário é obrigatório")
	}

	if usuario.Senha == "" {
		return nil, errors.New("senha do usuário é obrigatória")
	}

	// Verifica se já existe usuário com o mesmo email (unicidade)
	existente, err := s.usuarioRepo.FindByEmail(usuario.Email)
	if err == nil && existente != nil {
		return nil, errors.New("já existe um usuário com este email")
	}

	// Nota: a senha será criptografada automaticamente pelo hook BeforeCreate do modelo

	// Persiste o novo usuário
	err = s.usuarioRepo.Create(usuario)
	if err != nil {
		return nil, errors.New("erro ao criar usuário")
	}

	return usuario, nil
}

// Atualizar modifica os dados de um usuário existente
// Realiza validações e verifica a existência do usuário antes de atualizar
func (s *UsuarioServiceImpl) Atualizar(usuario *models.Usuario) (*models.Usuario, error) {
	// Verifica se o usuário existe
	_, err := s.usuarioRepo.FindByID(usuario.ID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Validações básicas dos campos obrigatórios
	if usuario.Nome == "" {
		return nil, errors.New("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return nil, errors.New("email do usuário é obrigatório")
	}

	// Verifica se já existe outro usuário com o mesmo email (unicidade)
	existente, err := s.usuarioRepo.FindByEmail(usuario.Email)
	if err == nil && existente != nil && existente.ID != usuario.ID {
		return nil, errors.New("já existe outro usuário com este email")
	}

	// Persiste as alterações
	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return nil, errors.New("erro ao atualizar usuário")
	}

	return usuario, nil
}

// Deletar remove um usuário do sistema (soft delete via GORM)
// Verifica a existência do usuário antes de remover
func (s *UsuarioServiceImpl) Deletar(id uint) error {
	// Verifica se o usuário existe
	_, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Remove o usuário
	err = s.usuarioRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar usuário")
	}

	return nil
}

// AtualizarUltimoLogin registra o momento do último acesso do usuário
// Útil para controle de sessão e análise de atividade
func (s *UsuarioServiceImpl) AtualizarUltimoLogin(id uint) error {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Define o timestamp atual como último login
	now := time.Now()
	usuario.UltimoLogin = &now

	// Persiste a alteração
	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return errors.New("erro ao atualizar último login")
	}

	return nil
}

// ValidarCredenciais verifica se o email e senha fornecidos são válidos
// Retorna o usuário se as credenciais estiverem corretas
func (s *UsuarioServiceImpl) ValidarCredenciais(email, senha string) (*models.Usuario, error) {
	if email == "" || senha == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	usuario, err := s.usuarioRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Nota: a verificação da senha será realizada no controller usando bcrypt
	// O service apenas retorna o usuário encontrado pelo email

	return usuario, nil
}

// AlterarSenha troca a senha de um usuário
// Verifica a senha atual antes de permitir a alteração
func (s *UsuarioServiceImpl) AlterarSenha(id uint, senhaAtual, novaSenha string) error {
	if senhaAtual == "" || novaSenha == "" {
		return errors.New("senha atual e nova senha são obrigatórias")
	}

	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Verifica se a senha atual está correta usando o método CompareSenha do modelo
	// Este método deve comparar a senha fornecida com o hash armazenado
	if !usuario.CompareSenha(senhaAtual) {
		return errors.New("senha atual incorreta")
	}

	// Define a nova senha - será criptografada pelo hook BeforeUpdate do modelo
	usuario.Senha = novaSenha

	// Persiste a alteração
	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return errors.New("erro ao atualizar senha")
	}

	return nil
}

// AlterarStatus ativa ou desativa um usuário
// Útil para controle de acesso sem remover o usuário do sistema
func (s *UsuarioServiceImpl) AlterarStatus(id uint, ativo bool) error {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Define o novo status
	usuario.Ativo = ativo

	// Persiste a alteração
	err = s.usuarioRepo.Update(usuario)
	if err != nil {
		return errors.New("erro ao alterar status do usuário")
	}

	return nil
}

// AtualizarAvatar atualiza o avatar de um usuário
// Recebe o ID do usuário e o caminho para a nova imagem do avatar
func (s *UsuarioServiceImpl) AtualizarAvatar(id uint, avatarPath string) error {
	usuario, err := s.usuarioRepo.FindByID(id)
	if err != nil {
		return err
	}
	usuario.Avatar = avatarPath
	return s.usuarioRepo.Update(usuario)
}
