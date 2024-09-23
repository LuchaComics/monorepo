package controller

type AccountCreateRequestIDO struct {
	Name           string `json:"name"`
	WalletPassword string `json:"wallet_password"`
}

type AccountDetailResponseIDO struct {
	Name           string `json:"name"`
	WalletFilepath string `json:"wallet_filepath"`
	WalletAddress  string `json:"wallet_address"`
}
