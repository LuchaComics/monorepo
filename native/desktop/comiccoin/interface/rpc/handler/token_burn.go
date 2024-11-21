package handler

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
)

type TokenBurnArgs struct {
	ChainID               uint16
	FromAccountAddress    *common.Address
	AccountWalletPassword string
	TokenID               *big.Int
}

type TokenBurnReply struct {
}

func (impl *ComicCoinRPCServer) TokenBurn(args *TokenBurnArgs, reply *TokenBurnReply) error {
	pass, secureErr := sstring.NewSecureString(args.AccountWalletPassword)
	if secureErr != nil {
		return secureErr
	}
	err := impl.tokenBurnService.Execute(
		context.Background(),
		args.ChainID,
		args.FromAccountAddress,
		pass,
		args.TokenID,
	)
	if err != nil {
		return err
	}

	// Fill reply pointer to send the data back
	*reply = TokenBurnReply{}
	return nil
}
