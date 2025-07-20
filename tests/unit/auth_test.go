package tests

import (
	"OficinaMecanica/controllers"
	"OficinaMecanica/models"
	"OficinaMecanica/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"OficinaMecanica/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock do repositório de usuário
type MockUsuarioRepository struct {
	mock.Mock
}

func (m *MockUsuarioRepository) FindAll() ([]models.Usuario, error) {
	args := m.Called()
	return args.Get(0).([]models.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByID(id uint) (*models.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByEmail(email string) (*models.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) Create(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) Update(usuario *models.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	// Configurar o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	// Mock da função de geração de token
	originalGerarTokenFn := utils.GerarTokenFn
	defer func() { utils.GerarTokenFn = originalGerarTokenFn }() // Restaura após o teste

	utils.GerarTokenFn = func(usuario models.Usuario) (string, error) {
		return "mock-jwt-token", nil
	}

	// Criar um mock do repositório
	mockRepo := new(MockUsuarioRepository)

	// Criar o serviço com o mock do repositório

	// Criar o controlador de autenticação sem argumentos
	authController := controllers.NewAuthController()

	// Criar um roteador Gin para o teste
	router := gin.Default()
	router.POST("/login", authController.Login)

	// Configurar expectativas do mock - quando buscar por email "test@example.com"
	// deve retornar um usuário com senha hash correspondente a "password123"
	hashSenha, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUsuario := &models.Usuario{
		ID:    1,
		Nome:  "Test User",
		Email: "test@example.com",
		Senha: string(hashSenha),
		Cargo: "user",
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(mockUsuario, nil)
	// Adicionar expectativa para FindByID - é chamado por AtualizarUltimoLogin
	mockRepo.On("FindByID", uint(1)).Return(mockUsuario, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Usuario")).Return(nil)

	// Criar uma requisição de login
	loginJSON := `{"email":"test@example.com","senha":"password123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	// Executar a requisição
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusOK, resp.Code)

	// Verificar a resposta
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)

	// Verificar se retornou um token
	assert.Contains(t, response, "token")

	// Verificar se o mock foi chamado conforme esperado
	mockRepo.AssertExpectations(t)
}
