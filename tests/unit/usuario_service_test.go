package tests

import (
	"OficinaMecanica/models"
	"OficinaMecanica/services"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock do repositório de usuário
type UsuarioRepositoryMock struct {
	mock.Mock
}

func (m *UsuarioRepositoryMock) FindByEmail(email string) (*models.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Usuario), args.Error(1)
}

func (m *UsuarioRepositoryMock) FindByID(id uint) (*models.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Usuario), args.Error(1)
}

func (m *UsuarioRepositoryMock) Create(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *UsuarioRepositoryMock) Update(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *UsuarioRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UsuarioRepositoryMock) FindAll() ([]models.Usuario, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Usuario), args.Error(1)
}

func init() {
	// Configurar viper para ler o arquivo test.env
	viper.SetConfigFile("test.env")
	if err := viper.ReadInConfig(); err != nil {
		// Tentar outro caminho
		viper.SetConfigFile("../../test.env")
		if err := viper.ReadInConfig(); err != nil {
			// Usar configurações padrão para testes
			viper.Set("DB_HOST", "localhost")
			viper.Set("DB_PORT", "3306")
			viper.Set("DB_USER", "root")
			viper.Set("DB_PASSWORD", "A1b2c3d4-")
			viper.Set("DB_NAME", "oficinateste_bd")
		}
	}
}

func TestValidarCredenciais(t *testing.T) {
	// Criar mock do repositório
	mockRepo := new(UsuarioRepositoryMock)

	// Criar o serviço usando o mock
	usuarioService := services.NewUsuarioService(mockRepo)

	// Configurar um usuário de teste com senha "senha123"
	senhaHash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	usuario := &models.Usuario{
		ID:    1,
		Email: "test@example.com",
		Senha: string(senhaHash),
	}

	// Configurar o mock para retornar o usuário quando pesquisado pelo email
	mockRepo.On("FindByEmail", "test@example.com").Return(usuario, nil)
	mockRepo.On("FindByEmail", "wrong@example.com").Return(nil, errors.New("usuário não encontrado"))

	// Caso 1: Credenciais corretas
	usuarioRetornado, err := usuarioService.ValidarCredenciais("test@example.com", "senha123")
	assert.NoError(t, err)
	assert.NotNil(t, usuarioRetornado)
	assert.Equal(t, uint(1), usuarioRetornado.ID)

	// Caso 2: Email incorreto
	usuarioRetornado, err = usuarioService.ValidarCredenciais("wrong@example.com", "senha123")
	assert.Error(t, err)
	assert.Nil(t, usuarioRetornado)

	// Caso 3: Senha incorreta
	usuarioRetornado, err = usuarioService.ValidarCredenciais("test@example.com", "senha_errada")
	assert.Error(t, err)
	assert.Nil(t, usuarioRetornado)

	// Verificar se o mock foi chamado conforme esperado
	mockRepo.AssertExpectations(t)
}

func TestAtualizarUltimoLogin(t *testing.T) {
	// Criar mock do repositório
	mockRepo := new(MockUsuarioRepository)

	// Criar o serviço usando o mock
	usuarioService := services.NewUsuarioService(mockRepo)

	// Criar um usuário para o teste
	usuario := &models.Usuario{
		ID:    1,
		Email: "test@example.com",
	}

	// Configurar o mock para retornar o usuário quando pesquisado pelo ID
	mockRepo.On("FindByID", uint(1)).Return(usuario, nil)

	// Configurar o mock para aceitar a atualização do usuário
	// Nota: Aqui usamos mock.AnythingOfType porque o timestamp será atualizado
	mockRepo.On("Update", mock.AnythingOfType("*models.Usuario")).Return(nil)

	// Executar a atualização do último login
	err := usuarioService.AtualizarUltimoLogin(1)
	assert.NoError(t, err)

	// Verificar se o mock foi chamado conforme esperado
	mockRepo.AssertExpectations(t)
}
