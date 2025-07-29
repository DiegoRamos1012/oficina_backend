package services

import (
	"errors"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type ClienteService interface {
	BuscarTodos() ([]models.Cliente, error)
	BuscarPorID(id uint) (*models.Cliente, error)
	Criar(cliente *models.Cliente) (*models.Cliente, error)
	Atualizar(cliente *models.Cliente) (*models.Cliente, error)
	Deletar(id uint) error
	BuscarClienteComVeiculos(id uint) (*models.ClienteVeiculosDTO, error)
}

type ClienteServiceImpl struct {
	clienteRepo repositories.ClienteRepositoryGorm
}

func NewClienteService(clienteRepo repositories.ClienteRepositoryGorm) ClienteService {
	return &ClienteServiceImpl{
		clienteRepo: clienteRepo,
	}
}

func (s *ClienteServiceImpl) BuscarTodos() ([]models.Cliente, error) {
	return s.clienteRepo.FindAll()
}

func (s *ClienteServiceImpl) BuscarPorID(id uint) (*models.Cliente, error) {
	cliente, err := s.clienteRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}
	return cliente, nil
}

func (s *ClienteServiceImpl) Criar(cliente *models.Cliente) (*models.Cliente, error) {
	err := s.clienteRepo.Create(cliente)
	if err != nil {
		return nil, errors.New("erro ao criar cliente")
	}
	return cliente, nil
}

func (s *ClienteServiceImpl) Atualizar(cliente *models.Cliente) (*models.Cliente, error) {
	err := s.clienteRepo.Update(cliente)
	if err != nil {
		return nil, errors.New("erro ao atualizar cliente")
	}
	return cliente, nil
}

func (s *ClienteServiceImpl) Deletar(id uint) error {
	err := s.clienteRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar cliente")
	}
	return nil
}

func (s *ClienteServiceImpl) BuscarClienteComVeiculos(id uint) (*models.ClienteVeiculosDTO, error) {
	cliente, err := s.clienteRepo.FindWithVeiculos(id)
	if err != nil {
		return nil, errors.New("cliente não encontrado ou erro ao buscar veículos")
	}

	// No GORM, cliente já vem com veículos preenchidos pela relação
	// mas mantemos o DTO para consistência com a API anterior
	dto := &models.ClienteVeiculosDTO{
		Cliente:  *cliente,
		Veiculos: cliente.Veiculos,
	}

	return dto, nil
}

