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
	usuarioRepo := repositories.NewUsuarioRepository(db)
	clienteRepo := repositories.NewClienteRepositoryGorm(db)
	veiculoRepo := repositories.NewVeiculoRepository(db)
	estoqueRepo := repositories.NewEstoqueRepository(db)
	funcionarioRepo := repositories.NewFuncionarioRepository(db)
	ordemServicoRepo := repositories.NewOrdemServicoRepository(db)

	// Serviços
	usuarioService := services.NewUsuarioService(usuarioRepo)
	clienteService := services.NewClienteService(clienteRepo)
	veiculoService := services.NewVeiculoService(veiculoRepo)
	estoqueService := services.NewEstoqueService(estoqueRepo)
	funcionarioService := services.NewFuncionarioService(funcionarioRepo)
	ordemServicoService := services.NewOrdemServicoService(ordemServicoRepo, veiculoRepo, clienteRepo, funcionarioRepo, estoqueRepo)

	// Controllers
	authController := controllers.NewAuthController(usuarioService)
	usuarioController := controllers.NewUsuarioController(usuarioService)
	clienteController := controllers.NewClienteController(clienteService)
	veiculoController := controllers.NewVeiculoController(veiculoService)
	estoqueController := controllers.NewEstoqueController(estoqueService)
	funcionarioController := controllers.NewFuncionarioController(funcionarioService)
	ordemServicoController := controllers.NewOrdemServicoController(ordemServicoService)

	// Rotas públicas
	public := r.Group("/api")
	{
		// Rotas de autenticação
		public.POST("/login", authController.Login)
		public.POST("/register", authController.Register)
		public.GET("/validate-token", middlewares.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(200, gin.H{"valid": true})
		})
	}

	// Rotas protegidas por autenticação
	authorized := r.Group("/api")
	authorized.Use(middlewares.AuthMiddleware())
	{
		// Rotas de usuários
		usuarios := authorized.Group("/usuarios")
		{
			usuarios.GET("/", usuarioController.BuscarTodos)
			usuarios.GET("/:id", usuarioController.BuscarPorID)
			usuarios.POST("/", usuarioController.Criar)
			usuarios.PUT("/:id", usuarioController.Atualizar)
			usuarios.DELETE("/:id", usuarioController.Deletar)
			usuarios.POST("/:id/avatar", usuarioController.UploadAvatar) // Rota para upload de avatar
		}

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

		// Rotas de ordens de serviço
		os := authorized.Group("/ordens-servico")
		{
			os.GET("/", ordemServicoController.BuscarTodas)
			os.GET("/:id", ordemServicoController.BuscarPorID)
			os.GET("/numero/:numero", ordemServicoController.BuscarPorNumero)
			os.GET("/cliente/:clienteId", ordemServicoController.BuscarPorCliente)
			os.GET("/veiculo/:veiculoId", ordemServicoController.BuscarPorVeiculo)
			os.GET("/status/:status", ordemServicoController.BuscarPorStatus)
			os.POST("/", ordemServicoController.Criar)
			os.PUT("/:id", ordemServicoController.Atualizar)
			os.PATCH("/:id/status", ordemServicoController.AtualizarStatus)
			os.DELETE("/:id", ordemServicoController.Deletar)

			// Rotas para itens da OS
			os.GET("/:id/itens", ordemServicoController.BuscarItens)
			os.POST("/:id/itens", ordemServicoController.AdicionarItem)
			os.PUT("/:id/itens/:itemId", ordemServicoController.AtualizarItem)
			os.DELETE("/:id/itens/:itemId", ordemServicoController.RemoverItem)

			// Ações específicas
			os.POST("/:id/concluir", ordemServicoController.ConcluirOS)
			os.POST("/:id/cancelar", ordemServicoController.CancelarOS)
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
