package services

import (
	"errors"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type VeiculoService interface {
	BuscarTodos() ([]models.Veiculo, error)
	BuscarPorID(id uint) (*models.Veiculo, error)
	Criar(veiculo *models.Veiculo) (*models.Veiculo, error)
	Atualizar(veiculo *models.Veiculo) (*models.Veiculo, error)
	Deletar(id uint) error
	BuscarPorPlaca(placa string) (*models.Veiculo, error)
	BuscarPorClienteID(clienteID uint) ([]models.Veiculo, error)
}

type VeiculoServiceImpl struct {
	veiculoRepo repositories.VeiculoRepository
	clienteRepo repositories.ClienteRepositoryGorm
}

func NewVeiculoService(veiculoRepo repositories.VeiculoRepository) VeiculoService {
	return &VeiculoServiceImpl{
		veiculoRepo: veiculoRepo,
	}
}

func (s *VeiculoServiceImpl) BuscarTodos() ([]models.Veiculo, error) {
	return s.veiculoRepo.FindAll()
}

func (s *VeiculoServiceImpl) BuscarPorID(id uint) (*models.Veiculo, error) {
	veiculo, err := s.veiculoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("veículo não encontrado")
	}
	return veiculo, nil
}

func (s *VeiculoServiceImpl) Criar(veiculo *models.Veiculo) (*models.Veiculo, error) {
	// Validações antes de criar o veículo
	if veiculo.Placa == "" {
		return nil, errors.New("placa do veículo é obrigatória")
	}

	// Verificar se já existe veículo com a mesma placa
	existente, err := s.veiculoRepo.FindByPlaca(veiculo.Placa)
	if err == nil && existente != nil {
		return nil, errors.New("já existe um veículo com esta placa")
	}

	// Verificar se o cliente existe
	if veiculo.ClienteID > 0 {
		_, err := s.clienteRepo.FindByID(veiculo.ClienteID)
		if err != nil {
			return nil, errors.New("cliente não encontrado")
		}
	}

	err = s.veiculoRepo.Create(veiculo)
	if err != nil {
		return nil, errors.New("erro ao criar veículo")
	}

	return veiculo, nil
}

func (s *VeiculoServiceImpl) Atualizar(veiculo *models.Veiculo) (*models.Veiculo, error) {
	// Verificar se o veículo existe
	_, err := s.veiculoRepo.FindByID(veiculo.ID)
	if err != nil {
		return nil, errors.New("veículo não encontrado")
	}

	// Validações
	if veiculo.Placa == "" {
		return nil, errors.New("placa do veículo é obrigatória")
	}

	// Verificar se existe outro veículo com a mesma placa
	existente, err := s.veiculoRepo.FindByPlaca(veiculo.Placa)
	if err == nil && existente != nil && existente.ID != veiculo.ID {
		return nil, errors.New("já existe outro veículo com esta placa")
	}

	// Verificar se o cliente existe
	if veiculo.ClienteID > 0 {
		_, err := s.clienteRepo.FindByID(veiculo.ClienteID)
		if err != nil {
			return nil, errors.New("cliente não encontrado")
		}
	}

	err = s.veiculoRepo.Update(veiculo)
	if err != nil {
		return nil, errors.New("erro ao atualizar veículo")
	}

	return veiculo, nil
}

func (s *VeiculoServiceImpl) Deletar(id uint) error {
	// Verificar se o veículo existe
	_, err := s.veiculoRepo.FindByID(id)
	if err != nil {
		return errors.New("veículo não encontrado")
	}

	err = s.veiculoRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar veículo")
	}

	return nil
}

func (s *VeiculoServiceImpl) BuscarPorPlaca(placa string) (*models.Veiculo, error) {
	if placa == "" {
		return nil, errors.New("placa não pode ser vazia")
	}

	veiculo, err := s.veiculoRepo.FindByPlaca(placa)
	if err != nil {
		return nil, errors.New("veículo não encontrado")
	}

	return veiculo, nil
}

func (s *VeiculoServiceImpl) BuscarPorClienteID(clienteID uint) ([]models.Veiculo, error) {
	if clienteID == 0 {
		return nil, errors.New("ID do cliente é obrigatório")
	}

	// Verificar se o cliente existe
	_, err := s.clienteRepo.FindByID(clienteID)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}

	veiculos, err := s.veiculoRepo.FindByClienteID(clienteID)
	if err != nil {
		return nil, errors.New("erro ao buscar veículos do cliente")
	}

	return veiculos, nil
}
