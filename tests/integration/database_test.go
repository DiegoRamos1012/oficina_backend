package tests

import (
	"OficinaMecanica/database"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Configurar viper para ler o arquivo test.env
	viper.SetConfigName("test")
	viper.SetConfigType("env")

	// Adicionar o diretório atual como local de busca
	viper.AddConfigPath(".") // Adiciona o diretório atual

	// Encontrar o diretório raiz do projeto
	rootDir := filepath.Join("..", "..")
	viper.AddConfigPath(rootDir)

	// Adicionar também o diretório de integração
	integrationDir := filepath.Dir(".")
	viper.AddConfigPath(integrationDir)

	absPath, _ := filepath.Abs(".")
	fmt.Printf("Diretório atual: %s\n", absPath)

	// Carregar o arquivo test.env
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Erro ao ler o arquivo de configuração: %v\n", err)
		// Se não conseguir ler o arquivo, use as variáveis de ambiente padrão para o teste
		viper.Set("DB_HOST", "localhost")
		viper.Set("DB_PORT", "3306")
		viper.Set("DB_USER", "root")
		viper.Set("DB_PASSWORD", "A1b2c3d4-")
		viper.Set("DB_NAME", "oficinateste_bd")
	} else {
		configFile := viper.ConfigFileUsed()
		fmt.Printf("Arquivo de configuração carregado com sucesso: %s\n", configFile)
		fmt.Printf("Usando banco de dados: %s\n", viper.GetString("DB_NAME"))
	}

	// Permitir sobreescrever com variáveis de ambiente do sistema
	viper.AutomaticEnv()
}

func TestDatabaseConnection(t *testing.T) {
	// Não precisamos mais definir as variáveis de ambiente aqui
	// pois já foram definidas na função init()

	// Estabelecer conexão com o banco de dados
	db, err := database.ConnectDB()

	// Verificar se a conexão foi estabelecida sem erros
	if !assert.NoError(t, err, "Erro ao conectar ao banco de dados") {
		t.FailNow()
	}
	assert.NotNil(t, db, "A conexão com o banco de dados é nil")

	// Verificar se o banco está respondendo com uma consulta simples
	sqlDB, err := db.DB()
	if !assert.NoError(t, err, "Erro ao obter a conexão SQL do GORM") {
		t.FailNow()
	}

	// Testar se o banco está respondendo
	err = sqlDB.Ping()
	assert.NoError(t, err, "O banco de dados não está respondendo ao ping")
}

func TestDatabaseMigrations(t *testing.T) {
	// Não precisamos mais definir as variáveis de ambiente aqui
	// pois já foram definidas na função init()

	// Estabelecer conexão com o banco de dados
	db, err := database.ConnectDB()
	if !assert.NoError(t, err, "Erro ao conectar ao banco de dados") {
		t.FailNow()
	}

	// Executar migrações no banco de teste
	err = database.SetupMigrations(db)
	if !assert.NoError(t, err, "Erro ao executar as migrações") {
		t.FailNow()
	}

	// Verificar se as tabelas foram criadas
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?",
		viper.GetString("DB_NAME")).Scan(&count)

	// Deve haver pelo menos as tabelas básicas
	assert.True(t, count >= 5, "Deve haver pelo menos 5 tabelas no banco de dados")

	// Verificar se tabelas específicas existem
	var exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ?)",
		viper.GetString("DB_NAME"), "usuarios").Scan(&exists)
	assert.True(t, exists, "A tabela 'usuarios' deve existir")

	db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ?)",
		viper.GetString("DB_NAME"), "clientes").Scan(&exists)
	assert.True(t, exists, "A tabela 'clientes' deve existir")
}

// TestMain é usado para configurar e limpar o ambiente de teste
func TestMain(m *testing.M) {
	// Executar os testes
	code := m.Run()

	// Limpar o banco de dados após os testes
	cleanupTestDatabase()

	// Finalizar com o código de saída dos testes
	os.Exit(code)
}

// cleanupTestDatabase limpa o banco de dados após os testes
func cleanupTestDatabase() {
	// Estabelecer conexão com o banco de dados
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados para limpeza: %v\n", err)
		return
	}

	// Obter conexão SQL do GORM
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Erro ao obter conexão SQL: %v\n", err)
		return
	}

	// Garantir que a conexão seja fechada ao final
	defer sqlDB.Close()

	// Lista das tabelas para truncar (você pode adicionar mais conforme necessário)
	tables := []string{"usuarios", "clientes", "veiculos", "estoque", "funcionarios", "ordens_servico"}

	// Desabilitar verificações de chave estrangeira para permitir truncar tabelas
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Truncar cada tabela
	for _, table := range tables {
		result := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if result.Error != nil {
			fmt.Printf("Erro ao truncar tabela %s: %v\n", table, result.Error)
		}
	}

	// Reabilitar verificações de chave estrangeira
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	fmt.Println("Banco de dados de teste limpo com sucesso")
}
