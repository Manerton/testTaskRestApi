package postgresql

import (
	"fmt"
	"main/internal/model/wallet"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPosgreSQL(connectStr string) (*gorm.DB, error) {

	const op = "storage.postgresql.NewPostreSQL"

	// try to connect posgresql
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commad := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""
	db.Exec(commad)

	// migration models
	err = db.AutoMigrate(&wallet.Wallet{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
