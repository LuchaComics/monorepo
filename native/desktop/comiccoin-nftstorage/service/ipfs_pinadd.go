package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/password"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/usecase"
)

type IPFSPinAddService struct {
	config                 *config.Config
	logger                 *slog.Logger
	jwtProvider            jwt.Provider
	passwordProvider       password.Provider
	ipfsGetNodeIDUseCase   *usecase.IPFSGetNodeIDUseCase
	ipfsPinAddUsecase      *usecase.IPFSPinAddUseCase
	upsertPinObjectUseCase *usecase.UpsertPinObjectUseCase
}

func NewIPFSPinAddService(
	cfg *config.Config,
	logger *slog.Logger,
	jwtp jwt.Provider,
	passp password.Provider,
	uc1 *usecase.IPFSGetNodeIDUseCase,
	uc2 *usecase.IPFSPinAddUseCase,
	uc3 *usecase.UpsertPinObjectUseCase,
) *IPFSPinAddService {
	return &IPFSPinAddService{cfg, logger, jwtp, passp, uc1, uc2, uc3}
}

type IPFSPinAddRequestIDO struct {
	ApiKey      string `bson:"api_key" json:"api_key"`
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

// IPFSPinAddResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type IPFSPinAddResponseIDO struct {
	RequestID uint64            `bson:"requestid" json:"requestid"`
	Status    string            `bson:"status" json:"status"`
	Created   time.Time         `bson:"created,omitempty" json:"created,omitempty"`
	Delegates []string          `bson:"delegates" json:"delegates"`
	Info      map[string]string `bson:"info" json:"info"`
	CID       string            `bson:"cid" json:"cid"`
	Name      string            `bson:"name" json:"name"`
	Origins   []string          `bson:"origins" json:"origins"`
	Meta      map[string]string `bson:"meta" json:"meta"`
}

func (s *IPFSPinAddService) Execute(ctx context.Context, req *IPFSPinAddRequestIDO) (*IPFSPinAddResponseIDO, error) {
	//
	// STEP 1:
	// Simple validation.
	//

	e := make(map[string]string)

	if req.ApiKey == "" {
		e["api_key"] = "missing value"
	}
	if req.Filename == "" {
		e["filename"] = "missing value"
	}
	if req.ContentType == "" {
		e["content_type"] = "missing value"
	}
	if req.Data == nil {
		e["data"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Advanced validation.
	//

	apiKeyDecoded, err := s.jwtProvider.ProcessJWTToken(req.ApiKey)
	if err != nil {
		s.logger.Error("Failed processing JWT token",
			slog.Any("error", err))
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", fmt.Sprintf("bad formatting: %v", err))
	}
	apiKeyPayload := strings.Split(apiKeyDecoded, "@")
	if len(apiKeyPayload) < 2 {
		s.logger.Error("api_key - corrupted payload: bad structure")
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: bad structure")
	}
	if apiKeyPayload[0] == "" {
		s.logger.Error("api_key - corrupted payload: missing `chain_id`")
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: missing `chain_id`")
	}
	if apiKeyPayload[1] == "" {
		s.logger.Error("api_key - corrupted payload: missing `secret`")
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: missing `secret`")
	}
	chainID := apiKeyPayload[0]
	if chainID != fmt.Sprintf("%v", constants.ComicCoinChainID) {
		s.logger.Error("api_key - invalid: `chain_id` does not match mainnet value")
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "invalid: `chain_id` does not match mainnet value")
	}

	// Verify the api key secret and project hashed secret match.
	passwordMatch, _ := s.passwordProvider.ComparePasswordAndHash(apiKeyPayload[1], s.config.App.AppSecret.String())
	if passwordMatch == false {
		s.logger.Error("password - does not match")
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "unauthorized")
	}

	//
	// STEP 2:
	// Check to see if we are able to send to our IPFS node. If not abandon
	// this execution immediately.
	//

	nodeID, err := s.ipfsGetNodeIDUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting ID from the IPFS node we are using", slog.Any("error", err))
		return nil, err
	}

	s.logger.Debug("IPFS node is connection to us",
		slog.Any("node_id", nodeID))

	//
	// STEP 3:
	// Submit to IPFS and get CID.
	//

	cid, err := s.ipfsPinAddUsecase.Execute(req.Data)
	if err != nil {
		s.logger.Error("Failed pin adding to IPFS.",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 4:
	// Define our object and save to our database.
	//

	ipAdress, _ := ctx.Value(constants.SessionIPAddress).(string)

	origins := make([]string, 0)
	meta := make(map[string]string, 0)

	// Handle meta. We will attach meta along with some custom fields.
	meta["filename"] = req.Filename
	meta["content_type"] = req.ContentType
	meta["content_length"] = fmt.Sprintf("%v", len(req.Data))

	// Create our meta record in the database.
	pinobj := &domain.PinObject{
		// Core fields required for a `pin` in IPFS.
		Status:    domain.StatusPinned,
		CID:       cid,
		RequestID: uint64(time.Now().UnixMilli()),
		Name:      "", // Blank b/c it's optional.
		Created:   time.Now(),
		Origins:   origins,
		Meta:      meta,
		Delegates: make([]string, 0),
		Info:      make(map[string]string, 0),

		// Extension (a.k.a. not part of the IPFS spec).
		Filename: req.Filename,
		// ContentType:           req.Meta["content_type"],
		CreatedFromIPAddress:  ipAdress,
		ModifiedAt:            time.Now(),
		ModifiedFromIPAddress: ipAdress,
	}

	// Save to database.
	if err := s.upsertPinObjectUseCase.Execute(ctx, pinobj); err != nil {
		s.logger.Error("database create error",
			slog.Any("error", err))
		return nil, err
	}

	s.logger.Debug("Saved to local database",
		slog.Any("cid", cid),
		slog.Any("requestid", pinobj.RequestID))

	//
	// STEP 5:
	// Return a response which is helpful for the IPFS gateway spec.
	//

	res := &IPFSPinAddResponseIDO{
		RequestID: pinobj.RequestID,
		Status:    pinobj.Status,
		Created:   pinobj.Created,
		Delegates: pinobj.Delegates,
		Info:      pinobj.Info,
		CID:       pinobj.CID,
		Name:      pinobj.Name,
		Origins:   pinobj.Origins,
		Meta:      pinobj.Meta,
	}

	s.logger.Debug("Sending OK response back",
		slog.Any("RequestID", pinobj.RequestID),
		slog.Any("Status", pinobj.Status),
		slog.Any("Created", pinobj.Created),
		slog.Any("Delegates", pinobj.Delegates),
		slog.Any("Info", pinobj.Info),
		slog.Any("CID", pinobj.CID),
		slog.Any("Name", pinobj.Name),
		slog.Any("Origins", pinobj.Origins),
		slog.Any("Meta", pinobj.Meta))

	return res, nil
}
