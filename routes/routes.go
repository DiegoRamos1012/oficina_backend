package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"OficinaMecanica/controllers"
	"OficinaMecanica/middlewares"
	"OficinaMecanica/repositories"
	"OficinaMecanica/services"
)

func SetupRoutes(r *gin.Engine) {
    // Obtendo conexão com banco de dados
    db := getDBConnection()
    
    // Repositórios
    clienteRepo := repositories.NewClienteRepository(db)
    veiculoRepo := repositories.NewVeiculoRepository(db)
    
    // Serviços
    clienteService := services.NewClienteService(clienteRepo)
    veiculoService := services.NewVeiculoService(veiculoRepo)
    
    // Controllers
    clienteController := controllers.NewClienteController(clienteService)
    veiculoController := controllers.NewVeiculoController(veiculoService)
    authController := controllers.NewAuthController()
    
    // Rotas públicas
    public := r.Group("/api")
    {
        public.POST("/login", authController.Login)
        public.POST("/register", authController.Register)
    }
    
    // Rotas protegidas por autenticação
    authorized := r.Group("/api")
    authorized.Use(middlewares.AuthMiddleware())
    {
        // Rotas de clientes
        clientes := authorized.Group("/clientes")
        {
            clientes.GET("/", clienteController.BuscarTodos)
            clientes.GET("/:id", clienteController.BuscarPorID)
            clientes.POST("/", clienteController.Criar)
            clientes.PUT("/:id", clienteController.Atualizar)
            clientes.DELETE("/:id", clienteController.Deletar)
        }
        
        // Rotas de veículos
        veiculos := authorized.Group("/veiculos")
        {
            veiculos.GET("/", veiculoController.BuscarTodos)
            veiculos.GET("/:id", veiculoController.BuscarPorID)
            veiculos.POST("/", veiculoController.Criar)
            veiculos.PUT("/:id", veiculoController.Atualizar)
            veiculos.DELETE("/:id", veiculoController.Deletar)
            veiculos.GET("/cliente/:clienteId", veiculoController.BuscarPorCliente)
        }
    }
}

func getDBConnection() *sql.DB {
    // Implementação da conexão com o banco de dados
    // ...
    return nil // Placeholder
}