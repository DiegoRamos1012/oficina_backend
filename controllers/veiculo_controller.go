package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

// VeiculoController gerencia as requisições HTTP relacionadas aos veículos
// Responsável por receber requisições, validar dados e retornar respostas apropriadas
type VeiculoController struct {
	veiculoService services.VeiculoService // Serviço injetado via construtor
}

// NewVeiculoController cria uma nova instância do controlador de veículos
// Implementa o padrão de injeção de dependência para facilitar os testes
func NewVeiculoController(veiculoService services.VeiculoService) *VeiculoController {
	return &VeiculoController{
		veiculoService: veiculoService,
	}
}

// BuscarTodos retorna todos os veículos cadastrados
// Este endpoint não recebe parâmetros e retorna um array com todos os veículos
func (c *VeiculoController) BuscarTodos(ctx *gin.Context) {
	// Solicita ao serviço que busque todos os veículos
	veiculos, err := c.veiculoService.BuscarTodos()
	if err != nil {
		// Em caso de erro, retorna status 500 (erro interno)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar veículos"})
		return
	}

	// Retorna a lista de veículos com status 200 (OK)
	ctx.JSON(http.StatusOK, veiculos)
}

// BuscarPorID retorna um veículo específico pelo ID
// Recebe o ID como parâmetro na URL e retorna os detalhes do veículo correspondente
func (c *VeiculoController) BuscarPorID(ctx *gin.Context) {
	// Converte o parâmetro ID da URL de string para inteiro
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// Retorna erro 400 (bad request) se o ID não for um número válido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca o veículo pelo ID usando o serviço
	veiculo, err := c.veiculoService.BuscarPorID(uint(id))
	if err != nil {
		// Retorna erro 404 (not found) se o veículo não for encontrado
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Veículo não encontrado"})
		return
	}

	// Retorna os detalhes do veículo com status 200 (OK)
	ctx.JSON(http.StatusOK, veiculo)
}

// BuscarPorPlaca retorna um veículo específico pela placa
// Recebe a placa como parâmetro na URL e retorna os detalhes do veículo correspondente
func (c *VeiculoController) BuscarPorPlaca(ctx *gin.Context) {
	// Extrai a placa do parâmetro da URL
	placa := ctx.Param("placa")
	if placa == "" {
		// Retorna erro 400 se a placa não for fornecida
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Placa não informada"})
		return
	}

	// Busca o veículo pela placa usando o serviço
	veiculo, err := c.veiculoService.BuscarPorPlaca(placa)
	if err != nil {
		// Retorna erro 404 se o veículo não for encontrado
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Retorna os detalhes do veículo com status 200 (OK)
	ctx.JSON(http.StatusOK, veiculo)
}

// BuscarPorCliente retorna os veículos de um cliente específico
// Recebe o ID do cliente como parâmetro na URL e retorna a lista de seus veículos
func (c *VeiculoController) BuscarPorCliente(ctx *gin.Context) {
	// Converte o parâmetro ID do cliente da URL de string para inteiro
	clienteID, err := strconv.Atoi(ctx.Param("clienteId"))
	if err != nil {
		// Retorna erro 400 (bad request) se o ID não for um número válido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	// Busca os veículos pelo ID do cliente usando o serviço
	veiculos, err := c.veiculoService.BuscarPorClienteID(uint(clienteID))
	if err != nil {
		// Retorna erro 404 se o cliente não for encontrado ou não tiver veículos
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Retorna a lista de veículos do cliente com status 200 (OK)
	ctx.JSON(http.StatusOK, veiculos)
}

// Criar adiciona um novo veículo
// Recebe os dados do veículo no corpo da requisição em formato JSON
func (c *VeiculoController) Criar(ctx *gin.Context) {
	var veiculo models.Veiculo

	// Faz o binding do JSON recebido para o modelo de Veiculo
	// Isso também aplica as validações definidas nas tags "binding" do modelo
	if err := ctx.ShouldBindJSON(&veiculo); err != nil {
		// Retorna erro 400 (bad request) se os dados não forem válidos
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Solicita ao serviço que crie o novo veículo
	veiculoCriado, err := c.veiculoService.Criar(&veiculo)
	if err != nil {
		// Retorna erro 400 com a mensagem específica do erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o veículo criado com status 201 (created)
	ctx.JSON(http.StatusCreated, veiculoCriado)
}

// Atualizar modifica um veículo existente
// Recebe o ID do veículo na URL e os novos dados no corpo da requisição
func (c *VeiculoController) Atualizar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var veiculo models.Veiculo

	// Faz o binding do JSON recebido para o modelo
	if err := ctx.ShouldBindJSON(&veiculo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Define o ID do veículo com o valor extraído da URL
	// Isso garante que o veículo correto será atualizado
	veiculo.ID = uint(id)

	// Solicita ao serviço que atualize o veículo
	veiculoAtualizado, err := c.veiculoService.Atualizar(&veiculo)
	if err != nil {
		// Se houver erro na atualização, retorna status 400 com a mensagem de erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o veículo atualizado com status 200 (OK)
	ctx.JSON(http.StatusOK, veiculoAtualizado)
}

// Deletar remove um veículo
// Recebe o ID do veículo a ser removido como parâmetro na URL
func (c *VeiculoController) Deletar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Solicita ao serviço que remova o veículo
	err = c.veiculoService.Deletar(uint(id))
	if err != nil {
		// Se o veículo não existir ou houver outro erro, retorna status 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna status 204 (no content) indicando sucesso sem conteúdo de retorno
	ctx.Status(http.StatusNoContent)
}
