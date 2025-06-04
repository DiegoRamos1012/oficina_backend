package repositories

import (
    "database/sql"
    "time"
    
    "OficinaMecanica/models"
)

type ClienteRepository interface {
    BuscarTodos() ([]models.Cliente, error)
    BuscarPorID(id int) (models.Cliente, error)
    Criar(cliente models.Cliente) (models.Cliente, error)
    Atualizar(cliente models.Cliente) (models.Cliente, error)
    Deletar(id int) error
}

type ClienteRepositoryImpl struct {
    db *sql.DB
}

func NewClienteRepository(db *sql.DB) ClienteRepository {
    return &ClienteRepositoryImpl{db: db}
}

func (r *ClienteRepositoryImpl) BuscarTodos() ([]models.Cliente, error) {
    query := `SELECT id, nome, cpf, email, telefone, data_nascimento, 
              endereco, criado_em, atualizado_em FROM clientes`
              
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var clientes []models.Cliente
    
    for rows.Next() {
        var cliente models.Cliente
        err := rows.Scan(
            &cliente.ID, 
            &cliente.Nome, 
            &cliente.CPF, 
            &cliente.Email, 
            &cliente.Telefone, 
            &cliente.DataNascimento, 
            &cliente.Endereco, 
            &cliente.CriadoEm, 
            &cliente.AtualizadoEm,
        )
        if err != nil {
            return nil, err
        }
        clientes = append(clientes, cliente)
    }
    
    return clientes, nil
}

func (r *ClienteRepositoryImpl) Criar(cliente models.Cliente) (models.Cliente, error) {
    query := `INSERT INTO clientes (nome, cpf, email, telefone, data_nascimento, endereco, criado_em, atualizado_em) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

    agora := time.Now()
    cliente.CriadoEm = agora
    cliente.AtualizadoEm = agora

    result, err := r.db.Exec(query, 
        cliente.Nome, 
        cliente.CPF, 
        cliente.Email, 
        cliente.Telefone, 
        cliente.DataNascimento,
        cliente.Endereco, 
        cliente.CriadoEm, 
        cliente.AtualizadoEm)
        
    if err != nil {
        return models.Cliente{}, err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return models.Cliente{}, err
    }
    
    cliente.ID = int(id)
    return cliente, nil
}