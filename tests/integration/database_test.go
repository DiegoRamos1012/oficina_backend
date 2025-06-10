package tests

import (
    "OficinaMecanica/database"
    "path/filepath"
    "testing"

    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
)

func init() {
    // Configurar viper para ler o arquivo .env na raiz do projeto
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    
    // Encontrar o diretório raiz do projeto (subindo da pasta tests)
    rootDir := filepath.Join("..")
    viper.AddConfigPath(rootDir)
    
    // Carregar o arquivo .env
    if err := viper.ReadInConfig(); err != nil {
        // Se não conseguir ler o arquivo, use as variáveis de ambiente padrão para o teste
        viper.Set("DB_HOST", "localhost")
        viper.Set("DB_PORT", "3306")
        viper.Set("DB_USER", "test_user")
        viper.Set("DB_PASSWORD", "test_password")
        viper.Set("DB_NAME", "oficina_mecanica_test")
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