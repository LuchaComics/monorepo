package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/emailer/mailgun" // TODO: Remove
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// ComicSubmissionController Interface for submission business logic controller.
type ComicSubmissionController interface {
	Create(ctx context.Context, req *ComicSubmissionCreateRequestIDO) (*submission_s.ComicSubmission, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*submission_s.ComicSubmission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*submission_s.ComicSubmission, error)
	UpdateByID(ctx context.Context, req *ComicSubmissionUpdateRequestIDO) (*submission_s.ComicSubmission, error)
	ListByFilter(ctx context.Context, f *submission_s.ComicSubmissionPaginationListFilter) (*submission_s.ComicSubmissionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *submission_s.ComicSubmissionPaginationListFilter) ([]*submission_s.ComicSubmissionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*submission_s.ComicSubmission, error)
	SetCustomer(ctx context.Context, submissionID primitive.ObjectID, customerID primitive.ObjectID) (*submission_s.ComicSubmission, error)
	CreateComment(ctx context.Context, submissionID primitive.ObjectID, content string) (*submission_s.ComicSubmission, error)
	// CreateFileAttachment(ctx context.Context, req *ComicSubmissionFileAttachmentCreateRequestIDO) (*submission_s.ComicSubmission, error)
	GetQRCodePNGImage(ctx context.Context, payload string) ([]byte, error)
	GetQRCodePNGImageOfRegisteryURLByCPSRN(ctx context.Context, cpsrn string) ([]byte, error)
}

type ComicSubmissionControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Password              password.Provider
	CPSRN                 cpsrn.Provider
	CBFFBuilder           pdfbuilder.CBFFBuilder
	PCBuilder             pdfbuilder.PCBuilder
	CCIMGBuilder          pdfbuilder.CCIMGBuilder
	CCSCBuilder           pdfbuilder.CCSCBuilder
	CCBuilder             pdfbuilder.CCBuilder
	CCUGBuilder           pdfbuilder.CCUGBuilder
	Emailer               mg.Emailer // TODO: Remove
	TemplatedEmailer      templatedemailer.TemplatedEmailer
	Kmutex                kmutex.Provider
	DbClient              *mongo.Client
	UserStorer            user_s.UserStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
	StoreStorer           store_s.StoreStorer
	CreditStorer          credit_s.CreditStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	kmux kmutex.Provider,
	cpsrnP cpsrn.Provider,
	cbffb pdfbuilder.CBFFBuilder,
	pcb pdfbuilder.PCBuilder,
	ccimg pdfbuilder.CCIMGBuilder,
	ccsc pdfbuilder.CCSCBuilder,
	cc pdfbuilder.CCBuilder,
	ccug pdfbuilder.CCUGBuilder,
	emailer mg.Emailer, // TODO: Remove
	client *mongo.Client,
	te templatedemailer.TemplatedEmailer,
	usr_storer user_s.UserStorer,
	sub_storer submission_s.ComicSubmissionStorer,
	org_storer store_s.StoreStorer,
	credit_storer credit_s.CreditStorer,
) ComicSubmissionController {
	loggerp.Debug("submission controller initialization started...")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae`
	// r := &pdfbuilder.CBFFBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	SubmissionDate:                     time.Now(),
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverYear:                     2024, // 1 = 'No Cover Date Year
	// 	IssueCoverMonth:                    1,    // 13 = label: "No Cover Date Month"
	// 	PublisherName:                      "Some publisher",
	// 	SpecialNotes:                       text,
	// 	GradingNotes:                       text,
	// 	CreasesFinding:                     "NM",
	// 	TearsFinding:                       "PR",
	// 	MissingPartsFinding:                "PR",
	// 	StainsFinding:                      "PR",
	// 	DistortionFinding:                  "PR",
	// 	PaperQualityFinding:                "PR",
	// 	SpineFinding:                       "PR",
	// 	CoverFinding:                       "PR",
	// 	GradingScale:                       3, // GradingScaleLetter = 1 | GradingScaleNumber = 2 | GradingScaleCPSPercentage = 3
	// 	OverallNumberGrade:                 75,
	// 	CpsPercentageGrade:                 75,
	// 	ShowsSignsOfTamperingOrRestoration: false,
	// 	OverallLetterGrade:                 "NM",
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// 	InspectorFirstName:                      "Bartlomiej",
	// 	InspectorLastName:                       "Mika",
	// 	InspectorStoreName:                      "Mika Software Corporation",
	// 	IsKeyIssue:                         true,
	// 	KeyIssue:                           2, // 2=First appearance of
	// 	KeyIssueOther:                      "",
	// 	KeyIssueDetail:                     "Batman",
	// }
	// res, err := cbffb.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae`
	// r := &pdfbuilder.PCBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	SubmissionDate:                     time.Now(),
	// 	SeriesTitle:                        "The Amazing Spider-Man",
	// 	IssueVol:                           "1",
	// 	IssueNo:                            "375",
	// 	IssueCoverYear:                     1993,
	// 	IssueCoverMonth:                    3,
	// 	PublisherName:                      "Some publisher",
	// 	KeyIssue:                           2, // 2=1st appearance of
	// 	KeyIssueDetail:                     "Anne Weying (She-Venom)",
	// 	SpecialNotes:                       text,
	// 	GradingNotes:                       text,
	// 	CreasesFinding:                     "PR",
	// 	TearsFinding:                       "PR",
	// 	MissingPartsFinding:                "PR",
	// 	StainsFinding:                      "PR",
	// 	DistortionFinding:                  "PR",
	// 	PaperQualityFinding:                "PR",
	// 	SpineFinding:                       "PR",
	// 	CoverFinding:                       "PR",
	// 	GradingScale:                       3, // 1=GradingScaleLetter | 2=GradingScaleNumber | 3=GradingScaleCPSPercentage
	// 	CpsPercentageGrade:                 98,
	// 	ShowsSignsOfTamperingOrRestoration: false,
	// 	OverallLetterGrade:                 "NM",
	// 	OverallNumberGrade:                 10,
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// 	InspectorFirstName:                 "Bartlomiej",
	// 	InspectorLastName:                  "Mika",
	// 	InspectorStoreName:                 "Mika Software Corporation",
	// 	Signatures: []*submission_s.ComicSubmissionSignature{
	// 		{
	// 			Role: "Writer",
	// 			Name: "Frank Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Brian Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Zoe Herbert",
	// 		},
	// 	},
	// }
	// res, err := pcb.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCIMGBuilderRequestDTO{
	// 	CPSRN:                "788346-26649-1-1000",
	// 	SeriesTitle:          "Winter World",
	// 	IssueVol:             "Vol 1",
	// 	IssueNo:              "#1",
	// 	IssueCoverYear:       2023,
	// 	IssueCoverMonth:      1,
	// 	PublisherName:        "Some publisher",
	// 	PrimaryLabelDetails:  2, // 2=Regular Edition
	// 	InspectorStoreName: "Mika Software Corp.",
	// 	SpecialNotes:         "XXXXXXXX XXXXXXX XXXXXXXX",
	// 	Signatures: []*s_d.SubmissionSignature{
	// 		{
	// 			Role: "Writer",
	// 			Name: "Frank Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Brian Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Zoe Herbert",
	// 		},
	// 	},
	// }
	// res, err := ccimg.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCSCBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverYear:                     2023,
	// 	IssueCoverMonth:                    1,
	// 	PublisherName:                      "Some publisher",
	// 	PrimaryLabelDetails:                2, // 2=Regular Edition
	// 	GradingScale:                       1,
	// 	ShowsSignsOfTamperingOrRestoration: true,
	// 	OverallLetterGrade:                 "NM",
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// 	InspectorStoreName:               "Mika Software Corp.",
	// 	SpecialNotes:                       "XXXXXXXX XXXXXXX XXXXXXXX XXXXXXX XXXXXXXXX",
	// 	Signatures: []*s_d.SubmissionSignature{
	// 		{
	// 			Role: "Writer",
	// 			Name: "Frank Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Brian Herbert",
	// 		}, {
	// 			Role: "Writer",
	// 			Name: "Zoe Herbert",
	// 		},
	// 	},
	// }
	//
	// res, err := ccsc.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCBuilderRequestDTO{
	// 	CPSRN:                            "788346-26649-1-1000",
	// 	SeriesTitle:                      "Winter World",
	// 	IssueVol:                         "Vol 1",
	// 	IssueNo:                          "#1",
	// 	IssueCoverYear:                   2023,
	// 	IssueCoverMonth:                  1,
	// 	PublisherName:                    "Some publisher",
	// 	PrimaryLabelDetails:              2, // 2=Regular Edition
	// 	GradingScale:                     1,
	// 	OverallLetterGrade:               "vf",
	// 	IsOverallLetterGradeNearMintPlus: false,
	// 	OverallNumberGrade:               10,
	// 	CpsPercentageGrade:               100,
	// 	InspectorStoreName:             "Mika Software Corp.",
	// 	SpecialNotes:                     "XXXXXXXX XXXXXXX XXXXXXXX XXXXXXX XXXXXXXXX XXXXXXXXXXXX XXXXXXXXXXXXXX XXXXXXXXX XXXXXXXX XXXXXXX XXXXXXXX XXXXXXX XXXXXXXXX XXXXXXXXXXXX XXXXXXXXXXXXXX XXXXXXXXX XXXXXXXXXXXXXX XXXXXXXXX",
	// }
	// res, err := cc.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCUGBuilderRequestDTO{
	// 	CPSRN:                            "788346-26649-1-1000",
	// 	SeriesTitle:                      "Winter World",
	// 	IssueVol:                         "Vol 1",
	// 	IssueNo:                          "#1",
	// 	IssueCoverYear:                   2023,
	// 	IssueCoverMonth:                  1,
	// 	PublisherName:                    "Some publisher",
	// 	PrimaryLabelDetails:              2, // 2=Regular Edition
	// 	GradingScale:                     3, // 1=Letter 2=Number 3=CPS
	// 	OverallLetterGrade:               "vf",
	// 	IsOverallLetterGradeNearMintPlus: false,
	// 	OverallNumberGrade:               7,
	// 	CpsPercentageGrade:               100,
	// 	InspectorStoreName:             "Mika Software Corp.",
	// 	SpecialNotes:                     "XXXXXXXX XXXXXXX XXXXXXXX XXXXXXX XXXXXXXXX XXXXXXXXXXXX XXXXXXXXXXXXXX XXXXXXXXX",
	// }
	// res, err := ccug.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	// ------------------------------------------------------------------------//
	// ------------------------------------------------------------------------//
	// ------------------------------------------------------------------------//

	s := &ComicSubmissionControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Password:              passwordp,
		Kmutex:                kmux,
		CPSRN:                 cpsrnP,
		CBFFBuilder:           cbffb,
		PCBuilder:             pcb,
		CCIMGBuilder:          ccimg,
		CCSCBuilder:           ccsc,
		CCBuilder:             cc,
		CCUGBuilder:           ccug,
		Emailer:               emailer, // TODO: Remove
		TemplatedEmailer:      te,
		DbClient:              client,
		UserStorer:            usr_storer,
		ComicSubmissionStorer: sub_storer,
		StoreStorer:           org_storer,
		CreditStorer:          credit_storer,
	}
	s.Logger.Debug("submission controller initialized")
	return s
}
