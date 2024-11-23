package walletRepository

import (
	"fmt"
	"main/internal/model/wallet"

	"gorm.io/gorm"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

func (r *WalletRepository) GetWalletByID(tx *gorm.DB, uid string) (wallet.Wallet, error) {
	const op = "internal.wallerRepository.GetWalletByID"
	resWallet := wallet.Wallet{ID: uid}
	if err := tx.First(&resWallet).Error; err != nil {
		return wallet.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}
	return resWallet, nil
}

func (r *WalletRepository) CreateWallet(tx *gorm.DB, balance int64) (string, error) {
	const op = "internal.wallerRepository.CreateWallet"
	creatingWaller := wallet.Wallet{Balance: balance}
	if err := tx.Create(&creatingWaller).Error; err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return creatingWaller.ID, nil
}

func (r *WalletRepository) UpdateWallet(tx *gorm.DB, uid string, balance int64) (string, error) {
	const op = "internal.wallerRepository.UpdateWallet"
	updatingWaller := wallet.Wallet{ID: uid}
	if err := tx.Model(&updatingWaller).Update("balance", balance).Error; err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return updatingWaller.ID, nil
}
