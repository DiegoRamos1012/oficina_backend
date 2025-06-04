// filepath: d:\Users\santos.diego\Oficina Mecânica\oficina_backend\controllers\cliente_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

type ClienteController struct {
	clienteService services.ClienteService
}

func NewClienteController(clienteService services.ClienteService) *ClienteController {
	return &ClienteController{
		clienteService: clienteService,
	}
}

func (c *ClienteController) BuscarTodos(ctx *gin.Context) {
	clientes, err := c.clienteService.BuscarTodos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	ctx.JSON(http.StatusOK, clientes)
}

func (c *ClienteController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	cliente, err := c.clienteService.BuscarPorID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, cliente)
}

func (c *ClienteController) Criar(ctx *gin.Context) {
	var cliente models.Cliente

	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	clienteCriado, err := c.clienteService.Criar(cliente)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar cliente"})
		return
	}

	ctx.JSON(http.StatusCreated, clienteCriado)
}
