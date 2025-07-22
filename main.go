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
	// 1. Carregar configurações e definir ambiente
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// 2. Configurar logs baseados no ambiente
	setupLogs(config.Environment)

	// 3. Definir modo do Gin antes de criar o router
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Inicializar o banco de dados
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// 5. Executar migrations
	if err := database.SetupMigrations(db); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	// 6. Inicializar o router
	r := gin.Default()

	// 7. Configurar middlewares globais (CORS, etc.)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 8. Log ambiente e porta
	log.Printf("Ambiente: %s, Servidor na porta: %s", config.Environment, config.ServerPort)

	// 9. Configurar rotas
	routes.SetupRoutes(r)

	// 10. Servir arquivos estáticos para uploads de avatar
	r.Static("/uploads", "./uploads")

	// 11. Iniciar o servidor
	log.Printf("Servidor iniciado na porta %s", config.ServerPort)
	r.Run(":" + config.ServerPort)
}

// setupLogs configura o comportamento dos logs dependendo do ambiente
func setupLogs(env string) {
	if env == "production" {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Println("Iniciando servidor em ambiente de produção")
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("======================================================")
		log.Println("🔧 INICIANDO SISTEMA OFICINA MECÂNICA - MODO DESENVOLVIMENTO")
		log.Println("- Projeto feito por Diego Ramos dos Santos. Github: Diego1012 -")
		log.Println("======================================================")
		log.Println("📌 RESUMO DO PROJETO:")
		log.Println(" • Backend em Go com Gin Framework")
		log.Println(" • Banco de dados MySQL com GORM ORM")
		log.Println(" • Autenticação via JWT")
		log.Println(" • API RESTful para gestão de oficina mecânica")
		log.Println("======================================================")
		log.Println("🔍 ESTRUTURA PRINCIPAL:")
		log.Println(" • /models     - Entidades e estruturas de dados")
		log.Println(" • /controllers - Manipuladores de requisições HTTP")
		log.Println(" • /services   - Lógica de negócio")
		log.Println(" • /repositories - Acesso a dados")
		log.Println(" • /middlewares - Interceptadores de requisições")
		log.Println(" • /database   - Configuração e migrações do BD")
		log.Println("======================================================")
		log.Println("🚀 ENDPOINTS PRINCIPAIS:")
		log.Println(" • POST /api/login          - Autenticação")
		log.Println(" • GET  /api/clientes       - Lista clientes")
		log.Println(" • GET  /api/veiculos       - Lista veículos")
		log.Println(" • GET  /api/ordens-servico - Lista ordens de serviço")
		log.Println("======================================================")
		log.Println("⚙️  FLUXO DE INICIALIZAÇÃO:")
		log.Println(" 1. Carregamento de configurações (.env)")
		log.Println(" 2. Conexão com banco de dados")
		log.Println(" 3. Verificação e execução de migrações")
		log.Println(" 4. Configuração de rotas e middlewares")
		log.Println(" 5. Inicialização do servidor HTTP")
		log.Println("======================================================")
	}
}
