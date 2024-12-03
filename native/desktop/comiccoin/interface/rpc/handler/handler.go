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
	getByBlockTransactionTimestampService *service.GetByBlockTransactionTimestampService
	blockDataGetByHashService             *service.BlockDataGetByHashService
	tokenListByOwnerService               *service.TokenListByOwnerService
	exportWalletService                   *service.ExportWalletService
	importWalletService                   *service.ImportWalletService
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
	s10 *service.GetByBlockTransactionTimestampService,
	s11 *service.BlockDataGetByHashService,
	s12 *service.TokenListByOwnerService,
	s13 *service.ExportWalletService,
	s14 *service.ImportWalletService,
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
		getByBlockTransactionTimestampService: s10,
		blockDataGetByHashService:             s11,
		tokenListByOwnerService:               s12,
		exportWalletService:                   s13,
		importWalletService:                   s14,
	}

	return port
}
