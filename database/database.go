package database

import (
	"fmt"
	"log"

	"OficinaMecanica/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB estabelece conexão com o banco de dados usando GORM
func ConnectDB() (*gorm.DB, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	// Configurando o logger para mostrar os logs de SQL
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// Obtém a conexão SQL subjacente para configurar o pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configurações de pool de conexões
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)

	log.Println("Conexão com banco de dados estabelecida com sucesso")

	// Executa migrações, se necessário
	if err := SetupMigrations(db); err != nil {
		log.Printf("Aviso: Erro nas migrações: %v", err)
		// Continuamos mesmo se houver erro nas migrações (decisão do desenvolvedor)
	}

	return db, nil
}
