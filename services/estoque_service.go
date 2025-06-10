package services

import (
	"errors"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type EstoqueService interface {
	BuscarTodos() ([]models.Estoque, error)
	BuscarPorID(id uint) (*models.Estoque, error)
	Criar(estoque *models.Estoque) (*models.Estoque, error)
	Atualizar(estoque *models.Estoque) (*models.Estoque, error)
	Deletar(id uint) error
	BuscarPorCategoria(categoria string) ([]models.Estoque, error)
	BuscarBaixoEstoque() ([]models.Estoque, error)
}

type EstoqueServiceImpl struct {
	estoqueRepo repositories.EstoqueRepository
}

func NewEstoqueService(estoqueRepo repositories.EstoqueRepository) EstoqueService {
	return &EstoqueServiceImpl{
		estoqueRepo: estoqueRepo,
	}
}

func (s *EstoqueServiceImpl) BuscarTodos() ([]models.Estoque, error) {
	return s.estoqueRepo.FindAll()
}

func (s *EstoqueServiceImpl) BuscarPorID(id uint) (*models.Estoque, error) {
	item, err := s.estoqueRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("item não encontrado")
	}
	return item, nil
}

func (s *EstoqueServiceImpl) Criar(estoque *models.Estoque) (*models.Estoque, error) {
	// Validações antes de criar o item
	if estoque.Nome == "" {
		return nil, errors.New("nome do item é obrigatório")
	}

	if estoque.PrecoVenda < estoque.PrecoUnitario {
		return nil, errors.New("preço de venda não pode ser menor que o preço de custo")
	}

	err := s.estoqueRepo.Create(estoque)
	if err != nil {
		return nil, errors.New("erro ao criar item no estoque")
	}

	return estoque, nil
}

func (s *EstoqueServiceImpl) Atualizar(estoque *models.Estoque) (*models.Estoque, error) {
	// Verificar se o item existe
	_, err := s.estoqueRepo.FindByID(estoque.ID)
	if err != nil {
		return nil, errors.New("item não encontrado")
	}

	// Aplicar validações
	if estoque.Nome == "" {
		return nil, errors.New("nome do item é obrigatório")
	}

	if estoque.PrecoVenda < estoque.PrecoUnitario {
		return nil, errors.New("preço de venda não pode ser menor que o preço de custo")
	}

	err = s.estoqueRepo.Update(estoque)
	if err != nil {
		return nil, errors.New("erro ao atualizar item")
	}

	return estoque, nil
}

func (s *EstoqueServiceImpl) Deletar(id uint) error {
	// Verificar se o item existe
	_, err := s.estoqueRepo.FindByID(id)
	if err != nil {
		return errors.New("item não encontrado")
	}

	err = s.estoqueRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar item")
	}

	return nil
}

func (s *EstoqueServiceImpl) BuscarPorCategoria(categoria string) ([]models.Estoque, error) {
	if categoria == "" {
		return nil, errors.New("categoria não pode ser vazia")
	}

	itens, err := s.estoqueRepo.FindByCategoria(categoria)
	if err != nil {
		return nil, errors.New("erro ao buscar itens por categoria")
	}

	return itens, nil
}

func (s *EstoqueServiceImpl) BuscarBaixoEstoque() ([]models.Estoque, error) {
	itens, err := s.estoqueRepo.FindBaixoEstoque()
	if err != nil {
		return nil, errors.New("erro ao buscar itens com estoque baixo")
	}

	return itens, nil
}
