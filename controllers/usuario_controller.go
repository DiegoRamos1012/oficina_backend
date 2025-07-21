package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

type UsuarioController struct {
	usuarioService services.UsuarioService
}

func NewUsuarioController(usuarioService services.UsuarioService) *UsuarioController {
	return &UsuarioController{
		usuarioService: usuarioService,
	}
}

func (c *UsuarioController) BuscarTodos(ctx *gin.Context) {
	usuarios, err := c.usuarioService.BuscarTodos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}

	// Remover senhas dos resultados
	usuariosSemSenha := make([]gin.H, len(usuarios))
	for i, u := range usuarios {
		usuariosSemSenha[i] = gin.H{
			"id":               u.ID,
			"nome":             u.Nome,
			"email":            u.Email,
			"cargo":            u.Cargo,
			"ultimo_login":     u.UltimoLogin,
			"data_criacao":     u.CreatedAt,
			"data_atualizacao": u.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, usuariosSemSenha)
}

func (c *UsuarioController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	usuario, err := c.usuarioService.BuscarPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Remover senha do resultado
	ctx.JSON(http.StatusOK, gin.H{
		"id":               usuario.ID,
		"nome":             usuario.Nome,
		"email":            usuario.Email,
		"cargo":            usuario.Cargo,
		"ultimo_login":     usuario.UltimoLogin,
		"data_criacao":     usuario.CreatedAt,
		"data_atualizacao": usuario.UpdatedAt,
	})
}

func (c *UsuarioController) Criar(ctx *gin.Context) {
	var usuario models.Usuario

	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verificar se o email já está em uso
	_, err := c.usuarioService.BuscarPorEmail(usuario.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email já está em uso"})
		return
	}

	usuarioCriado, err := c.usuarioService.Criar(&usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	// Remover senha do resultado
	ctx.JSON(http.StatusCreated, gin.H{
		"id":               usuarioCriado.ID,
		"nome":             usuarioCriado.Nome,
		"email":            usuarioCriado.Email,
		"cargo":            usuarioCriado.Cargo,
		"ultimo_login":     usuarioCriado.UltimoLogin,
		"data_criacao":     usuarioCriado.CreatedAt,
		"data_atualizacao": usuarioCriado.UpdatedAt,
	})
}

func (c *UsuarioController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	type UpdateUsuarioDTO struct {
		Nome  string `json:"nome"`
		Email string `json:"email"`
	}
	var input UpdateUsuarioDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	usuario, err := c.usuarioService.BuscarPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	if input.Nome != "" {
		usuario.Nome = input.Nome
	}
	if input.Email != "" {
		usuario.Email = input.Email
	}

	usuarioAtualizado, err := c.usuarioService.Atualizar(usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":               usuarioAtualizado.ID,
		"nome":             usuarioAtualizado.Nome,
		"email":            usuarioAtualizado.Email,
		"cargo":            usuarioAtualizado.Cargo,
		"ultimo_login":     usuarioAtualizado.UltimoLogin,
		"data_criacao":     usuarioAtualizado.CreatedAt,
		"data_atualizacao": usuarioAtualizado.UpdatedAt,
	})
}

func (c *UsuarioController) AlterarSenha(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	type SenhaRequest struct {
		SenhaAtual string `json:"senha_atual" binding:"required"`
		NovaSenha  string `json:"nova_senha" binding:"required,min=6"`
	}

	var req SenhaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err = c.usuarioService.AlterarSenha(uint(id), req.SenhaAtual, req.NovaSenha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) Desativar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.usuarioService.AlterarStatus(uint(id), false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao desativar usuário"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) Ativar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.usuarioService.AlterarStatus(uint(id), true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ativar usuário"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.usuarioService.Deletar(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar usuário"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UsuarioController) UploadAvatar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}

	// Garante que a pasta exista
	os.MkdirAll("uploads/avatars", os.ModePerm)

	// Busca o usuário para deletar o avatar antigo
	usuario, err := c.usuarioService.BuscarPorID(uint(id))
	if err == nil && usuario.Avatar != "" {
		if _, statErr := os.Stat(usuario.Avatar); statErr == nil {
			_ = os.Remove(usuario.Avatar)
		}
	}

	filename := fmt.Sprintf("uploads/avatars/%d_%s", id, file.Filename)
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
		return
	}

	err = c.usuarioService.AtualizarAvatar(uint(id), filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar avatar"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"avatar": filename})
}
