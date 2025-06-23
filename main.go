package main

import (
	"OficinaMecanica/configs"
	"OficinaMecanica/database"
	"OficinaMecanica/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar o router
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Carregar configurações
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Executar migrations apenas se as tabelas não existirem
	if err := database.SetupMigrations(db); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	// Configurar rotas
	routes.SetupRoutes(r)

	// Iniciar o servidor
	log.Printf("Servidor iniciado na porta %s", config.ServerPort)
	r.Run(":" + config.ServerPort)
}
