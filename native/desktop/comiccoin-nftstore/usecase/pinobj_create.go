package usecase

import (
	"log/slog"
	"mime/multipart"
	"time"
)

type CreatePinObjectUseCase struct {
	logger *slog.Logger
}

func NewCreatePinObjectUseCase(logger *slog.Logger) *CreatePinObjectUseCase {
	return &CreatePinObjectUseCase{logger}
}

type PinObjectCreateRequestIDO struct {
	Name    string
	Origins []string          `bson:"origins" json:"origins"`
	Meta    map[string]string `bson:"meta" json:"meta"`
	File    multipart.File    // Outside of IPFS pinning spec.
}

// PinObjectCreateResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type PinObjectCreateResponseIDO struct {
	RequestID uint64            `bson:"requestid" json:"requestid"`
	Status    string            `bson:"status" json:"status"`
	Created   time.Time         `bson:"created,omitempty" json:"created,omitempty"`
	Delegates []string          `bson:"delegates" json:"delegates"`
	Info      map[string]string `bson:"info" json:"info"`
}

func (uc *CreatePinObjectUseCase) Execute(req *PinObjectCreateRequestIDO) (*PinObjectCreateResponseIDO, error) {
	return nil, nil
	// //
	// // STEP 1:
	// // Validation.
	// //
	//
	// e := make(map[string]string)
	//
	// if req.Meta == nil {
	// 	e["meta"] = "missing value"
	// } else {
	// 	if req.Meta["filename"] == "" {
	// 		e["meta"] = "missing `filename` value"
	// 	}
	// 	if req.Meta["content_type"] == "" {
	// 		e["meta"] = "missing `content_type` value"
	// 	}
	// }
	// if req.File == nil {
	// 	e["file"] = "missing value"
	// }
	// if len(e) != 0 {
	// 	return nil, httperror.NewForBadRequest(&e)
	// }
	//
	// //
	// // STEP 2:
	// // Define our object.
	// //
	//
	// // Create our meta record in the database.
	// // res := &domain.PinObject{
	// // 	// Core fields required for a `pin` in IPFS.
	// // 	Status:    domain.StatusPinned,
	// // 	CID:       req.CID,
	// // 	RequestID: time.Now().UnixMilli(),
	// // 	Name:      req.Name,
	// // 	Created:   time.Now(),
	// // 	Origins:   req.Origins,
	// // 	Meta:      req.Meta,
	// // 	Delegates: make([]string, 0),
	// // 	Info:      make(map[string]string, 0),
	// //
	// // 	// Extension
	// // 	TenantID:              orgID,
	// // 	TenantName:            orgName,
	// // 	TenantTimezone:        orgTimezone,
	// // 	ID:                    primitive.NewObjectID(),
	// // 	ProjectID:             req.ProjectID,
	// // 	CreatedFromIPAddress:  ipAdress,
	// // 	ModifiedAt:            time.Now(),
	// // 	ModifiedFromIPAddress: ipAdress,
	// //
	// // 	// S3
	// // 	Filename: req.Meta["filename"],
	// // 	// ObjectKey: objectKey,
	// // 	// ObjectURL: "",
	// // }
	//
	// // Save to database.
	// if err := uc.PinObjectStorer.Create(sessCtx, res); err != nil {
	// 	impl.Logger.Error("database create error", slog.Any("error", err))
	// 	return nil, err
	// }
	// return res, nil
}
