package main

import (
	"fmt"
	"log"

	"OficinaMecanica/database"
	"github.com/spf13/viper"
)

func main() {
	// Configurar Viper para usar test.env
	viper.SetConfigFile("test.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Erro ao ler arquivo test.env: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Executar migrações
	if err := database.SetupMigrations(db); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	fmt.Println("Migrações executadas com sucesso no banco de teste!")
}
