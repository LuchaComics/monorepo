package controller

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

// ComicSubmissionCreateRequestIDO represents the user submitted data into our
// system for what comicbook submission they want.
type ComicSubmissionCreateRequestIDO struct {
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
	CustomerID                         primitive.ObjectID                 `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
}

// comicSubmissionFromCreate takes the request and converts it into our apps datastructure.
func comicSubmissionFromCreate(req *ComicSubmissionCreateRequestIDO) *s_d.ComicSubmission {
	cs := &s_d.ComicSubmission{
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
		CustomerID:                         req.CustomerID,
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

// Create function submits the comic book submission into our database.
func (impl *ComicSubmissionControllerImpl) Create(ctx context.Context, req *ComicSubmissionCreateRequestIDO) (*s_d.ComicSubmission, error) {
	// DEVELOPERS NOTE:
	// Every submission needs to have a unique `CPS Registry Number` (CPRN)
	// generated. The following needs to happen to generate the unique CPRN:
	// 1. Make the `Create` function be `atomic` and thus lock this function.
	// 2. Count total submissions in system.
	// 3. Generate CPRN.
	// 4. Apply the CPRN to the submission.
	// 5. Unlock this `Create` function to be usable again by other calls
	impl.Kmutex.Lock("CPS-BACKEND-SUBMISSION-INSERTION") // Step 1
	defer func() {
		impl.Kmutex.Unlock("CPS-BACKEND-SUBMISSION-INSERTION") // Step 5
	}()

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

		m := comicSubmissionFromCreate(req) // Convert into our data-structure.

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

		// DEVELOPERS NOTE:
		// Every submission creation is dependent on the `role` of the logged in
		// user in our system; however, the root administrator has the ability to
		// assign whatever store you want.
		switch userRole {

		case u_d.UserRoleRoot:
			impl.Logger.Debug("admin picking custom store")
		case u_d.UserRoleRetailer:
			impl.Logger.Debug("retailer assigning their store (auto-assigning `store_id`)")
			m.StoreID = sessCtx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		case u_d.UserRoleCustomer:
			impl.Logger.Debug("customer picking custom store (auto-assigning `store_id`)")

			// Force the following fields for logged in customer accounts.
			m.StoreID = loggedInUser.StoreID
			m.CustomerID = loggedInUser.ID
			m.CustomerFirstName = loggedInUser.FirstName
			m.CustomerLastName = loggedInUser.LastName
		default:
			impl.Logger.Error("unsupported role", slog.Any("role", userRole))
			return nil, fmt.Errorf("unsupported role via: %v", userRole)
		}

		// Lookup the store.
		org, err := impl.StoreStorer.GetByID(sessCtx, m.StoreID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if org == nil {
			impl.Logger.Error("database get by id does not exist", slog.Any("store id", m.StoreID))
			return nil, fmt.Errorf("does not exist for store id: %v", m.StoreID)
		}

		// Lookup the store owner.
		orgOwner, err := impl.UserStorer.GetByID(sessCtx, org.CreatedByUserID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if orgOwner == nil {
			impl.Logger.Error("database get by id does not exist", slog.Any("store id", m.StoreID))
			return nil, fmt.Errorf("does not exist for created by id: %v", m.StoreID)
		}

		m.StoreID = org.ID
		m.StoreName = org.Name
		m.StoreSpecialCollection = org.SpecialCollection
		m.StoreTimezone = org.Timezone

		// Generate the new `CSPRN` code and classificiation code to use for this submission record.
		csprn, csprnClassification, err := impl.generateCSRPN(sessCtx, org.SpecialCollection, m.ServiceType, userRole)
		if err != nil {
			impl.Logger.Error("csprn generation error", slog.Any("error", err))
			return nil, err
		}
		m.CPSRN = csprn
		m.CPSRNClassification = csprnClassification

		// Add defaults.
		m.ID = primitive.NewObjectID()
		m.CreatedByUserID = sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		m.CreatedByUserRole = userRole
		m.CreatedAt = time.Now()
		m.ModifiedByUserID = sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		m.ModifiedByUserRole = userRole
		m.ModifiedAt = time.Now()
		m.SubmissionDate = time.Now()
		m.Item = fmt.Sprintf("%v, %v, %v", m.SeriesTitle, m.IssueVol, m.IssueNo)

		// Attach a copy of the inspector to our record.
		m.InspectorID = orgOwner.ID
		m.InspectorFirstName = orgOwner.FirstName
		m.InspectorLastName = orgOwner.LastName

		// Attach a copy of the customer to our record
		if !m.CustomerID.IsZero() {
			customer, err := impl.UserStorer.GetByID(sessCtx, m.CustomerID)
			if err != nil {
				impl.Logger.Error("get customer user error", slog.Any("error", err))
				return nil, err
			}
			m.CustomerID = customer.ID
			m.CustomerFirstName = customer.FirstName
			m.CustomerLastName = customer.LastName
		}

		//
		// Credit - Lookup the retailer admin that posted the comic book submission
		//          and check to see if they have any available credits to burn. If
		//          they do then we will burn it and they do not have to purchase
		//          from us.
		//

		switch userRole {
		case u_d.UserRoleRetailer, u_d.UserRoleCustomer:

			// STEP 1: Credits.

			// The following code will lookup to see if the retailer user has a
			// credit they can burn on this submission.
			credit, err := impl.CreditStorer.GetNextAvailable(sessCtx, userID, m.ServiceType)
			if err != nil {
				impl.Logger.Error("get next available credit error", slog.Any("error", err))
				return nil, err
			}

			if credit != nil {
				impl.Logger.Debug("found credit for submission",
					slog.Any("creditID", credit.ID),
					slog.Any("comicSubmissionID", m.ID))

				// Keep a record in the comic submission that a credit was burned
				// so the retailer partner does not need to purchase.
				m.CreditID = credit.ID

				// Claim the credit so the credit can no longer be used again.
				credit.Status = credit_s.StatusClaimed
				credit.ClaimedByComicSubmissionID = m.ID
				credit.ModifiedAt = time.Now()
				if err := impl.CreditStorer.UpdateByID(sessCtx, credit); err != nil {
					impl.Logger.Error("database credit update error", slog.Any("error", err))
					return nil, err
				}

				impl.Logger.Debug("applied credit to submission",
					slog.Any("creditID", credit.ID),
					slog.Any("comicSubmissionID", m.ID))
			}

			// STEP 2: Pre-Screening

			// If a retailer makes a submission for `pre-screening` submission
			// then we need to change the status specific to this case of `completed
			// by retailer partner`.
			if m.ServiceType == s_d.ServiceTypePreScreening {
				m.Status = s_d.StatusCompletedByRetailPartner
			}
		}

		// Save to our database.
		if err := impl.ComicSubmissionStorer.Create(sessCtx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		//
		// Generate `Findings Form`.
		//

		ffObjectKey, ffObjectURL, ffObjectURLExpiry, err := impl.generateAndUploadFindingsFormPDF(sessCtx, m)
		if err != nil {
			impl.Logger.Error("generate and upload findings form error error", slog.Any("error", err))
			return nil, err
		}
		m.FindingsFormObjectKey = ffObjectKey
		m.FindingsFormObjectURL = ffObjectURL
		m.FindingsFormObjectURLExpiry = ffObjectURLExpiry
		m.ModifiedAt = time.Now()

		//
		// Generate `Label` based on the `service type`.
		//

		lObjectKey, lObjectURL, lObjectURLExpiry, err := impl.generateAndUploadLabelPDF(sessCtx, m)
		if err != nil {
			impl.Logger.Error("generate and upload findings form error error", slog.Any("error", err))
			return nil, err
		}
		m.LabelObjectKey = lObjectKey
		m.LabelObjectURL = lObjectURL
		m.LabelObjectURLExpiry = lObjectURLExpiry
		m.ModifiedAt = time.Now()

		//
		// Update database
		//

		if err := impl.ComicSubmissionStorer.UpdateByID(sessCtx, m); err != nil {
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
			m.LabelObjectKey = "[hidden]"
			m.LabelObjectURL = "[hidden]"
			m.LabelObjectURLExpiry = time.Now()
		}

		//
		// Send notification.
		//

		// The following code will send the email notifications to the correct
		// CPS staff individuals if the submission is NOT CBFF.
		if m.ServiceType != domain.ServiceTypePreScreening {
			if err := impl.sendNewComicSubmissionEmails(m); err != nil {
				impl.Logger.Error("database update error", slog.Any("error", err))
				// Do not return error, just keep it in the server logs.
			}
		}

		// return nil, httperror.NewForBadRequestWithSingleField("message", "halted by programmer") // For debugging purposes only.

		return m, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*submission_s.ComicSubmission), nil
}
