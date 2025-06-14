package database

import (
	"log"
	"time"

	"OficinaMecanica/models"

	"gorm.io/gorm"
)

// SetupMigrations configura e executa as migrações do banco de dados
func SetupMigrations(db *gorm.DB) error {
	// Configurações para otimização das migrações
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	// Log de início
	start := time.Now()
	log.Println("Iniciando migrações do banco de dados...")

	// Executando migrações em ordem apropriada para respeitar dependências
	err := db.AutoMigrate(
		// 1. Tabelas independentes primeiro
		&models.Usuario{},
		&models.Cliente{},
		&models.Estoque{},

		// 2. Tabelas com dependências
		&models.Funcionario{},
		&models.Veiculo{},

		// 3. Tabelas que dependem das anteriores
		&models.OrdemServico{},
		&models.ItemOrdemServico{},
	)

	if err != nil {
		log.Printf("Erro nas migrações: %v", err)
		return err
	}

	// Adicionar índices para otimização de consultas frequentes
	err = addOptimizationIndexes(db)
	if err != nil {
		log.Printf("Erro ao adicionar índices: %v", err)
		return err
	}

	// Log de conclusão
	log.Printf("Migrações concluídas em %v", time.Since(start))
	return nil
}

// addOptimizationIndexes adiciona índices para otimização de consultas comuns
func addOptimizationIndexes(db *gorm.DB) error {
	// Verificar se índice já existe antes de criar
	createIndexIfNotExists := func(indexName, tableName, columns string) error {
		var count int64
		// Verificar se o índice já existe
		db.Raw("SELECT COUNT(1) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			tableName, indexName).Count(&count)

		if count == 0 {
			// Índice não existe, então cria
			sql := "CREATE INDEX " + indexName + " ON " + tableName + "(" + columns + ")"
			if err := db.Exec(sql).Error; err != nil {
				return err
			}
		}
		return nil
	}

	// Criar os índices se não existirem
	if err := createIndexIfNotExists("idx_veiculos_placa", "veiculos", "placa"); err != nil {
		return err
	}

	if err := createIndexIfNotExists("idx_clientes_nome", "clientes", "nome"); err != nil {
		return err
	}

	if err := createIndexIfNotExists("idx_usuarios_email", "usuarios", "email"); err != nil {
		return err
	}

	return nil
}
