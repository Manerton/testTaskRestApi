package wallet

const (
	WITHDRAW string = "WITHDRAW"
	DEPOSIT  string = "DEPOSIT"
)

type Wallet struct {
	ID      string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Balance int64  `json:"balance"`
}
