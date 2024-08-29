package controller

import (
	"context"

	"log/slog"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/adapter/storage/ipfs"
)

func (impl *IpfsGatewayControllerImpl) GetIpfsNodeInfo(ctx context.Context) (*ipfs_storage.IpfsNodeInfo, error) {
	info, err := impl.IPFS.Id(ctx)
	if err != nil {
		impl.Logger.Error("get content by cid via ipfs error", slog.Any("error", err))
		return nil, err
	}
	impl.Logger.Debug("", slog.Any("info", info))
	return info, err
}
