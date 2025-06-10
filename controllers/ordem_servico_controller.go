package controllers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"

    "OficinaMecanica/models"
    "OficinaMecanica/services"
)

// OrdemServicoController gerencia as requisições HTTP relacionadas às ordens de serviço
// Responsável por receber requisições, validar dados e retornar respostas apropriadas
type OrdemServicoController struct {
    osService services.OrdemServicoService // Serviço de ordem de serviço injetado
}

// NewOrdemServicoController cria uma nova instância do controlador de ordens de serviço
// Implementa o padrão de injeção de dependência para facilitar os testes
func NewOrdemServicoController(osService services.OrdemServicoService) *OrdemServicoController {
    return &OrdemServicoController{
        osService: osService,
    }
}

// BuscarTodas retorna todas as ordens de serviço cadastradas
// Este endpoint pode ter parâmetros de consulta para filtrar os resultados
func (c *OrdemServicoController) BuscarTodas(ctx *gin.Context) {
    // Verificar se há filtros por status
    status := ctx.Query("status")
    if status != "" {
        c.BuscarPorStatus(ctx)
        return
    }

    // Verificar filtro por período
    dataInicio := ctx.Query("dataInicio")
    dataFim := ctx.Query("dataFim")
    if dataInicio != "" && dataFim != "" {
        c.BuscarPorPeriodo(ctx)
        return
    }

    // Buscar todas sem filtro
    ordens, err := c.osService.BuscarTodas()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ordens de serviço"})
        return
    }

    ctx.JSON(http.StatusOK, ordens)
}

// BuscarPorID retorna uma ordem de serviço específica pelo ID
// Recebe o ID como parâmetro na URL e retorna os detalhes da ordem correspondente
func (c *OrdemServicoController) BuscarPorID(ctx *gin.Context) {
    // Extrai e converte o ID da URL
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Busca a ordem pelo ID
    os, err := c.osService.BuscarPorID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Ordem de serviço não encontrada"})
        return
    }

    ctx.JSON(http.StatusOK, os)
}

// BuscarPorNumero retorna uma ordem de serviço específica pelo número
// Recebe o número como parâmetro na URL e retorna os detalhes da ordem correspondente
func (c *OrdemServicoController) BuscarPorNumero(ctx *gin.Context) {
    // Extrai o número da OS da URL
    numero := ctx.Param("numero")
    if numero == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Número da OS não informado"})
        return
    }

    // Busca a ordem pelo número
    os, err := c.osService.BuscarPorNumeroOS(numero)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Ordem de serviço não encontrada"})
        return
    }

    ctx.JSON(http.StatusOK, os)
}

// BuscarPorCliente retorna as ordens de serviço de um cliente específico
// Recebe o ID do cliente como parâmetro na URL e retorna sua lista de ordens
func (c *OrdemServicoController) BuscarPorCliente(ctx *gin.Context) {
    // Extrai e converte o ID do cliente
    clienteID, err := strconv.Atoi(ctx.Param("clienteId"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
        return
    }

    // Busca as ordens pelo ID do cliente
    ordens, err := c.osService.BuscarPorCliente(uint(clienteID))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, ordens)
}

// BuscarPorVeiculo retorna as ordens de serviço de um veículo específico
// Recebe o ID do veículo como parâmetro na URL e retorna sua lista de ordens
func (c *OrdemServicoController) BuscarPorVeiculo(ctx *gin.Context) {
    // Extrai e converte o ID do veículo
    veiculoID, err := strconv.Atoi(ctx.Param("veiculoId"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do veículo inválido"})
        return
    }

    // Busca as ordens pelo ID do veículo
    ordens, err := c.osService.BuscarPorVeiculo(uint(veiculoID))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, ordens)
}

// BuscarPorStatus retorna as ordens de serviço com um status específico
// Recebe o status como parâmetro na query e retorna a lista de ordens
func (c *OrdemServicoController) BuscarPorStatus(ctx *gin.Context) {
    // Extrai o status da query
    status := ctx.Query("status")
    if status == "" {
        status = ctx.Param("status")
    }
    
    if status == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status não informado"})
        return
    }

    // Busca as ordens pelo status
    ordens, err := c.osService.BuscarPorStatus(status)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, ordens)
}

// BuscarPorPeriodo retorna as ordens de serviço em um período específico
// Recebe datas de início e fim como parâmetros na query e retorna a lista de ordens
func (c *OrdemServicoController) BuscarPorPeriodo(ctx *gin.Context) {
    // Extrai as datas da query
    dataInicio := ctx.Query("dataInicio")
    dataFim := ctx.Query("dataFim")
    
    if dataInicio == "" || dataFim == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data de início e fim são obrigatórias"})
        return
    }

    // Converte as strings para time.Time
    inicio, err := time.Parse("2006-01-02", dataInicio)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data de início inválida. Use o formato AAAA-MM-DD"})
        return
    }

    // Ajusta a data final para o final do dia
    fim, err := time.Parse("2006-01-02", dataFim)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data de fim inválida. Use o formato AAAA-MM-DD"})
        return
    }
    fim = fim.Add(24*time.Hour - time.Second) // Final do dia

    // Busca as ordens pelo período
    ordens, err := c.osService.BuscarPorPeriodo(inicio, fim)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, ordens)
}

// Criar adiciona uma nova ordem de serviço
// Recebe os dados da ordem no corpo da requisição
func (c *OrdemServicoController) Criar(ctx *gin.Context) {
    var os models.OrdemServico

    // Faz o binding do JSON para o modelo
    if err := ctx.ShouldBindJSON(&os); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
        return
    }

    // Cria a ordem de serviço
    osCriada, err := c.osService.Criar(&os)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, osCriada)
}

// Atualizar modifica uma ordem de serviço existente
// Recebe o ID da ordem na URL e os novos dados no corpo
func (c *OrdemServicoController) Atualizar(ctx *gin.Context) {
    // Extrai e converte o ID da URL
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var os models.OrdemServico

    // Faz o binding do JSON para o modelo
    if err := ctx.ShouldBindJSON(&os); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
        return
    }

    // Define o ID usando o valor da URL
    os.ID = uint(id)

    // Atualiza a ordem de serviço
    osAtualizada, err := c.osService.Atualizar(&os)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, osAtualizada)
}

// AtualizarStatus modifica o status de uma ordem de serviço
// Recebe o ID da ordem na URL e o novo status no corpo
func (c *OrdemServicoController) AtualizarStatus(ctx *gin.Context) {
    // Extrai e converte o ID da URL
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Extrai o novo status do corpo
    var dados struct {
        Status string `json:"status" binding:"required"`
    }

    if err := ctx.ShouldBindJSON(&dados); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status não informado"})
        return
    }

    // Atualiza o status
    osAtualizada, err := c.osService.AtualizarStatus(uint(id), dados.Status)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, osAtualizada)
}

// Deletar remove uma ordem de serviço
// Recebe o ID da ordem como parâmetro na URL
func (c *OrdemServicoController) Deletar(ctx *gin.Context) {
    // Extrai e converte o ID da URL
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Remove a ordem de serviço
    err = c.osService.Deletar(uint(id))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.Status(http.StatusNoContent)
}

// AdicionarItem adiciona um item à ordem de serviço
// Recebe o ID da ordem na URL e os dados do item no corpo
func (c *OrdemServicoController) AdicionarItem(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    var item models.ItemOrdemServico

    // Faz o binding do JSON para o modelo
    if err := ctx.ShouldBindJSON(&item); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados do item inválidos: " + err.Error()})
        return
    }

    // Adiciona o item à OS
    itemAdicionado, err := c.osService.AdicionarItem(uint(id), &item)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, itemAdicionado)
}

// RemoverItem remove um item da ordem de serviço
// Recebe os IDs da ordem e do item na URL
func (c *OrdemServicoController) RemoverItem(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    osID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    // Extrai e converte o ID do item
    itemID, err := strconv.Atoi(ctx.Param("itemId"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do item inválido"})
        return
    }

    // Remove o item
    err = c.osService.RemoverItem(uint(osID), uint(itemID))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.Status(http.StatusNoContent)
}

// AtualizarItem modifica um item da ordem de serviço
// Recebe os IDs da ordem e do item na URL e os novos dados no corpo
func (c *OrdemServicoController) AtualizarItem(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    osID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    // Extrai e converte o ID do item
    itemID, err := strconv.Atoi(ctx.Param("itemId"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID do item inválido"})
        return
    }

    var item models.ItemOrdemServico

    // Faz o binding do JSON para o modelo
    if err := ctx.ShouldBindJSON(&item); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados do item inválidos: " + err.Error()})
        return
    }

    // Define o ID do item usando o valor da URL
    item.ID = uint(itemID)
    item.OrdemServicoID = uint(osID)

    // Atualiza o item
    itemAtualizado, err := c.osService.AtualizarItem(uint(osID), &item)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, itemAtualizado)
}

// BuscarItens retorna todos os itens de uma ordem de serviço
// Recebe o ID da ordem na URL e retorna a lista de itens
func (c *OrdemServicoController) BuscarItens(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    // Busca os itens da OS
    itens, err := c.osService.BuscarItens(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, itens)
}

// ConcluirOS marca uma ordem de serviço como concluída
// Recebe o ID da ordem na URL
func (c *OrdemServicoController) ConcluirOS(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    // Conclui a OS
    osAtualizada, err := c.osService.ConcluirOS(uint(id))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, osAtualizada)
}

// CancelarOS marca uma ordem de serviço como cancelada
// Recebe o ID da ordem na URL
func (c *OrdemServicoController) CancelarOS(ctx *gin.Context) {
    // Extrai e converte o ID da OS
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da ordem de serviço inválido"})
        return
    }

    // Cancela a OS
    osAtualizada, err := c.osService.CancelarOS(uint(id))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, osAtualizada)
}