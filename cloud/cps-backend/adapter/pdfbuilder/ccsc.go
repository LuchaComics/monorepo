package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// Signature Collection.

type CCSCBuilderRequestDTO struct {
	CPSRN                              string                          `bson:"cpsrn" json:"cpSrn"`
	Filename                           string                          `bson:"filename" json:"filename"`
	SubmissionDate                     time.Time                       `bson:"submission_date" json:"submission_date"`
	Item                               string                          `bson:"item" json:"item"`
	SeriesTitle                        string                          `bson:"series_title" json:"series_title"`
	IssueVol                           string                          `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string                          `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64                           `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8                            `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      string                          `bson:"publisher_name" json:"publisher_name"`
	SpecialNotes                       string                          `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string                          `bson:"grading_notes" json:"grading_notes"`
	CreasesFinding                     string                          `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string                          `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string                          `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string                          `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string                          `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string                          `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string                          `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string                          `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration bool                            `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8                            `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string                          `bson:"overall_letter_grade" json:"overall_letter_grade"`
	IsOverallLetterGradeNearMintPlus   bool                            `bson:"is_overall_letter_grade_near_mint_plus" json:"is_overall_letter_grade_near_mint_plus"`
	OverallNumberGrade                 float64                         `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64                         `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
	InspectorFirstName                 string                          `bson:"inspector_first_name" json:"inspector_first_name"`
	InspectorLastName                  string                          `bson:"inspector_last_name" json:"inspector_last_name"`
	InspectorStoreName                 string                          `bson:"inspector_store_name" json:"inspector_store_name"`
	Signatures                         []*s_d.ComicSubmissionSignature `bson:"signatures" json:"signatures,omitempty"`
	PrimaryLabelDetails                int8                            `bson:"primary_label_details" json:"primary_label_details"`
	PrimaryLabelDetailsOther           string                          `bson:"primary_label_details_other" json:"primary_label_details_other"`
}

// CCSCBuilder interface for building the "CPS C-Capsule Indie Mint Gem" edition document.
type CCSCBuilder interface {
	GeneratePDF(dto *CCSCBuilderRequestDTO) (*PDFBuilderResponseDTO, error)
}

type ccscBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewCCSCBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) CCSCBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for CCSC initializing...")
	_, err := os.Stat(cfg.PDFBuilder.CCSCTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &ccscBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.CCSCTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *ccscBuilder) GeneratePDF(r *CCSCBuilderRequestDTO) (*PDFBuilderResponseDTO, error) {
	var err error

	// Open our PDF invoice template and create clone it for the PDF invoice we will be building with.
	pdf := gofpdf.New("P", "mm", "A4", "")
	tpl1 := gofpdi.ImportPage(pdf, bdr.PDFTemplateFilePath, 1, "/MediaBox")

	pdf.AddPage()

	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 210, 300)

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetFont("Courier", "", 12)

	// Set the transformation matrix to rotate 180 degrees
	pdf.TransformBegin()
	pdf.TransformRotate(180, 190, 27) // angle=180, x=190, y=27

	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)

	// Print the text
	pdf.Text(190, 27, r.CPSRN) // x=190, y=27

	pdf.SetTextColor(0, 0, 0) // Set font color to black.

	pdf.TransformEnd()

	//
	// TITLE
	//

	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetXY(60, 51)
	pdf.Cell(0, 0, fmt.Sprintf("%v %v", r.SeriesTitle, r.IssueNo))

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(115, 51)
	pdf.Cell(0, 0, r.PublisherName)

	//
	// LEFT SIDE
	//

	pdf.SetFont("Helvetica", "B", 14)

	// ROW 1
	pdf.SetXY(60, 60)
	pdf.Cell(0, 0, "Volume:")
	pdf.SetXY(81, 60)
	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.IssueVol))
	pdf.SetTextColor(0, 0, 0) // Set font color to black.

	var issueDate string = "-"
	if r.IssueCoverMonth < 13 && r.IssueCoverMonth > 0 {
		month := fmt.Sprintf("%v", time.Month(int(r.IssueCoverMonth)))
		if r.IssueCoverYear > 1 {
			if r.IssueCoverYear == 2 {
				issueDate = "1899 or before"
			} else {
				issueDate = fmt.Sprintf("%v %v", month, int(r.IssueCoverYear))
			}
		} else { // No cover date year.
			// Do nothing
		}
	} else {
		// No cover year date.
		// Do nothing.
	}
	pdf.SetXY(60, 66)
	pdf.Cell(0, 0, "Date:")
	pdf.SetXY(75, 66)
	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)
	pdf.Cell(0, 0, issueDate)
	pdf.SetTextColor(0, 0, 0) // Set font color to black.

	////
	//// RIGHT SIDE
	////

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)
	pdf.SetXY(115, 59)
	switch r.PrimaryLabelDetails {
	case s_d.PrimaryLabelDetailsOther:
		pdf.Cell(0, 0, r.PrimaryLabelDetailsOther)
	case s_d.PrimaryLabelDetailsRegularEdition:
		pdf.Cell(0, 0, "Regular Edition")
	case s_d.PrimaryLabelDetailsDirectEdition:
		pdf.Cell(0, 0, "Direct Edition")
	case s_d.PrimaryLabelDetailsNewsstandEdition:
		pdf.Cell(0, 0, "Newstand Edition")
	case s_d.PrimaryLabelDetailsVariantCover:
		pdf.Cell(0, 0, "Variant Cover")
	case s_d.PrimaryLabelDetailsCanadianPriceVariant:
		pdf.Cell(0, 0, "Canadian Price Variant")
	case s_d.PrimaryLabelDetailsFacsimile:
		pdf.Cell(0, 0, "Facsimile")
	case s_d.PrimaryLabelDetailsReprint:
		pdf.Cell(0, 0, "Reprint")
	default:
		return nil, fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}
	pdf.SetTextColor(0, 0, 0) // Set font color to black.

	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)

	// Special notes. Max 100 characters.
	pdf.SetFont("Helvetica", "", 6)
	pdf.SetXY(115, 65)
	pdf.Cell(0, 0, r.SpecialNotes)

	//
	// Signature
	//

	if len(r.Signatures) >= 1 {
		ln1 := fmt.Sprintf("Signature of %v %v authenticated by CPS.", r.Signatures[0].Role, r.Signatures[0].Name)
		pdf.SetXY(115, 65+3)
		pdf.Cell(0, 0, ln1)
	}
	if len(r.Signatures) >= 2 {
		ln2 := fmt.Sprintf("Signature of %v %v authenticated by CPS.", r.Signatures[1].Role, r.Signatures[1].Name)
		pdf.SetXY(115, 65+6)
		pdf.Cell(0, 0, ln2)
	}
	if len(r.Signatures) >= 3 {
		ln3 := fmt.Sprintf("Signature of %v %v authenticated by CPS.", r.Signatures[2].Role, r.Signatures[2].Name)
		pdf.SetXY(115, 65+9)
		pdf.Cell(0, 0, ln3)
	}

	pdf.SetTextColor(0, 0, 0) // Set font color to black.

	//
	// // Retailer store name.
	// pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)
	// pdf.SetXY(0, 73)
	// pdf.CellFormat(0, 0, r.InspectorStoreName, "", 0, "R", false, 0, "") // Draw the right-aligned text
	// // pdf.Cell(0, 0,)
	// pdf.SetTextColor(0, 0, 0) // Set font color to black.

	//
	// GRADING
	//

	switch r.GradingScale {
	case s_d.GradingScaleCPSPercentage:
		pdf.SetFont("Helvetica", "", 24)
		if r.CpsPercentageGrade <= 9 {
			pdf.SetXY(29, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		} else if r.CpsPercentageGrade <= 99 && r.CpsPercentageGrade > 9 {
			pdf.SetXY(27, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		} else {
			pdf.SetXY(24, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		}
	case s_d.GradingScaleNumber:
		if r.OverallNumberGrade == 10 {
			pdf.SetFont("Helvetica", "", 60)
			pdf.SetXY(21.5, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
		} else {
			pdf.SetFont("Helvetica", "", 45)
			pdf.SetXY(23, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%.1f", r.OverallNumberGrade)) // Force 2 decimal places.
		}
	case s_d.GradingScaleLetter:
		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("Helvetica", "", 30)
			pdf.SetXY(26, 57)
			pdf.Cell(0, 0, "NM")

			pdf.SetFont("Helvetica", "", 10)
			pdf.SetXY(22.5, 65)
			pdf.Cell(0, 0, "Near Mint Plus")

			pdf.SetFont("Helvetica", "B", 22) // Start subscript.
			pdf.SetXY(41, 50)
			pdf.Cell(0, 0, "+")
		} else {
			pdf.SetFont("Helvetica", "", 30)
			pdf.SetXY(27, 56)
			pdf.Cell(0, 0, strings.ToUpper(r.OverallLetterGrade))

			pdf.SetFont("Helvetica", "", 14)
			switch r.OverallLetterGrade {
			case "pr", "PR":
				fallthrough
			case "fr", "FR":
				fallthrough
			case "fn", "FN":
				fallthrough
			case "gd", "GD":
				// CASE 1 OF 2: One word description. (Ex: "Fine")
				pdf.SetXY(29, 65)
				pdf.Cell(0, 0, constants.SubmissionOverallLetterGrades[r.OverallLetterGrade])
			case "vg", "VG":
				fallthrough
			case "vf", "VF":
				fallthrough
			case "nm", "NM":
				// CASE 1 OF 2: Two word description. (Ex: "Very Fine")
				pdf.SetXY(23, 65)
				pdf.Cell(0, 0, constants.SubmissionOverallLetterGrades[r.OverallLetterGrade])
			}
		}

		// case s_d.GradingScaleNumber:
		// 	pdf.SetXY(171, 153.5)
		// 	pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
		// case s_d.GradingScaleCPSPercentage:
		// 	pdf.SetXY(171, 153.5)
		// 	pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	pdf.SetFont("Helvetica", "", 10) // Set back the previous font.

	////
	//// Generate the file and save it to the file.
	////

	fileName := fmt.Sprintf("%s.pdf", r.CPSRN)
	filePath := fmt.Sprintf("%s/%s", bdr.DataDirectoryPath, fileName)

	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return nil, err
	}

	////
	//// Open the file and read all the binary data.
	////

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bin, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	////
	//// Return the generate invoice.
	////

	return &PDFBuilderResponseDTO{
		FileName: fileName,
		FilePath: filePath,
		Content:  bin,
	}, err
}
