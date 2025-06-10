package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"OficinaMecanica/configs"
	"OficinaMecanica/models"
)

// Variável que pode ser substituída em testes
var GerarTokenFn = gerarTokenImpl

// GerarToken gera um token JWT para um usuário
func GerarToken(usuario models.Usuario) (string, error) {
	return GerarTokenFn(usuario)
}

// Implementação real da geração de token
func gerarTokenImpl(usuario models.Usuario) (string, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":         usuario.ID,
		"nome de usuário": usuario.Nome,
		"cargo":           usuario.Cargo,
		"exp":             time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
		"issued_at":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}

func ValidarToken(tokenString string) (*jwt.Token, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
}

func ExtrairUserID(token *jwt.Token) (int, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false
	}

	return int(userID), true
}