package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

// FuncionarioController gerencia as requisições HTTP relacionadas aos funcionários
// Responsável por receber requisições, validar dados e retornar respostas apropriadas
type FuncionarioController struct {
	funcionarioService services.FuncionarioService // Serviço injetado via construtor
}

// NewFuncionarioController cria uma nova instância do controlador de funcionários
// Implementa o padrão de injeção de dependência para facilitar os testes
func NewFuncionarioController(funcionarioService services.FuncionarioService) *FuncionarioController {
	return &FuncionarioController{
		funcionarioService: funcionarioService,
	}
}

// BuscarTodos retorna todos os funcionários cadastrados
// Este endpoint não recebe parâmetros e retorna um array com todos os funcionários
func (c *FuncionarioController) BuscarTodos(ctx *gin.Context) {
	// Solicita ao serviço que busque todos os funcionários
	funcionarios, err := c.funcionarioService.BuscarTodos()
	if err != nil {
		// Em caso de erro, retorna status 500 (erro interno)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar funcionários"})
		return
	}

	// Retorna a lista de funcionários com status 200 (OK)
	ctx.JSON(http.StatusOK, funcionarios)
}

// BuscarPorID retorna um funcionário específico pelo ID
// Recebe o ID como parâmetro na URL e retorna os detalhes do funcionário correspondente
func (c *FuncionarioController) BuscarPorID(ctx *gin.Context) {
	// Converte o parâmetro ID da URL de string para inteiro
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// Retorna erro 400 (bad request) se o ID não for um número válido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca o funcionário pelo ID usando o serviço
	funcionario, err := c.funcionarioService.BuscarPorID(uint(id))
	if err != nil {
		// Retorna erro 404 (not found) se o funcionário não for encontrado
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Funcionário não encontrado"})
		return
	}

	// Retorna os detalhes do funcionário com status 200 (OK)
	ctx.JSON(http.StatusOK, funcionario)
}

// BuscarPorCPF retorna um funcionário específico pelo CPF
// Recebe o CPF como parâmetro na URL e retorna os detalhes do funcionário correspondente
func (c *FuncionarioController) BuscarPorCPF(ctx *gin.Context) {
	// Extrai o CPF do parâmetro da URL
	cpf := ctx.Param("cpf")
	if cpf == "" {
		// Retorna erro 400 se o CPF não for fornecido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CPF não informado"})
		return
	}

	// Busca o funcionário pelo CPF usando o serviço
	funcionario, err := c.funcionarioService.BuscarPorCPF(cpf)
	if err != nil {
		// Retorna erro 404 se o funcionário não for encontrado ou CPF inválido
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Retorna os detalhes do funcionário com status 200 (OK)
	ctx.JSON(http.StatusOK, funcionario)
}

// BuscarPorCargo retorna funcionários filtrados pelo cargo
// Recebe o cargo como parâmetro na URL e retorna a lista de funcionários com esse cargo
func (c *FuncionarioController) BuscarPorCargo(ctx *gin.Context) {
	// Extrai o cargo do parâmetro da URL
	cargo := ctx.Param("cargo")
	if cargo == "" {
		// Retorna erro 400 se o cargo não for fornecido
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cargo não informado"})
		return
	}

	// Busca os funcionários pelo cargo usando o serviço
	funcionarios, err := c.funcionarioService.BuscarPorCargo(cargo)
	if err != nil {
		// Retorna erro 400 com a mensagem específica do erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna a lista de funcionários com status 200 (OK)
	ctx.JSON(http.StatusOK, funcionarios)
}

// Criar adiciona um novo funcionário
// Recebe os dados do funcionário no corpo da requisição em formato JSON
func (c *FuncionarioController) Criar(ctx *gin.Context) {
	var funcionario models.Funcionario

	// Faz o binding do JSON recebido para o modelo de Funcionario
	// Isso também aplica as validações definidas nas tags "binding" do modelo
	if err := ctx.ShouldBindJSON(&funcionario); err != nil {
		// Retorna erro 400 (bad request) se os dados não forem válidos
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Solicita ao serviço que crie o novo funcionário
	funcionarioCriado, err := c.funcionarioService.Criar(&funcionario)
	if err != nil {
		// Retorna erro 400 com a mensagem específica do erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o funcionário criado com status 201 (created)
	ctx.JSON(http.StatusCreated, funcionarioCriado)
}

// Atualizar modifica um funcionário existente
// Recebe o ID do funcionário na URL e os novos dados no corpo da requisição
func (c *FuncionarioController) Atualizar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var funcionario models.Funcionario

	// Faz o binding do JSON recebido para o modelo
	if err := ctx.ShouldBindJSON(&funcionario); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Define o ID do funcionário com o valor extraído da URL
	// Isso garante que o funcionário correto será atualizado
	funcionario.ID = uint(id)

	// Solicita ao serviço que atualize o funcionário
	funcionarioAtualizado, err := c.funcionarioService.Atualizar(&funcionario)
	if err != nil {
		// Se houver erro na atualização, retorna status 400 com a mensagem de erro
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna o funcionário atualizado com status 200 (OK)
	ctx.JSON(http.StatusOK, funcionarioAtualizado)
}

// Deletar remove um funcionário
// Recebe o ID do funcionário a ser removido como parâmetro na URL
func (c *FuncionarioController) Deletar(ctx *gin.Context) {
	// Extrai e valida o ID do parâmetro da URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Solicita ao serviço que remova o funcionário
	err = c.funcionarioService.Deletar(uint(id))
	if err != nil {
		// Se o funcionário não existir ou houver outro erro, retorna status 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retorna status 204 (no content) indicando sucesso sem conteúdo de retorno
	ctx.Status(http.StatusNoContent)
}
