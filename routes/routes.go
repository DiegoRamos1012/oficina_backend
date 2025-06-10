package routes

import (
	"OficinaMecanica/controllers"
	"OficinaMecanica/database"
	"OficinaMecanica/middlewares"
	"OficinaMecanica/repositories"
	"OficinaMecanica/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine) {
	// Obtendo conexão com banco de dados
	db := getDBConnection()

	// Repositórios
	clienteRepo := repositories.NewClienteRepositoryGorm(db)
	veiculoRepo := repositories.NewVeiculoRepository(db)
	estoqueRepo := repositories.NewEstoqueRepository(db)
	funcionarioRepo := repositories.NewFuncionarioRepository(db)

	// Serviços
	clienteService := services.NewClienteService(clienteRepo)
	veiculoService := services.NewVeiculoService(veiculoRepo)
	estoqueService := services.NewEstoqueService(estoqueRepo)
	funcionarioService := services.NewFuncionarioService(funcionarioRepo)

	// Controllers
	clienteController := controllers.NewClienteController(clienteService)
	veiculoController := controllers.NewVeiculoController(veiculoService)
	estoqueController := controllers.NewEstoqueController(estoqueService)
	funcionarioController := controllers.NewFuncionarioController(funcionarioService)
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

		// Rotas de estoque
		estoque := authorized.Group("/estoque")
		{
			estoque.GET("/", estoqueController.BuscarTodos)
			estoque.GET("/:id", estoqueController.BuscarPorID)
			estoque.POST("/", estoqueController.Criar)
			estoque.PUT("/:id", estoqueController.Atualizar)
			estoque.DELETE("/:id", estoqueController.Deletar)
			estoque.GET("/categoria/:categoria", estoqueController.BuscarPorCategoria)
			estoque.GET("/baixo-estoque", estoqueController.BuscarBaixoEstoque)
		}

		// Rotas de funcionários
		funcionarios := authorized.Group("/funcionarios")
		{
			funcionarios.GET("/", funcionarioController.BuscarTodos)
			funcionarios.GET("/:id", funcionarioController.BuscarPorID)
			funcionarios.POST("/", funcionarioController.Criar)
			funcionarios.PUT("/:id", funcionarioController.Atualizar)
			funcionarios.DELETE("/:id", funcionarioController.Deletar)
			funcionarios.GET("/cpf/:cpf", funcionarioController.BuscarPorCPF)
			funcionarios.GET("/cargo/:cargo", funcionarioController.BuscarPorCargo)
		}
	}
}

// getDBConnection retorna uma conexão com o banco de dados usando GORM
func getDBConnection() *gorm.DB {
	db, err := database.ConnectDB()
	if err != nil {
		// Em produção, você deve lidar com este erro de forma mais robusta
		// Considere usar um logger estruturado e talvez reiniciar o aplicativo
		panic("Falha ao conectar ao banco de dados: " + err.Error())
	}
	return db
}
