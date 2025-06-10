package services

import (
	"errors"
	"fmt"
	"time"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type OrdemServicoService interface {
	BuscarTodas() ([]models.OrdemServico, error)
	BuscarPorID(id uint) (*models.OrdemServico, error)
	Criar(os *models.OrdemServico) (*models.OrdemServico, error)
	Atualizar(os *models.OrdemServico) (*models.OrdemServico, error)
	AtualizarStatus(id uint, novoStatus string) (*models.OrdemServico, error)
	Deletar(id uint) error
	BuscarPorCliente(clienteID uint) ([]models.OrdemServico, error)
	BuscarPorVeiculo(veiculoID uint) ([]models.OrdemServico, error)
	BuscarPorStatus(status string) ([]models.OrdemServico, error)
	BuscarPorPeriodo(inicio, fim time.Time) ([]models.OrdemServico, error)
	BuscarPorNumeroOS(numeroOS string) (*models.OrdemServico, error)
	AdicionarItem(osID uint, item *models.ItemOrdemServico) (*models.ItemOrdemServico, error)
	RemoverItem(osID uint, itemID uint) error
	AtualizarItem(osID uint, item *models.ItemOrdemServico) (*models.ItemOrdemServico, error)
	BuscarItens(osID uint) ([]models.ItemOrdemServico, error)
	ConcluirOS(id uint) (*models.OrdemServico, error)
	CancelarOS(id uint) (*models.OrdemServico, error)
}

type OrdemServicoServiceImpl struct {
	osRepo          repositories.OrdemServicoRepository
	veiculoRepo     repositories.VeiculoRepository
	clienteRepo     repositories.ClienteRepositoryGorm
	funcionarioRepo repositories.FuncionarioRepository
	estoqueRepo     repositories.EstoqueRepository
}

func NewOrdemServicoService(
	osRepo repositories.OrdemServicoRepository,
	veiculoRepo repositories.VeiculoRepository,
	clienteRepo repositories.ClienteRepositoryGorm,
	funcionarioRepo repositories.FuncionarioRepository,
	estoqueRepo repositories.EstoqueRepository,
) OrdemServicoService {
	return &OrdemServicoServiceImpl{
		osRepo:          osRepo,
		veiculoRepo:     veiculoRepo,
		clienteRepo:     clienteRepo,
		funcionarioRepo: funcionarioRepo,
		estoqueRepo:     estoqueRepo,
	}
}

func (s *OrdemServicoServiceImpl) BuscarTodas() ([]models.OrdemServico, error) {
	return s.osRepo.FindAll()
}

func (s *OrdemServicoServiceImpl) BuscarPorID(id uint) (*models.OrdemServico, error) {
	os, err := s.osRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}
	return os, nil
}

func (s *OrdemServicoServiceImpl) Criar(os *models.OrdemServico) (*models.OrdemServico, error) {
	// Validações básicas
	if os.VeiculoID == 0 {
		return nil, errors.New("veículo é obrigatório")
	}
	if os.ClienteID == 0 {
		return nil, errors.New("cliente é obrigatório")
	}
	if os.Descricao == "" {
		return nil, errors.New("descrição do serviço é obrigatória")
	}

	// Verificar se o veículo existe
	veiculo, err := s.veiculoRepo.FindByID(os.VeiculoID)
	if err != nil {
		return nil, errors.New("veículo não encontrado")
	}

	// Verificar se o cliente existe
	_, err = s.clienteRepo.FindByID(os.ClienteID)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}

	// Verificar se o funcionário existe, se fornecido
	if os.FuncionarioID > 0 {
		funcionario, err := s.funcionarioRepo.FindByID(os.FuncionarioID)
		if err != nil {
			return nil, errors.New("funcionário não encontrado")
		}
		os.Funcionario = *funcionario
	}

	// Garantir que o veículo pertence ao cliente
	if veiculo.ClienteID != os.ClienteID {
		return nil, errors.New("o veículo não pertence ao cliente informado")
	}

	// Definir valores padrão
	if os.DataEntrada.IsZero() {
		os.DataEntrada = time.Now()
	}
	if os.Status == "" {
		os.Status = "aberta"
	}

	// Persistir a ordem de serviço
	err = s.osRepo.Create(os)
	if err != nil {
		return nil, errors.New("erro ao criar ordem de serviço: " + err.Error())
	}

	return os, nil
}

func (s *OrdemServicoServiceImpl) Atualizar(os *models.OrdemServico) (*models.OrdemServico, error) {
	// Verificar se a OS existe
	osExistente, err := s.osRepo.FindByID(os.ID)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	// Não permitir alterar OS concluídas ou canceladas
	if osExistente.Status == "concluida" || osExistente.Status == "cancelada" {
		return nil, errors.New("não é possível alterar uma ordem de serviço concluída ou cancelada")
	}

	// Validações básicas
	if os.Descricao == "" {
		return nil, errors.New("descrição do serviço é obrigatória")
	}

	// Atualizar apenas campos permitidos
	osExistente.DataPrevisao = os.DataPrevisao
	osExistente.Descricao = os.Descricao
	osExistente.Diagnostico = os.Diagnostico
	osExistente.ValorServico = os.ValorServico
	osExistente.ValorDesconto = os.ValorDesconto
	osExistente.FormaPagamento = os.FormaPagamento
	osExistente.Observacoes = os.Observacoes
	osExistente.ServicosRealizados = os.ServicosRealizados

	// Se mudar o status, verificar se é uma transição válida
	if os.Status != "" && os.Status != osExistente.Status {
		if !isValidStatusTransition(osExistente.Status, os.Status) {
			return nil, fmt.Errorf("transição de status inválida: de %s para %s", osExistente.Status, os.Status)
		}
		osExistente.Status = os.Status

		// Se estiver concluindo a OS, registrar a data de conclusão
		if os.Status == "concluida" && osExistente.DataConclusao == nil {
			now := time.Now()
			osExistente.DataConclusao = &now
		}
	}

	// Persistir as alterações
	err = s.osRepo.Update(osExistente)
	if err != nil {
		return nil, errors.New("erro ao atualizar ordem de serviço: " + err.Error())
	}

	return osExistente, nil
}

func (s *OrdemServicoServiceImpl) AtualizarStatus(id uint, novoStatus string) (*models.OrdemServico, error) {
	// Buscar a OS
	os, err := s.osRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	// Validar o status
	if !isValidStatus(novoStatus) {
		return nil, errors.New("status inválido")
	}

	// Verificar transição válida
	if !isValidStatusTransition(os.Status, novoStatus) {
		return nil, fmt.Errorf("transição de status inválida: de %s para %s", os.Status, novoStatus)
	}

	// Atualizar status
	os.Status = novoStatus

	// Se concluindo a OS, registrar data de conclusão
	if novoStatus == "concluida" && os.DataConclusao == nil {
		now := time.Now()
		os.DataConclusao = &now
	}

	// Persistir as alterações
	err = s.osRepo.Update(os)
	if err != nil {
		return nil, errors.New("erro ao atualizar status: " + err.Error())
	}

	return os, nil
}

func (s *OrdemServicoServiceImpl) Deletar(id uint) error {
	// Verificar se a OS existe
	os, err := s.osRepo.FindByID(id)
	if err != nil {
		return errors.New("ordem de serviço não encontrada")
	}

	// Não permitir excluir OS concluídas ou em andamento
	if os.Status == "concluida" || os.Status == "emandamento" {
		return errors.New("não é possível excluir uma ordem de serviço concluída ou em andamento")
	}

	// Excluir a OS
	return s.osRepo.Delete(id)
}

func (s *OrdemServicoServiceImpl) BuscarPorCliente(clienteID uint) ([]models.OrdemServico, error) {
	// Verificar se o cliente existe
	_, err := s.clienteRepo.FindByID(clienteID)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}

	return s.osRepo.FindByClienteID(clienteID)
}

func (s *OrdemServicoServiceImpl) BuscarPorVeiculo(veiculoID uint) ([]models.OrdemServico, error) {
	// Verificar se o veículo existe
	_, err := s.veiculoRepo.FindByID(veiculoID)
	if err != nil {
		return nil, errors.New("veículo não encontrado")
	}

	return s.osRepo.FindByVeiculoID(veiculoID)
}

func (s *OrdemServicoServiceImpl) BuscarPorStatus(status string) ([]models.OrdemServico, error) {
	// Validar o status
	if !isValidStatus(status) {
		return nil, errors.New("status inválido")
	}

	return s.osRepo.FindByStatus(status)
}

func (s *OrdemServicoServiceImpl) BuscarPorPeriodo(inicio, fim time.Time) ([]models.OrdemServico, error) {
	// Validar o período
	if inicio.After(fim) {
		return nil, errors.New("data inicial deve ser anterior à data final")
	}

	return s.osRepo.FindByPeriodo(inicio, fim)
}

func (s *OrdemServicoServiceImpl) BuscarPorNumeroOS(numeroOS string) (*models.OrdemServico, error) {
	if numeroOS == "" {
		return nil, errors.New("número da OS é obrigatório")
	}

	return s.osRepo.FindByNumeroOS(numeroOS)
}

func (s *OrdemServicoServiceImpl) AdicionarItem(osID uint, item *models.ItemOrdemServico) (*models.ItemOrdemServico, error) {
	// Verificar se a OS existe
	os, err := s.osRepo.FindByID(osID)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	// Não permitir adicionar itens a OS concluídas ou canceladas
	if os.Status == "concluida" || os.Status == "cancelada" {
		return nil, errors.New("não é possível adicionar itens a uma OS concluída ou cancelada")
	}

	// Verificar se o item existe no estoque
	estoqueItem, err := s.estoqueRepo.FindByID(item.EstoqueID)
	if err != nil {
		return nil, errors.New("item de estoque não encontrado")
	}

	// Verificar se há quantidade suficiente
	if estoqueItem.Quantidade < item.Quantidade {
		return nil, errors.New("quantidade insuficiente em estoque")
	}

	// Definir valores do item
	item.OrdemServicoID = osID
	if item.ValorUnitario <= 0 {
		item.ValorUnitario = estoqueItem.PrecoVenda
	}
	item.ValorTotal = float64(item.Quantidade) * item.ValorUnitario

	// Adicionar o item
	err = s.osRepo.AddItem(item)
	if err != nil {
		return nil, errors.New("erro ao adicionar item: " + err.Error())
	}

	// Atualizar estoque
	estoqueItem.Quantidade -= item.Quantidade
	err = s.estoqueRepo.Update(estoqueItem)
	if err != nil {
		return nil, errors.New("erro ao atualizar estoque: " + err.Error())
	}

	// Atualizar valor de peças da OS
	os.ValorPecas += item.ValorTotal
	err = s.osRepo.Update(os)
	if err != nil {
		return nil, errors.New("erro ao atualizar valor total da OS: " + err.Error())
	}

	return item, nil
}

func (s *OrdemServicoServiceImpl) RemoverItem(osID uint, itemID uint) error {
	// Verificar se a OS existe
	os, err := s.osRepo.FindByID(osID)
	if err != nil {
		return errors.New("ordem de serviço não encontrada")
	}

	// Não permitir remover itens de OS concluídas ou canceladas
	if os.Status == "concluida" || os.Status == "cancelada" {
		return errors.New("não é possível remover itens de uma OS concluída ou cancelada")
	}

	// Buscar itens da OS
	itens, err := s.osRepo.FindItens(osID)
	if err != nil {
		return errors.New("erro ao buscar itens da OS")
	}

	// Verificar se o item pertence à OS
	var itemParaRemover *models.ItemOrdemServico
	for _, i := range itens {
		if i.ID == itemID {
			itemParaRemover = &i
			break
		}
	}

	if itemParaRemover == nil {
		return errors.New("item não encontrado na ordem de serviço")
	}

	// Devolver ao estoque
	estoqueItem, err := s.estoqueRepo.FindByID(itemParaRemover.EstoqueID)
	if err != nil {
		return errors.New("item de estoque não encontrado")
	}

	// Atualizar estoque
	estoqueItem.Quantidade += itemParaRemover.Quantidade
	err = s.estoqueRepo.Update(estoqueItem)
	if err != nil {
		return errors.New("erro ao atualizar estoque: " + err.Error())
	}

	// Atualizar valor da OS
	os.ValorPecas -= itemParaRemover.ValorTotal
	if os.ValorPecas < 0 {
		os.ValorPecas = 0
	}
	err = s.osRepo.Update(os)
	if err != nil {
		return errors.New("erro ao atualizar valor da OS: " + err.Error())
	}

	// Remover o item
	return s.osRepo.RemoveItem(itemID)
}

func (s *OrdemServicoServiceImpl) AtualizarItem(osID uint, item *models.ItemOrdemServico) (*models.ItemOrdemServico, error) {
	// Verificar se a OS existe
	os, err := s.osRepo.FindByID(osID)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	// Não permitir atualizar itens de OS concluídas ou canceladas
	if os.Status == "concluida" || os.Status == "cancelada" {
		return nil, errors.New("não é possível atualizar itens de uma OS concluída ou cancelada")
	}

	// Buscar o item atual
	itens, err := s.osRepo.FindItens(osID)
	if err != nil {
		return nil, errors.New("erro ao buscar itens da OS")
	}

	var itemAtual *models.ItemOrdemServico
	for _, i := range itens {
		if i.ID == item.ID {
			itemAtual = &i
			break
		}
	}

	if itemAtual == nil {
		return nil, errors.New("item não encontrado na ordem de serviço")
	}

	// Verificar o estoque
	estoqueItem, err := s.estoqueRepo.FindByID(itemAtual.EstoqueID)
	if err != nil {
		return nil, errors.New("item de estoque não encontrado")
	}

	// Calcular diferença de quantidade
	diferencaQuantidade := item.Quantidade - itemAtual.Quantidade

	// Verificar se há estoque suficiente
	if diferencaQuantidade > 0 && estoqueItem.Quantidade < diferencaQuantidade {
		return nil, errors.New("quantidade insuficiente em estoque")
	}

	// Atualizar estoque
	estoqueItem.Quantidade -= diferencaQuantidade
	err = s.estoqueRepo.Update(estoqueItem)
	if err != nil {
		return nil, errors.New("erro ao atualizar estoque: " + err.Error())
	}

	// Calcular novo valor total do item
	valorTotalAnterior := itemAtual.ValorTotal
	item.ValorTotal = float64(item.Quantidade) * item.ValorUnitario

	// Atualizar o item
	err = s.osRepo.UpdateItem(item)
	if err != nil {
		return nil, errors.New("erro ao atualizar item: " + err.Error())
	}

	// Atualizar valor da OS
	os.ValorPecas = os.ValorPecas - valorTotalAnterior + item.ValorTotal
	err = s.osRepo.Update(os)
	if err != nil {
		return nil, errors.New("erro ao atualizar valor da OS: " + err.Error())
	}

	return item, nil
}

func (s *OrdemServicoServiceImpl) BuscarItens(osID uint) ([]models.ItemOrdemServico, error) {
	// Verificar se a OS existe
	_, err := s.osRepo.FindByID(osID)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	return s.osRepo.FindItens(osID)
}

func (s *OrdemServicoServiceImpl) ConcluirOS(id uint) (*models.OrdemServico, error) {
	return s.AtualizarStatus(id, "concluida")
}

func (s *OrdemServicoServiceImpl) CancelarOS(id uint) (*models.OrdemServico, error) {
	// Buscar a OS
	os, err := s.osRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	// Não permitir cancelar OS concluídas
	if os.Status == "concluida" {
		return nil, errors.New("não é possível cancelar uma OS concluída")
	}

	// Devolver itens ao estoque
	itens, err := s.osRepo.FindItens(id)
	if err != nil {
		return nil, errors.New("erro ao buscar itens da OS")
	}

	// Para cada item, devolver ao estoque
	for _, item := range itens {
		estoqueItem, err := s.estoqueRepo.FindByID(item.EstoqueID)
		if err != nil {
			return nil, errors.New("item de estoque não encontrado: " + err.Error())
		}

		estoqueItem.Quantidade += item.Quantidade
		err = s.estoqueRepo.Update(estoqueItem)
		if err != nil {
			return nil, errors.New("erro ao atualizar estoque: " + err.Error())
		}
	}

	// Atualizar status da OS
	return s.AtualizarStatus(id, "cancelada")
}

// Funções auxiliares
func isValidStatus(status string) bool {
	validStatus := []string{"aberta", "emandamento", "concluida", "cancelada"}
	for _, s := range validStatus {
		if s == status {
			return true
		}
	}
	return false
}

func isValidStatusTransition(atual, novo string) bool {
	// Define transições válidas
	transicoes := map[string][]string{
		"aberta":      {"emandamento", "cancelada"},
		"emandamento": {"concluida", "cancelada"},
		"concluida":   {}, // Não permite transição após concluída
		"cancelada":   {}, // Não permite transição após cancelada
	}

	// Verifica se a transição é válida
	permitidas, ok := transicoes[atual]
	if !ok {
		return false
	}

	// Se o status não mudar, é sempre permitido
	if atual == novo {
		return true
	}

	// Verifica se o novo status está entre os permitidos
	for _, s := range permitidas {
		if s == novo {
			return true
		}
	}

	return false
}
