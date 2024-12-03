package domain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	auth_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ComicCoincRPCClient struct {
	// Empty
}

type ComicCoincRPCClientRepository interface {
	GetTimestamp(ctx context.Context) (uint64, error)

	GetNonFungibleToken(ctx context.Context, nftID *big.Int, directoryPath string) (*NonFungibleToken, error)

	GetAccount(ctx context.Context, accountAddress *common.Address) (*auth_domain.Account, error)

	CreateAccount(
		ctx context.Context,
		password *sstring.SecureString,
		passwordRepeated *sstring.SecureString,
		label string,
	) (*auth_domain.Account, error)
	AccountListingByLocalWallets(ctx context.Context) ([]*auth_domain.Account, error)

	CoinTransfer(
		ctx context.Context,
		chainID uint16,
		fromAccountAddress *common.Address,
		accountWalletPassword *sstring.SecureString,
		to *common.Address,
		value uint64,
		data []byte,
	) error

	GetToken(ctx context.Context, tokenID *big.Int) (*auth_domain.Token, error)

	TokenTransfer(
		ctx context.Context,
		chainID uint16,
		fromAccountAddress *common.Address,
		accountWalletPassword *sstring.SecureString,
		to *common.Address,
		tokenID *big.Int,
	) error

	TokenBurn(
		ctx context.Context,
		chainID uint16,
		fromAccountAddress *common.Address,
		accountWalletPassword *sstring.SecureString,
		tokenID *big.Int,
	) error

	ListBlockTransactionsByAddress(
		ctx context.Context,
		address *common.Address,
	) ([]*auth_domain.BlockTransaction, error)

	GetBlockDataByHash(ctx context.Context, hash string) (*auth_domain.BlockData, error)

	ListTokensByOwnerAddress(
		ctx context.Context,
		ownerAddress *common.Address,
	) ([]*auth_domain.Token, error)

	ExportWallet(ctx context.Context, accountAddress *common.Address, filepath string) error

	ImportWallet(ctx context.Context, walletFilepath string) error
}
