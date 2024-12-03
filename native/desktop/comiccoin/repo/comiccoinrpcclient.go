package repo

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"net/rpc"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	auth_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ComicCoincRPCClientRepoConfigurationProvider interface {
	GetAddress() string // Retrieves the remote IPFS service address
	GetPort() string    // Retrieves the API key for authentication
}

// ComicCoincRPCClientRepoConfigurationProviderImpl is a struct that implements
// ComicCoincRPCClientRepoConfigurationProvider for configuration details.
type ComicCoincRPCClientRepoConfigurationProviderImpl struct {
	address string
	port    string // API key for accessing IPFS service
}

func NewComicCoincRPCClientRepoConfigurationProvider(address string, port string) ComicCoincRPCClientRepoConfigurationProvider {
	// Defensive code: Enforce `address` is set at minimum.
	if address == "" {
		log.Fatal("Missing `address` parameter.")
	}
	return &ComicCoincRPCClientRepoConfigurationProviderImpl{
		address: address,
		port:    port,
	}
}

// GetAddress retrieves the remote IPFS service address.
func (impl *ComicCoincRPCClientRepoConfigurationProviderImpl) GetAddress() string {
	return impl.address
}

// GetPort retrieves the API key for IPFS service authentication.
func (impl *ComicCoincRPCClientRepoConfigurationProviderImpl) GetPort() string {
	return impl.port
}

type ComicCoincRPCClientRepo struct {
	config    ComicCoincRPCClientRepoConfigurationProvider
	logger    *slog.Logger
	rpcClient *rpc.Client
}

func NewComicCoincRPCClientRepo(config ComicCoincRPCClientRepoConfigurationProvider, logger *slog.Logger) domain.ComicCoincRPCClientRepository {

	client, err := rpc.DialHTTP("tcp", config.GetAddress()+":"+config.GetPort())
	if err != nil {
		log.Fatal("NewComicCoincRPCClientRepo: RPC Dialing:", err)
	}

	return &ComicCoincRPCClientRepo{config, logger, client}
}

func (r *ComicCoincRPCClientRepo) GetTimestamp(ctx context.Context) (uint64, error) {
	var reply uint64

	type Args struct{}

	args := Args{}

	// Execute the remote procedure call.
	if err := r.rpcClient.Call("ComicCoinRPCServer.GiveServerTimestamp", args, &reply); err != nil {
		log.Fatal("arith error:", err)
	}

	// Return response from server.
	return reply, nil
}

func (r *ComicCoincRPCClientRepo) GetNonFungibleToken(ctx context.Context, nftID *big.Int, directoryPath string) (*domain.NonFungibleToken, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type GetNonFungibleTokenArgs struct {
		NonFungibleTokenID *big.Int
		DirectoryPath      string
	}

	type GetNonFungibleTokenReply struct {
		NonFungibleToken *domain.NonFungibleToken
	}

	// Construct our request / response.
	args := GetNonFungibleTokenArgs{
		NonFungibleTokenID: nftID,
		DirectoryPath:      directoryPath,
	}
	var reply GetNonFungibleTokenReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.GetNonFungibleToken", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.NonFungibleToken, nil
}

func (r *ComicCoincRPCClientRepo) GetAccount(ctx context.Context, accountAddress *common.Address) (*auth_domain.Account, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type GetAccountArgs struct {
		AccountAddress *common.Address
	}

	type GetAccountReply struct {
		Account *auth_domain.Account
	}

	// Construct our request / response.
	args := GetAccountArgs{
		AccountAddress: accountAddress,
	}
	var reply GetAccountReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.GetAccount", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.Account, nil
}

func (r *ComicCoincRPCClientRepo) CreateAccount(
	ctx context.Context,
	password *sstring.SecureString,
	passwordRepeated *sstring.SecureString,
	label string,
) (*auth_domain.Account, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type CreateAccountArgs struct {
		Password         string
		PasswordRepeated string
		Label            string
	}

	type CreateAccountReply struct {
		Account *auth_domain.Account
	}

	// Construct our request / response.
	args := CreateAccountArgs{
		Label:            label,
		Password:         password.String(),
		PasswordRepeated: passwordRepeated.String(),
	}
	var reply CreateAccountReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.CreateAccount", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.Account, nil
}

func (r *ComicCoincRPCClientRepo) AccountListingByLocalWallets(ctx context.Context) ([]*auth_domain.Account, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type AccountListingByLocalWalletsArgs struct {
	}

	type AccountListingByLocalWalletsReply struct {
		Accounts []*auth_domain.Account
	}

	// Construct our request / response.
	args := AccountListingByLocalWalletsArgs{}
	var reply AccountListingByLocalWalletsReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.AccountListingByLocalWallets", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.Accounts, nil
}

func (r *ComicCoincRPCClientRepo) CoinTransfer(
	ctx context.Context,
	chainID uint16,
	fromAccountAddress *common.Address,
	accountWalletPassword *sstring.SecureString,
	to *common.Address,
	value uint64,
	data []byte,
) error {
	// Define our request / response here by copy and pasting from the server codebase.
	type CoinTransferArgs struct {
		ChainID               uint16
		FromAccountAddress    *common.Address
		AccountWalletPassword string
		To                    *common.Address
		Value                 uint64
		Data                  []byte
	}

	type CoinTransferReply struct {
	}

	// Construct our request / response.
	args := CoinTransferArgs{
		ChainID:               chainID,
		FromAccountAddress:    fromAccountAddress,
		AccountWalletPassword: accountWalletPassword.String(),
		To:                    to,
		Value:                 value,
		Data:                  data,
	}
	var reply CoinTransferReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.CoinTransfer", args, &reply)
	if callError != nil {
		return callError
	}

	return nil
}

func (r *ComicCoincRPCClientRepo) GetToken(ctx context.Context, tokenID *big.Int) (*auth_domain.Token, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type GetTokenArgs struct {
		TokenID *big.Int
	}

	type GetTokenReply struct {
		Token *auth_domain.Token
	}

	// Construct our request / response.
	args := GetTokenArgs{
		TokenID: tokenID,
	}
	var reply GetTokenReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.GetToken", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.Token, nil
}

func (r *ComicCoincRPCClientRepo) TokenTransfer(
	ctx context.Context,
	chainID uint16,
	fromAccountAddress *common.Address,
	accountWalletPassword *sstring.SecureString,
	to *common.Address,
	tokenID *big.Int,
) error {
	// Define our request / response here by copy and pasting from the server codebase.
	type TokenTransferArgs struct {
		ChainID               uint16
		FromAccountAddress    *common.Address
		AccountWalletPassword string
		To                    *common.Address
		TokenID               *big.Int
	}

	type TokenTransferReply struct {
	}

	// Construct our request / response.
	args := TokenTransferArgs{
		ChainID:               chainID,
		FromAccountAddress:    fromAccountAddress,
		AccountWalletPassword: accountWalletPassword.String(),
		To:                    to,
		TokenID:               tokenID,
	}
	var reply TokenTransferReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.TokenTransfer", args, &reply)
	if callError != nil {
		return callError
	}

	return nil
}

func (r *ComicCoincRPCClientRepo) TokenBurn(
	ctx context.Context,
	chainID uint16,
	fromAccountAddress *common.Address,
	accountWalletPassword *sstring.SecureString,
	tokenID *big.Int,
) error {
	// Define our request / response here by copy and pasting from the server codebase.
	type TokenBurnArgs struct {
		ChainID               uint16
		FromAccountAddress    *common.Address
		AccountWalletPassword string
		TokenID               *big.Int
	}

	type TokenBurnReply struct {
	}

	// Construct our request / response.
	args := TokenBurnArgs{
		ChainID:               chainID,
		FromAccountAddress:    fromAccountAddress,
		AccountWalletPassword: accountWalletPassword.String(),
		TokenID:               tokenID,
	}
	var reply TokenBurnReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.TokenBurn", args, &reply)
	if callError != nil {
		return callError
	}

	return nil
}

func (r *ComicCoincRPCClientRepo) ListBlockTransactionsByAddress(
	ctx context.Context,
	accountAddress *common.Address,
) ([]*auth_domain.BlockTransaction, error) {
	// Define our request / response here by copy and pasting from the server codebase.

	type ListBlockTransactionsByAddressArgs struct {
		AccountAddress *common.Address
	}

	type ListBlockTransactionsByAddressReply struct {
		BlockTransactions []*auth_domain.BlockTransaction
	}

	// Construct our request / response.
	args := ListBlockTransactionsByAddressArgs{
		AccountAddress: accountAddress,
	}
	var reply ListBlockTransactionsByAddressReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.ListBlockTransactionsByAddress", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.BlockTransactions, nil
}

func (r *ComicCoincRPCClientRepo) GetBlockDataByHash(
	ctx context.Context,
	hash string,
) (*auth_domain.BlockData, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type BlockDataGetByHashArgs struct {
		Hash string
	}

	type BlockDataGetByHashReply struct {
		BlockData *auth_domain.BlockData
	}

	// Construct our request / response.
	args := BlockDataGetByHashArgs{
		Hash: hash,
	}
	var reply BlockDataGetByHashReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.BlockDataGetByHash", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.BlockData, nil
}

func (r *ComicCoincRPCClientRepo) ListTokensByOwnerAddress(
	ctx context.Context,
	ownerAddress *common.Address,
) ([]*auth_domain.Token, error) {
	// Define our request / response here by copy and pasting from the server codebase.
	type TokensListByOwnerAddressArgs struct {
		OwnerAddress *common.Address
	}

	type TokensListByOwnerAddressReply struct {
		Tokens []*auth_domain.Token
	}

	// Construct our request / response.
	args := TokensListByOwnerAddressArgs{
		OwnerAddress: ownerAddress,
	}
	var reply TokensListByOwnerAddressReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.ListTokensByOwnerAddress", args, &reply)
	if callError != nil {
		return nil, callError
	}

	return reply.Tokens, nil
}

func (r *ComicCoincRPCClientRepo) ExportWallet(ctx context.Context, accountAddress *common.Address, filepath string) error {
	// Define our request / response here by copy and pasting from the server codebase.
	type ExportWalletArgs struct {
		AccountAddress *common.Address
		FilePath       string
	}

	type ExportWalletReply struct{}

	// Construct our request / response.
	args := ExportWalletArgs{
		AccountAddress: accountAddress,
		FilePath:       filepath,
	}
	var reply ExportWalletReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.ExportWallet", args, &reply)
	if callError != nil {
		return callError
	}

	return nil
}

func (r *ComicCoincRPCClientRepo) ImportWallet(ctx context.Context, filepath string) error {
	// Define our request / response here by copy and pasting from the server codebase.
	type ImportWalletArgs struct {
		FilePath string
	}

	type ImportWalletReply struct{}

	// Construct our request / response.
	args := ImportWalletArgs{
		FilePath: filepath,
	}
	var reply ImportWalletReply

	// Execute the remote procedure call.
	callError := r.rpcClient.Call("ComicCoinRPCServer.ImportWallet", args, &reply)
	if callError != nil {
		return callError
	}

	return nil
}
