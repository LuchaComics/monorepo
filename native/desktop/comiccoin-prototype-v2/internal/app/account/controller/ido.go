package controller

type AccountCreateRequestIDO struct {
	Name           string `json:"name"`
	WalletPassword string `json:"wallet_password"`
}

type AccountDetailResponseIDO struct {
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
}
