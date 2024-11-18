package repo

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"net/rpc"

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
