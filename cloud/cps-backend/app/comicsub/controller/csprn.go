package controller

import (
	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ComicSubmissionControllerImpl) generateCSRPN(sessCtx mongo.SessionContext, specialCollection int8, serviceType int8, userRole int8) (string, string, error) {
	//-----------
	// Algorithn:
	// 1. Classify the comic book by the CSRPN generator.
	// 2. Get a total count of the comics submissions which have assigned the
	//    particular classification from (1) in the database.
	// 3. Generate a CSPRN
	//-----------

	// Step 1

	csprnClassification, err := c.CPSRN.Classify(specialCollection, serviceType, userRole)
	if err != nil {
		c.Logger.Error("count all submissions error", slog.Any("error", err))
		return "", "", err
	}

	// Step 2

	f := &domain.ComicSubmissionPaginationListFilter{CPSRNClassification: csprnClassification}
	count, err := c.ComicSubmissionStorer.CountByFilter(sessCtx, f) // Step 2
	if err != nil {
		c.Logger.Error("count all submissions error", slog.Any("error", err))
		return "", "", err
	}

	// Step 3

	csprn, err := c.CPSRN.Generate(specialCollection, serviceType, userRole, count)
	if err != nil {
		c.Logger.Error("count all submissions error", slog.Any("error", err))
		return "", "", err
	}

	c.Logger.Debug("Generated CPSRN based on retailer and comic submission service type",
		slog.String("cpsrn", csprn),
		slog.Int64("role", int64(userRole)),
		slog.Int64("specialCollection", int64(specialCollection)),
		slog.Int64("serviceType", int64(serviceType)),
		slog.Int64("count", count))

	return csprn, csprnClassification, nil
}
