package services

import (
    "errors"
    
    "OficinaMecanica/models"
    "OficinaMecanica/repositories"

)

type ClienteService interface {
    BuscarTodos() ([]models.Cliente, error)
    BuscarPorID(id int) (models.Cliente, error)
    Criar(cliente models.Cliente) (models.Cliente, error)
    Atualizar(cliente models.Cliente) (models.Cliente, error)
    Deletar(id int) error
    BuscarClienteComVeiculos(id int) (models.ClienteVeiculosDTO, error)
}

type ClienteServiceImpl struct {
    clienteRepo repositories.ClienteRepository
    veiculoRepo repositories.VeiculoRepository
}

func NewClienteService(clienteRepo repositories.ClienteRepository) ClienteService {
    return &ClienteServiceImpl{
        clienteRepo: clienteRepo,
    }
}

func (s *ClienteServiceImpl) BuscarTodos() ([]models.Cliente, error) {
    return s.clienteRepo.BuscarTodos()
}

func (s *ClienteServiceImpl) BuscarPorID(id int) (models.Cliente, error) {
    cliente, err := s.clienteRepo.BuscarPorID(id)
    if err != nil {
        return models.Cliente{}, errors.New("cliente não encontrado")
    }
    return cliente, nil
}

func (s *ClienteServiceImpl) Criar(cliente models.Cliente) (models.Cliente, error) {
    // Validar CPF
    if !validarCPF(cliente.CPF) {
        return models.Cliente{}, errors.New("CPF inválido")
    }
    
    return s.clienteRepo.Criar(cliente)
}

func (s *ClienteServiceImpl) BuscarClienteComVeiculos(id int) (models.ClienteVeiculosDTO, error) {
    cliente, err := s.clienteRepo.BuscarPorID(id)
    if err != nil {
        return models.ClienteVeiculosDTO{}, errors.New("cliente não encontrado")
    }
    
    veiculos, err := s.veiculoRepo.BuscarPorCliente(id)
    if err != nil {
        return models.ClienteVeiculosDTO{}, errors.New("erro ao buscar veículos do cliente")
    }
    
    return models.ClienteVeiculosDTO{
        Cliente:  cliente,
        Veiculos: veiculos,
    }, nil
}

// Função auxiliar para validação de CPF
func validarCPF(cpf string) bool {
    // Implementação da validação de CPF
    return true // Placeholder
}