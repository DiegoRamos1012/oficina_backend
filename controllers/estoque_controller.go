package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/models"
	"OficinaMecanica/services"
)

// EstoqueController gerencia as requisições HTTP relacionadas ao estoque
type EstoqueController struct {
	estoqueService services.EstoqueService
}

// NewEstoqueController cria uma nova instância do controlador de estoque
func NewEstoqueController(estoqueService services.EstoqueService) *EstoqueController {
	return &EstoqueController{
		estoqueService: estoqueService,
	}
}

// BuscarTodos retorna todos os itens do estoque
// @Summary Listar todos os itens do estoque
// @Description Retorna uma lista com todos os itens cadastrados no estoque
// @Tags estoque
// @Produce json
// @Success 200 {array} models.Estoque
// @Failure 500 {object} map[string]string "Erro ao buscar itens do estoque"
// @Router /estoque [get]
func (c *EstoqueController) BuscarTodos(ctx *gin.Context) {
	itens, err := c.estoqueService.BuscarTodos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens do estoque"})
		return
	}

	ctx.JSON(http.StatusOK, itens)
}

// BuscarPorID retorna um item específico do estoque pelo ID
// @Summary Buscar item por ID
// @Description Retorna um item do estoque pelo seu ID
// @Tags estoque
// @Produce json
// @Param id path int true "ID do item"
// @Success 200 {object} models.Estoque
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Item não encontrado"
// @Router /estoque/{id} [get]
func (c *EstoqueController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	item, err := c.estoqueService.BuscarPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// Criar adiciona um novo item ao estoque
// @Summary Adicionar novo item ao estoque
// @Description Cadastra um novo item no estoque
// @Tags estoque
// @Accept json
// @Produce json
// @Param item body models.Estoque true "Dados do item"
// @Success 201 {object} models.Estoque
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro ao criar item"
// @Router /estoque [post]
func (c *EstoqueController) Criar(ctx *gin.Context) {
	var item models.Estoque

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	itemCriado, err := c.estoqueService.Criar(&item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, itemCriado)
}

// Atualizar modifica um item existente no estoque
// @Summary Atualizar item do estoque
// @Description Atualiza os dados de um item existente no estoque
// @Tags estoque
// @Accept json
// @Produce json
// @Param id path int true "ID do item"
// @Param item body models.Estoque true "Novos dados do item"
// @Success 200 {object} models.Estoque
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 404 {object} map[string]string "Item não encontrado"
// @Failure 500 {object} map[string]string "Erro ao atualizar item"
// @Router /estoque/{id} [put]
func (c *EstoqueController) Atualizar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var item models.Estoque
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	item.ID = uint(id)
	itemAtualizado, err := c.estoqueService.Atualizar(&item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, itemAtualizado)
}

// Deletar remove um item do estoque
// @Summary Remover item do estoque
// @Description Remove um item do estoque pelo ID
// @Tags estoque
// @Param id path int true "ID do item"
// @Success 204 "Item removido com sucesso"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Item não encontrado"
// @Failure 500 {object} map[string]string "Erro ao remover item"
// @Router /estoque/{id} [delete]
func (c *EstoqueController) Deletar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.estoqueService.Deletar(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// BuscarPorCategoria retorna itens do estoque filtrados por categoria
// @Summary Buscar itens por categoria
// @Description Retorna todos os itens do estoque de uma categoria específica
// @Tags estoque
// @Produce json
// @Param categoria path string true "Nome da categoria"
// @Success 200 {array} models.Estoque
// @Failure 400 {object} map[string]string "Categoria inválida"
// @Failure 500 {object} map[string]string "Erro ao buscar itens"
// @Router /estoque/categoria/{categoria} [get]
func (c *EstoqueController) BuscarPorCategoria(ctx *gin.Context) {
	categoria := ctx.Param("categoria")
	if categoria == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Categoria não especificada"})
		return
	}

	itens, err := c.estoqueService.BuscarPorCategoria(categoria)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, itens)
}

// BuscarBaixoEstoque retorna itens com estoque abaixo do mínimo
// @Summary Listar itens com estoque baixo
// @Description Retorna todos os itens que estão com quantidade abaixo do estoque mínimo
// @Tags estoque
// @Produce json
// @Success 200 {array} models.Estoque
// @Failure 500 {object} map[string]string "Erro ao buscar itens com estoque baixo"
// @Router /estoque/baixo-estoque [get]
func (c *EstoqueController) BuscarBaixoEstoque(ctx *gin.Context) {
	itens, err := c.estoqueService.BuscarBaixoEstoque()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, itens)
}

// AtualizarQuantidade incrementa ou decrementa a quantidade de um item no estoque
// @Summary Atualizar quantidade em estoque
// @Description Atualiza a quantidade disponível de um item no estoque
// @Tags estoque
// @Accept json
// @Produce json
// @Param id path int true "ID do item"
// @Param dados body map[string]int true "Dados da atualização"
// @Success 200 {object} models.Estoque
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 404 {object} map[string]string "Item não encontrado"
// @Router /estoque/{id}/quantidade [patch]
func (c *EstoqueController) AtualizarQuantidade(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var dados struct {
		Quantidade int `json:"quantidade" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&dados); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	item, err := c.estoqueService.BuscarPorID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	// Atualiza a quantidade
	item.Quantidade = dados.Quantidade

	itemAtualizado, err := c.estoqueService.Atualizar(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, itemAtualizado)
}

// Estrutura para controle de limites de estoque
// Salva em arquivo local: controle_estoque_config.json

type ControleEstoqueConfig struct {
	LimiteBaixo int `json:"limite_baixo"`
	LimiteMedio int `json:"limite_medio"`
}

const controleEstoqueConfigFile = "controle_estoque_config.json"

// GET /estoque/controle-estoque
func (c *EstoqueController) BuscarControleEstoque(ctx *gin.Context) {
	config, err := lerControleEstoqueConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler controle de estoque: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

// POST /estoque/controle-estoque
func (c *EstoqueController) SalvarControleEstoque(ctx *gin.Context) {
	var config ControleEstoqueConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}
	if err := gravarControleEstoqueConfig(config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar controle de estoque: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

// Funções utilitárias para ler/gravar controle de estoque
func lerControleEstoqueConfig() (ControleEstoqueConfig, error) {
	var config ControleEstoqueConfig
	file, err := os.Open(controleEstoqueConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Retorna valores padrão se não existir
			return ControleEstoqueConfig{LimiteBaixo: 10, LimiteMedio: 20}, nil
		}
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}

func gravarControleEstoqueConfig(config ControleEstoqueConfig) error {
	file, err := os.Create(controleEstoqueConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}
