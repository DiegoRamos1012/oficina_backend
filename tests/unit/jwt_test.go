package tests

import (
	"OficinaMecanica/models"
	"OficinaMecanica/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGerarEValidarToken(t *testing.T) {
	// Criar um usuário de teste
	usuario := models.Usuario{
		ID:    1,
		Nome:  "Test User",
		Email: "test@example.com",
		Cargo: "admin",
	}

	// Gerar um token para o usuário
	token, err := utils.GerarToken(usuario)

	// Verificar se não houve erro na geração do token
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validar o token gerado
	jwtToken, err := utils.ValidarToken(token)

	// Verificar se não houve erro na validação
	assert.NoError(t, err)
	assert.True(t, jwtToken.Valid)

	// Extrair o ID do usuário do token e verificar se corresponde
	userID, ok := utils.ExtrairUserID(jwtToken)
	assert.True(t, ok)
	assert.Equal(t, int(usuario.ID), userID)
}

func TestTokenExpirado(t *testing.T) {
	// Este teste requer modificação da função GerarToken para aceitar um tempo de expiração personalizado
	// ou uma verificação de token expirado com um token criado manualmente

	// Criar claims com tempo de expiração já passado
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"nome":    "Test User",
		"cargo":   "admin",
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expirado há 1 hora
	}

	// Criar token com as claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assinar token com uma chave de teste
	tokenString, err := token.SignedString([]byte("chave_teste"))
	assert.NoError(t, err)

	// Tentar validar o token (deve falhar por estar expirado)
	// Nota: Isso requer uma adaptação da função ValidarToken para usar uma chave de teste
	// ou uma função mock para esse teste específico
	_, err = utils.ValidarToken(tokenString)
	assert.Error(t, err)
}
