// filepath: d:\Users\santos.diego\Oficina Mecânica\oficina_backend\controllers\cliente_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

// ClienteController gerencia as requisições HTTP relacionadas aos clientes
// Responsável por receber requisições, validar dados e retornar respostas apropriadas
type ClienteController struct {
	clienteService services.ClienteService // Serviço injetado via construtor
}

// NewClienteController cria uma nova instância do controlador de clientes
// Implementa o padrão de injeção de dependência para facilitar os testes
func NewClienteController(clienteService services.ClienteService) *ClienteController {
	return &ClienteController{
		clienteService: clienteService,
	}
}

// BuscarTodos retorna todos os clientes cadastrados
// Este endpoint não recebe parâmetros e retorna um array com todos os clientes
func (c *ClienteController) BuscarTodos(ctx *gin.Context) {
	// Solicita ao serviço que busque todos os clientes
	clientes, err := c.clienteService.BuscarTodos()
	if err != nil {
		// Em caso de erro, retorna status 500 (erro interno)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	// Retorna a lista de clientes com status 200 (OK)
	ctx.JSON(http.StatusOK, clientes)
}

// BuscarPorID retorna um cliente específico pelo ID
// Recebe o ID como parâmetro na URL e retorna os detalhes do cliente correspondente
func (c *ClienteController) BuscarPorID(ctx *gin.Context) {
	// Converte o parâmetro ID da URL de string para inteiro
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// Retorna erro 400 (bad request) se o ID não for um número válido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca o cliente pelo ID usando o serviço
	cliente, err := c.clienteService.BuscarPorID(uint(id))
	if err != nil {
		// Retorna erro 404 (not found) se o cliente não for encontrado
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	// Retorna os detalhes do cliente com status 200 (OK)
	ctx.JSON(http.StatusOK, cliente)
}

// Criar adiciona um novo cliente
// Recebe os dados do cliente no corpo da requisição em formato JSON
func (c *ClienteController) Criar(ctx *gin.Context) {
	var cliente models.Cliente

	// Faz o binding do JSON recebido para o modelo de Cliente
	// Isso também aplica as validações definidas nas tags "binding" do modelo
	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		// Retorna erro 400 (bad request) se os dados não forem válidos
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Solicita ao serviço que crie o novo cliente
	clienteCriado, err := c.clienteService.Criar(&cliente)
	if err != nil {
		// Retorna erro 400 com a mensagem específica do erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o cliente criado com status 201 (created)
	ctx.JSON(http.StatusCreated, clienteCriado)
}

// Atualizar modifica um cliente existente
// Recebe os dados atualizados do cliente no corpo da requisição
func (c *ClienteController) Atualizar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cliente models.Cliente

	// Faz o binding do JSON recebido para o modelo
	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Define o ID do cliente com o valor extraído da URL
	// Isso garante que o cliente correto será atualizado
	cliente.ID = uint(id)

	// Solicita ao serviço que atualize o cliente
	clienteAtualizado, err := c.clienteService.Atualizar(&cliente)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o cliente atualizado com status 200 (OK)
	ctx.JSON(http.StatusOK, clienteAtualizado)
}

// Deletar remove um cliente
// Recebe o ID do cliente a ser removido como parâmetro na URL
func (c *ClienteController) Deletar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Solicita ao serviço que remova o cliente
	err = c.clienteService.Deletar(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna status 204 (no content) indicando sucesso sem conteúdo de retorno
	ctx.Status(http.StatusNoContent)
}

// BuscarComVeiculos retorna um cliente junto com seus veículos registrados
// Este endpoint é útil para visualizar todos os veículos de um cliente de uma vez
func (c *ClienteController) BuscarComVeiculos(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca o cliente com seus veículos usando o serviço
	clienteDTO, err := c.clienteService.BuscarClienteComVeiculos(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Retorna o cliente com seus veículos
	ctx.JSON(http.StatusOK, clienteDTO)
}
