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

		// 2. Tabelas com dependências
		&models.Funcionario{},
		&models.Veiculo{},
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
	// Estas são operações brutas de SQL para índices compostos ou específicos
	// que podem não ser facilmente declarados nas tags do modelo

	// Índice para pesquisa de veículos por placa (pesquisa parcial comum)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_veiculos_placa ON veiculos(placa)").Error; err != nil {
		return err
	}

	// Índice para pesquisa de clientes por nome (pesquisa parcial comum)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_clientes_nome ON clientes(nome)").Error; err != nil {
		return err
	}

	// Índice para autenticação (login frequente)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_usuarios_email ON usuarios(email)").Error; err != nil {
		return err
	}

	return nil
}
