package tests

import (
    "OficinaMecanica/middlewares"
    "OficinaMecanica/models"
    "OficinaMecanica/utils"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
    // Configurar o modo de teste do Gin
    gin.SetMode(gin.TestMode)
    
    // Criar um roteador Gin com o middleware de autenticação
    router := gin.New()
    router.Use(middlewares.AuthMiddleware())
    
    // Adicionar uma rota protegida para teste
    router.GET("/protected", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "success"})
    })
    
    // Caso 1: Requisição sem token
    req, _ := http.NewRequest("GET", "/protected", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    
    // Deve retornar status não autorizado
    assert.Equal(t, http.StatusUnauthorized, resp.Code)
    
    // Caso 2: Requisição com token válido
    usuario := models.Usuario{
        ID:    1,
        Nome:  "Test User",
        Email: "test@example.com",
        Cargo: "admin",
    }
    
    token, err := utils.GerarToken(usuario)
    assert.NoError(t, err)
    
    req, _ = http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    
    // Deve retornar status OK
    assert.Equal(t, http.StatusOK, resp.Code)
    
    // Caso 3: Requisição com token inválido
    req, _ = http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer token_invalido")
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    
    // Deve retornar status não autorizado
    assert.Equal(t, http.StatusUnauthorized, resp.Code)
}