package services

import (
	"errors"

	"OficinaMecanica/models"
	"OficinaMecanica/repositories"
)

type FuncionarioService interface {
	BuscarTodos() ([]models.Funcionario, error)
	BuscarPorID(id uint) (*models.Funcionario, error)
	Criar(funcionario *models.Funcionario) (*models.Funcionario, error)
	Atualizar(funcionario *models.Funcionario) (*models.Funcionario, error)
	Deletar(id uint) error
	BuscarPorCPF(cpf string) (*models.Funcionario, error)
	BuscarPorCargo(cargo string) ([]models.Funcionario, error)
}

type FuncionarioServiceImpl struct {
	funcionarioRepo repositories.FuncionarioRepository
}

func NewFuncionarioService(funcionarioRepo repositories.FuncionarioRepository) FuncionarioService {
	return &FuncionarioServiceImpl{
		funcionarioRepo: funcionarioRepo,
	}
}

func (s *FuncionarioServiceImpl) BuscarTodos() ([]models.Funcionario, error) {
	return s.funcionarioRepo.FindAll()
}

func (s *FuncionarioServiceImpl) BuscarPorID(id uint) (*models.Funcionario, error) {
	funcionario, err := s.funcionarioRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("funcionário não encontrado")
	}
	return funcionario, nil
}

func (s *FuncionarioServiceImpl) Criar(funcionario *models.Funcionario) (*models.Funcionario, error) {
	// Validações básicas
	if funcionario.Nome == "" {
		return nil, errors.New("nome do funcionário é obrigatório")
	}

	if funcionario.CPF != "" {
		// Verificar se já existe funcionário com o mesmo CPF
		existente, err := s.funcionarioRepo.FindByCPF(funcionario.CPF)
		if err == nil && existente != nil {
			return nil, errors.New("já existe um funcionário com este CPF")
		}

		// Validar formato do CPF
		if !validarCPF(funcionario.CPF) {
			return nil, errors.New("CPF inválido")
		}
	}

	err := s.funcionarioRepo.Create(funcionario)
	if err != nil {
		return nil, errors.New("erro ao criar funcionário")
	}

	return funcionario, nil
}

func (s *FuncionarioServiceImpl) Atualizar(funcionario *models.Funcionario) (*models.Funcionario, error) {
	// Verificar se o funcionário existe
	_, err := s.funcionarioRepo.FindByID(funcionario.ID)
	if err != nil {
		return nil, errors.New("funcionário não encontrado")
	}

	// Validações básicas
	if funcionario.Nome == "" {
		return nil, errors.New("nome do funcionário é obrigatório")
	}

	if funcionario.CPF != "" {
		// Verificar se já existe outro funcionário com o mesmo CPF
		existente, err := s.funcionarioRepo.FindByCPF(funcionario.CPF)
		if err == nil && existente != nil && existente.ID != funcionario.ID {
			return nil, errors.New("já existe outro funcionário com este CPF")
		}

		// Validar formato do CPF
		if !validarCPF(funcionario.CPF) {
			return nil, errors.New("CPF inválido")
		}
	}

	err = s.funcionarioRepo.Update(funcionario)
	if err != nil {
		return nil, errors.New("erro ao atualizar funcionário")
	}

	return funcionario, nil
}

func (s *FuncionarioServiceImpl) Deletar(id uint) error {
	// Verificar se o funcionário existe
	_, err := s.funcionarioRepo.FindByID(id)
	if err != nil {
		return errors.New("funcionário não encontrado")
	}

	err = s.funcionarioRepo.Delete(id)
	if err != nil {
		return errors.New("erro ao deletar funcionário")
	}

	return nil
}

func (s *FuncionarioServiceImpl) BuscarPorCPF(cpf string) (*models.Funcionario, error) {
	if cpf == "" {
		return nil, errors.New("CPF não pode ser vazio")
	}

	// Validar formato do CPF
	if !validarCPF(cpf) {
		return nil, errors.New("CPF inválido")
	}

	funcionario, err := s.funcionarioRepo.FindByCPF(cpf)
	if err != nil {
		return nil, errors.New("funcionário não encontrado")
	}

	return funcionario, nil
}

func (s *FuncionarioServiceImpl) BuscarPorCargo(cargo string) ([]models.Funcionario, error) {
	if cargo == "" {
		return nil, errors.New("cargo não pode ser vazio")
	}

	funcionarios, err := s.funcionarioRepo.FindByCargo(cargo)
	if err != nil {
		return nil, errors.New("erro ao buscar funcionários por cargo")
	}

	return funcionarios, nil
}
