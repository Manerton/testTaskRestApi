package request

type RequestCreateWallet struct {
	Balance int64 `json:"balance"`
}

type RequestUpdateWallet struct {
	ID            string `json:"id"`
	TypeOperation string `json:"typeOperation"`
	Amount        int64  `json:"amount"`
}
