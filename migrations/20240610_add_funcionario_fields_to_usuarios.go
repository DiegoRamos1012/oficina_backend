package migrations

import (
	"gorm.io/gorm"
)

func AddFuncionarioFieldsToUsuarios(db *gorm.DB) error {
	// Adiciona as colunas se não existirem
	if err := db.Exec(`ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS data_admissao TIMESTAMP NULL`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS status VARCHAR(50) NULL`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS ferias BOOLEAN DEFAULT FALSE`).Error; err != nil {
		return err
	}
	return nil
}
