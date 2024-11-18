package handler

import "time"

type Args struct{}

func (impl *ComicCoinRPCServer) GiveServerTimestamp(args *Args, reply *uint64) error {
	// Fill reply pointer to send the data back
	*reply = uint64(time.Now().UTC().UnixMilli())
	return nil
}
