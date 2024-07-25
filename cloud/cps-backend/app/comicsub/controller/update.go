package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type ComicSubmissionUpdateRequestIDO struct {
	ID                                 primitive.ObjectID                 `bson:"id,omitempty" json:"id,omitempty"`
	StoreID                            primitive.ObjectID                 `bson:"store_id,omitempty" json:"store_id,omitempty"`
	ServiceType                        int8                               `bson:"service_type" json:"service_type"`
	SubmissionDate                     time.Time                          `bson:"submission_date" json:"submission_date"`
	SeriesTitle                        string                             `bson:"series_title" json:"series_title"`
	IssueVol                           string                             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string                             `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64                              `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8                               `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      int8                               `bson:"publisher_name" json:"publisher_name"`
	PublisherNameOther                 string                             `bson:"publisher_name_other" json:"publisher_name_other"`
	IsKeyIssue                         bool                               `bson:"is_key_issue" json:"is_key_issue"`
	KeyIssue                           int8                               `bson:"key_issue" json:"key_issue"`
	KeyIssueOther                      string                             `bson:"key_issue_other" json:"key_issue_other"`
	KeyIssueDetail                     string                             `bson:"key_issue_detail" json:"key_issue_detail"`
	IsInternationalEdition             bool                               `bson:"is_international_edition" json:"is_international_edition"`
	IsVariantCover                     bool                               `bson:"is_variant_cover" json:"is_variant_cover"`
	VariantCoverDetail                 string                             `bson:"variant_cover_detail" json:"variant_cover_detail"`
	Printing                           int8                               `bson:"printing" json:"printing"`
	PrimaryLabelDetails                int8                               `bson:"primary_label_details" json:"primary_label_details"`
	PrimaryLabelDetailsOther           string                             `bson:"primary_label_details_other" json:"primary_label_details_other"`
	SpecialNotes                       string                             `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string                             `bson:"grading_notes" json:"grading_notes"`
	CreasesFinding                     string                             `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string                             `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string                             `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string                             `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string                             `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string                             `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string                             `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string                             `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration int8                               `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8                               `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string                             `bson:"overall_letter_grade" json:"overall_letter_grade"`
	OverallNumberGrade                 float64                            `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64                            `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
	IsOverallLetterGradeNearMintPlus   bool                               `bson:"is_overall_letter_grade_near_mint_plus" json:"is_overall_letter_grade_near_mint_plus"`
	CollectibleType                    int8                               `bson:"collectible_type" json:"collectible_type"`
	Status                             int8                               `bson:"status" json:"status"`
	Signatures                         []*domain.ComicSubmissionSignature `bson:"signatures" json:"signatures,omitempty"`
}

func comicSubmissionFromModify(req *ComicSubmissionUpdateRequestIDO) *s_d.ComicSubmission {
	cs := &s_d.ComicSubmission{
		ID:                                 req.ID,
		StoreID:                            req.StoreID,
		ServiceType:                        req.ServiceType,
		SubmissionDate:                     req.SubmissionDate,
		SeriesTitle:                        req.SeriesTitle,
		IssueVol:                           req.IssueVol,
		IssueNo:                            req.IssueNo,
		IssueCoverYear:                     req.IssueCoverYear,
		IssueCoverMonth:                    req.IssueCoverMonth,
		PublisherName:                      req.PublisherName,
		PublisherNameOther:                 req.PublisherNameOther,
		IsKeyIssue:                         req.IsKeyIssue,
		KeyIssue:                           req.KeyIssue,
		KeyIssueOther:                      req.KeyIssueOther,
		KeyIssueDetail:                     req.KeyIssueDetail,
		IsInternationalEdition:             req.IsInternationalEdition,
		IsVariantCover:                     req.IsVariantCover,
		VariantCoverDetail:                 req.VariantCoverDetail,
		Printing:                           req.Printing,
		PrimaryLabelDetails:                req.PrimaryLabelDetails,
		PrimaryLabelDetailsOther:           req.PrimaryLabelDetailsOther,
		SpecialNotes:                       req.SpecialNotes,
		GradingNotes:                       req.GradingNotes,
		CreasesFinding:                     req.CreasesFinding,
		TearsFinding:                       req.TearsFinding,
		MissingPartsFinding:                req.MissingPartsFinding,
		StainsFinding:                      req.StainsFinding,
		DistortionFinding:                  req.DistortionFinding,
		PaperQualityFinding:                req.PaperQualityFinding,
		SpineFinding:                       req.SpineFinding,
		CoverFinding:                       req.CoverFinding,
		ShowsSignsOfTamperingOrRestoration: req.ShowsSignsOfTamperingOrRestoration,
		GradingScale:                       req.GradingScale,
		OverallLetterGrade:                 req.OverallLetterGrade,
		OverallNumberGrade:                 req.OverallNumberGrade,
		CpsPercentageGrade:                 req.CpsPercentageGrade,
		IsOverallLetterGradeNearMintPlus:   req.IsOverallLetterGradeNearMintPlus,
		CollectibleType:                    req.CollectibleType,
		Status:                             req.Status,
		Signatures:                         req.Signatures,
	}

	// Set defaults for indie mint gems.
	if cs.ServiceType == s_d.ServiceTypeCPSCapsuleIndieMintGem {
		cs.CreasesFinding = "nm"
		cs.TearsFinding = "nm"
		cs.MissingPartsFinding = "nm"
		cs.StainsFinding = "nm"
		cs.DistortionFinding = "nm"
		cs.PaperQualityFinding = "nm"
		cs.SpineFinding = "nm"
		cs.CoverFinding = "nm"
		cs.GradingScale = s_d.GradingScaleNumber
		cs.IsOverallLetterGradeNearMintPlus = false
		cs.OverallNumberGrade = 10
		cs.ShowsSignsOfTamperingOrRestoration = 2
	}
	return cs
}

func (impl *ComicSubmissionControllerImpl) UpdateByID(ctx context.Context, req *ComicSubmissionUpdateRequestIDO) (*domain.ComicSubmission, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// DEVELOPERS NOTE:
		// Every submission creation is dependent on the `role` of the logged in
		// user in our system so we need to extract it right away.
		userRole, _ := sessCtx.Value(constants.SessionUserRole).(int8)
		// userFirstName, _ := sessCtx.Value(constants.SessionUserFirstName).(string)
		// userLastName, _ := sessCtx.Value(constants.SessionUserLastName).(string)
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

		ns := comicSubmissionFromModify(req) // Convert into our data-structure.

		//
		// Fetch submission.
		//

		// Fetch the original submission.
		os, err := impl.ComicSubmissionStorer.GetByID(sessCtx, ns.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			impl.Logger.Warn("submission does not exist error", slog.Any("id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("id", fmt.Sprintf("submission does not exist for ID: %v", req.ID))
		}

		// Variable used to keep track of the current logged in user.
		loggedInUser, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if loggedInUser == nil {
			impl.Logger.Error("database get by id does not exist", slog.Any("user id", userID))
			return nil, fmt.Errorf("does not exist for logged in user by id: %v", userID)
		}

		//
		// Set store.
		//

		// DEVELOPERS NOTE:
		// Every submission creation is dependent on the `role` of the logged in
		// user in our system; however, the root administrator has the ability to
		// assign whatever store you want.
		switch userRole {
		case u_d.UserRoleRoot:
			impl.Logger.Debug("admin picking custom store")
		case u_d.UserRoleRetailer:
			impl.Logger.Debug("retailer assigning their store (auto-assigning `store_id`)")
			os.StoreID = sessCtx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		case u_d.UserRoleCustomer:
			impl.Logger.Debug("customer picking custom store (auto-assigning `store_id`)")

			// Force the following fields for logged in customer accounts.
			os.StoreID = loggedInUser.StoreID
			os.CustomerID = loggedInUser.ID
			os.CustomerFirstName = loggedInUser.FirstName
			os.CustomerLastName = loggedInUser.LastName
		default:
			impl.Logger.Error("unsupported role", slog.Any("role", userRole))
			return nil, fmt.Errorf("unsupported role via: %v", userRole)
		}

		// Lookup the store.
		org, err := impl.StoreStorer.GetByID(sessCtx, ns.StoreID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if org == nil {
			impl.Logger.Error("database get by id does not exist", slog.Any("store id", ns.StoreID))
			return nil, fmt.Errorf("does not exist for store id: %v", ns.StoreID)
		}

		// Lookup the store owner.
		orgOwner, err := impl.UserStorer.GetByID(sessCtx, org.CreatedByUserID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if orgOwner == nil {
			impl.Logger.Error("database get by id does not exist", slog.Any("created by user id", org.CreatedByUserID))
			return nil, fmt.Errorf("does not exist for created by id: %v", org.CreatedByUserID)
		}

		// Update the record.
		os.StoreID = org.ID
		os.StoreName = org.Name
		os.StoreSpecialCollection = org.SpecialCollection
		os.StoreTimezone = org.Timezone

		//
		// Update records in database.
		//

		// Modify our original submission.
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		os.ModifiedByUserRole = userRole
		os.ServiceType = ns.ServiceType
		os.SubmissionDate = ns.SubmissionDate
		os.Item = fmt.Sprintf("%v, %v, %v", ns.SeriesTitle, ns.IssueVol, ns.IssueNo)
		os.SeriesTitle = ns.SeriesTitle
		os.IssueVol = ns.IssueVol
		os.IssueNo = ns.IssueNo
		os.IssueCoverYear = ns.IssueCoverYear
		os.IssueCoverMonth = ns.IssueCoverMonth
		os.PublisherName = ns.PublisherName
		os.PublisherNameOther = ns.PublisherNameOther
		os.IsKeyIssue = ns.IsKeyIssue
		os.KeyIssue = ns.KeyIssue
		os.KeyIssueOther = ns.KeyIssueOther
		os.KeyIssueDetail = ns.KeyIssueDetail
		os.IsInternationalEdition = ns.IsInternationalEdition
		os.IsVariantCover = ns.IsVariantCover
		os.VariantCoverDetail = ns.VariantCoverDetail
		os.PrimaryLabelDetails = ns.PrimaryLabelDetails
		os.PrimaryLabelDetailsOther = ns.PrimaryLabelDetailsOther
		os.SpecialNotes = ns.SpecialNotes
		os.GradingNotes = ns.GradingNotes
		os.CreasesFinding = ns.CreasesFinding
		os.TearsFinding = ns.TearsFinding
		os.MissingPartsFinding = ns.MissingPartsFinding
		os.StainsFinding = ns.StainsFinding
		os.DistortionFinding = ns.DistortionFinding
		os.PaperQualityFinding = ns.PaperQualityFinding
		os.SpineFinding = ns.SpineFinding
		os.CoverFinding = ns.CoverFinding
		os.ShowsSignsOfTamperingOrRestoration = ns.ShowsSignsOfTamperingOrRestoration
		os.GradingScale = ns.GradingScale
		os.OverallLetterGrade = ns.OverallLetterGrade
		os.IsOverallLetterGradeNearMintPlus = ns.IsOverallLetterGradeNearMintPlus
		os.OverallNumberGrade = ns.OverallNumberGrade
		os.CpsPercentageGrade = ns.CpsPercentageGrade
		// os.UserFirstName = ns.UserFirstName     // NO NEED TO CHANGE AFTER FACT.
		// os.UserLastName = ns.UserLastName       // NO NEED TO CHANGE AFTER FACT.
		// os.UserStoreName = ns.UserStoreName // NO NEED TO CHANGE AFTER FACT.
		os.Item = fmt.Sprintf("%v, %v, %v", ns.SeriesTitle, ns.IssueVol, ns.IssueNo)
		os.Signatures = ns.Signatures

		// Attach a copy of the inspector to our record.
		os.InspectorID = orgOwner.ID
		os.InspectorFirstName = orgOwner.FirstName
		os.InspectorLastName = orgOwner.LastName

		// DEVELOPERS NOTE:
		// Enforce status protection based on user roles.
		switch userRole {
		case u_d.UserRoleRoot:
			os.Status = ns.Status
		}

		// Save to the database the modified submission.
		if err := impl.ComicSubmissionStorer.UpdateByID(sessCtx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		//
		// Delete previous files from object storage.
		//

		impl.Logger.Debug("Will delete previous uploaded findings form",
			slog.String("path", os.FindingsFormObjectKey))

		// Delete previous record from remote storage.
		if err := impl.S3.DeleteByKeys(sessCtx, []string{os.FindingsFormObjectKey}); err != nil {
			impl.Logger.Warn("object delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		impl.Logger.Debug("Will delete previous uploaded label",
			slog.String("path", os.LabelObjectKey))

		// Delete previous record from remote storage.
		if err := impl.S3.DeleteByKeys(sessCtx, []string{os.LabelObjectKey}); err != nil {
			impl.Logger.Warn("object delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		//
		// Generate `Findings Form`.
		//

		ffObjectKey, ffObjectURL, ffObjectURLExpiry, err := impl.generateAndUploadFindingsFormPDF(sessCtx, os)
		if err != nil {
			impl.Logger.Error("generate and upload findings form error error", slog.Any("error", err))
			return nil, err
		}
		os.FindingsFormObjectKey = ffObjectKey
		os.FindingsFormObjectURL = ffObjectURL
		os.FindingsFormObjectURLExpiry = ffObjectURLExpiry
		os.ModifiedAt = time.Now()

		//
		// Generate `Label` based on the `service type`.
		//

		lObjectKey, lObjectURL, lObjectURLExpiry, err := impl.generateAndUploadLabelPDF(sessCtx, os)
		if err != nil {
			impl.Logger.Error("generate and upload findings form error error", slog.Any("error", err))
			return nil, err
		}
		os.LabelObjectKey = lObjectKey
		os.LabelObjectURL = lObjectURL
		os.LabelObjectURLExpiry = lObjectURLExpiry
		os.ModifiedAt = time.Now()

		//
		// Signatures - Update `special notes` for the PDF.
		//

		if len(ns.Signatures) > 0 {
			var str string
			for _, s := range ns.Signatures {
				str += fmt.Sprintf("Signature of %v %v authenticated by CPS.", s.Role, s.Name)
			}
			os.SpecialNotes = fmt.Sprintf("%v %v", str, ns.SpecialNotes)
		}

		//
		// Update database
		//

		if err := impl.ComicSubmissionStorer.UpdateByID(sessCtx, os); err != nil {
			impl.Logger.Error("database update error", slog.Any("error", err))
			return nil, err
		}

		//
		// Security - Censor label data if the logged in user is retailer. We do
		//            this because if the retailer gets our label then they can
		//            print it themeselves!
		//

		switch userRole {
		case u_d.UserRoleRetailer, u_d.UserRoleCustomer:
			os.LabelObjectKey = "[hidden]"
			os.LabelObjectURL = "[hidden]"
			os.LabelObjectURLExpiry = time.Now()
		}

		return os, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*domain.ComicSubmission), nil
}
