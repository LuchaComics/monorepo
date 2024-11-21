package handler

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type ComicCoinRPCServer struct {
	logger                                *slog.Logger
	getAccountService                     *service.GetAccountService
	createAccountService                  *service.CreateAccountService
	accountListingByLocalWalletsService   *service.AccountListingByLocalWalletsService
	coinTransferService                   *service.CoinTransferService
	tokenGetService                       *service.TokenGetService
	tokenTransferService                  *service.TokenTransferService
	tokenBurnService                      *service.TokenBurnService
	getOrDownloadNonFungibleTokenService  *service.GetOrDownloadNonFungibleTokenService
	listBlockTransactionsByAddressService *service.ListBlockTransactionsByAddressService
}

func NewComicCoinRPCServer(
	logger *slog.Logger,
	s1 *service.GetAccountService,
	s2 *service.CreateAccountService,
	s3 *service.AccountListingByLocalWalletsService,
	s4 *service.CoinTransferService,
	s5 *service.TokenGetService,
	s6 *service.TokenTransferService,
	s7 *service.TokenBurnService,
	s8 *service.GetOrDownloadNonFungibleTokenService,
	s9 *service.ListBlockTransactionsByAddressService,
) *ComicCoinRPCServer {

	// Create a new RPC server instance.
	port := &ComicCoinRPCServer{
		logger:                                logger,
		getAccountService:                     s1,
		createAccountService:                  s2,
		accountListingByLocalWalletsService:   s3,
		coinTransferService:                   s4,
		tokenGetService:                       s5,
		tokenTransferService:                  s6,
		tokenBurnService:                      s7,
		getOrDownloadNonFungibleTokenService:  s8,
		listBlockTransactionsByAddressService: s9,
	}

	return port
}
