package repo

import (
	"context"
	"log"
	"log/slog"
	"net/rpc"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ComicCoincRPCClientRepo struct {
	logger *slog.Logger
}

func NewComicCoincRPCClientRepo(logger *slog.Logger) domain.ComicCoincRPCClientRepository {
	return &ComicCoincRPCClientRepo{logger}
}

type Args struct{}

func (r *ComicCoincRPCClientRepo) GetTimestamp(ctx context.Context) (uint64, error) {
	var reply uint64
	args := Args{}
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call("ComicCoinRPCServer.GiveServerTimestamp", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	return reply, nil
}
