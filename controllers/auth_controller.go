package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
	"OficinaMecanica/utils"
)

type AuthController struct {
	usuarioService services.UsuarioService
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

type RegisterRequest struct {
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required,min=6"`
	Cargo string `json:"cargo"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Implementação simplificada - em um projeto real, buscar o usuário no banco
	// e verificar a senha
	usuario, err := c.usuarioService.BuscarPorEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(req.Senha)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gerar token
	token, err := utils.GerarToken(*usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Atualizar último login
	c.usuarioService.AtualizarUltimoLogin(usuario.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"usuario": gin.H{
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
		Cargo: req.Cargo,
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
