package controller

import (
	"context"
	"log/slog"
	"time"

	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type IPFSGatewayResponseIDO struct {
	CID         string
	Filename    string
	ContentType string
	ModifiedAt  time.Time
	Content     []byte
}

func (impl *IpfsGatewayControllerImpl) GetByCID(ctx context.Context, cid string) (*pinobject_s.PinObject, error) {
	// Retrieve from our database the record for the specific id.
	m, err := impl.PinObjectStorer.GetByCID(ctx, cid)
	if err != nil {
		impl.Logger.Error("database get by cid error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		impl.Logger.Warn("does not exist", slog.String("cid", cid))
		return nil, httperror.NewForNotFoundWithSingleField("cid", "does not exist")
	}

	content, err := impl.IPFS.Get(ctx, m.CID)
	if err != nil {
		impl.Logger.Error("get content by cid via ipfs error", slog.Any("error", err))
		return nil, err
	}

	m.Content = content

	return m, err
}
