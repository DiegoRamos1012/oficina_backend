package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"OficinaMecanica/configs"
)

func ConnectDB() (*sql.DB, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conexão com banco de dados estabelecida com sucesso")

	// Configurações de pool de conexões
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	return db, nil
}
