package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
	"OficinaMecanica/utils"
)

type RegisterRequest struct {
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

type AuthController struct {
	usuarioService services.UsuarioService
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// SetUsuarioService permite injetar um serviço de usuário (útil para testes)
func (c *AuthController) SetUsuarioService(service services.UsuarioService) {
	c.usuarioService = service
}

// Login autentica um usuário e retorna um token JWT
func (c *AuthController) Login(ctx *gin.Context) {
	var loginRequest struct {
		Email string `json:"email" binding:"required,email"`
		Senha string `json:"senha" binding:"required"`
	}

	// Faz o binding do JSON da requisição
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Busca o usuário pelo email
	usuario, err := c.usuarioService.BuscarPorEmail(loginRequest.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Verifica se a senha está correta
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(loginRequest.Senha))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gera o token JWT
	token, err := utils.GerarToken(*usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Atualiza o timestamp de último login
	go c.usuarioService.AtualizarUltimoLogin(usuario.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    usuario.ID,
			"nome":  usuario.Nome,
			"email": usuario.Email,
			"cargo": usuario.Cargo,
		},
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Adicione este log:
		fmt.Printf("Erro ao fazer bind do JSON no Register: %+v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verificar se o email já está em uso
	_, err := c.usuarioService.BuscarPorEmail(req.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email já está em uso"})
		return
	}

	// Criar novo usuário
	usuario := models.Usuario{
		Nome:  req.Nome,
		Email: req.Email,
		Senha: req.Senha, // O hash será gerado pelo hook BeforeCreate
	}

	// Salvar usuário
	usuarioCriado, err := c.usuarioService.Criar(&usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar usuário"})
		return
	}

	// Gerar token
	token, err := utils.GerarToken(*usuarioCriado)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token": token,
		"usuario": gin.H{
			"id":    usuarioCriado.ID,
			"nome":  usuarioCriado.Nome,
			"email": usuarioCriado.Email,
			"cargo": usuarioCriado.Cargo,
		},
	})
}
