package walletService

import (
	"fmt"
	"main/internal/model/wallet"
	"main/internal/walletRepository"

	"gorm.io/gorm"
)

type WalletService struct {
	db         *gorm.DB
	repository *walletRepository.WalletRepository
}

func NewWalletService(db *gorm.DB, rep *walletRepository.WalletRepository) *WalletService {
	return &WalletService{
		db:         db,
		repository: rep,
	}
}

func (s *WalletService) GetWalletBalance(uid string) (int64, error) {
	const op = "internal.walletService.GetWalletBalance"
	resWallet, err := s.repository.GetWalletByID(s.db, uid)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return resWallet.Balance, nil
}

func (s *WalletService) CreateWallet(balance int64) (string, error) {
	const op = "internal.walletService.CreateWallet"
	resWallet, err := s.repository.CreateWallet(s.db, balance)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resWallet, nil
}

func (s *WalletService) UpdateBalance(uid string, amount int64, typeOperation string) (string, error) {
	const op = "internal.walletService.UpdateBalance"
	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return "", fmt.Errorf("%s: not create transaction %w", op, err)
	}
	tx.Set("gorm:query_option", "FOR UPDATE")
	// check have wallet by id
	resWallet, err := s.repository.GetWalletByID(tx, uid)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("%s: ERROR %w", op, err)
	}

	// check correct operation
	if typeOperation == wallet.WITHDRAW && resWallet.Balance < amount {
		tx.Rollback()
		return "", fmt.Errorf("%s: Balance < amount", op)
	}

	if typeOperation == wallet.WITHDRAW {
		resWallet.Balance -= amount
	} else if typeOperation == wallet.DEPOSIT {
		resWallet.Balance += amount
	} else {
		tx.Rollback()
		return "", fmt.Errorf("%s: %s - unsupported type operation", op, typeOperation)
	}

	// Update balance
	updatedId, err := s.repository.UpdateWallet(tx, uid, resWallet.Balance)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("%s: error update %w", op, err)
	}

	tx.Commit()
	return updatedId, nil
}
