package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Nome        string         `json:"nome" gorm:"not null;size:100" binding:"required"`
	Email       string         `json:"email" gorm:"not null;unique;size:100" binding:"required,email"`
	Senha       string         `json:"senha,omitempty" gorm:"not null;size:100" binding:"required"`
	Cargo       string         `json:"cargo" gorm:"size:20;default:'user'"`
	Ativo       bool           `json:"ativo" gorm:"default:true"` // Novo campo
	UltimoLogin *time.Time     `json:"ultimoLogin,omitempty"`
	CreatedAt   time.Time      `json:"criadoEm" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"atualizadoEm" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relação inversa (opcional)
	Funcionario *Funcionario `json:"funcionario,omitempty" gorm:"foreignKey:UsuarioID"`
}

func (u *Usuario) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Senha = string(hashedPassword)
	return nil
}

func (u *Usuario) BeforeUpdate(tx *gorm.DB) error {
	// Só criptografa a senha se ela foi alterada
	if tx.Statement.Changed("Senha") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Senha), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		tx.Statement.SetColumn("Senha", string(hashedPassword))
	}
	return nil
}

// CompareSenha verifica se a senha fornecida corresponde à senha hash armazenada
func (u *Usuario) CompareSenha(senha string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(senha))
	return err == nil
}

func (Usuario) TableName() string {
	return "usuarios"
}
