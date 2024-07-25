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

	"github.com/bartmika/timekit"
	"github.com/signintech/gopdf"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// CPS PEDIGREE COLLECTION

type PCBuilderRequestDTO struct {
	CPSRN                              string                          `bson:"cpsrn" json:"cpSrn"`
	Filename                           string                          `bson:"filename" json:"filename"`
	SubmissionDate                     time.Time                       `bson:"submission_date" json:"submission_date"`
	Item                               string                          `bson:"item" json:"item"`
	SeriesTitle                        string                          `bson:"series_title" json:"series_title"`
	IssueVol                           string                          `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string                          `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64                           `bson:"issue_cover_year" json:"issue_cover_year"`   // 1 = 'No Cover Date Year
	IssueCoverMonth                    int8                            `bson:"issue_cover_month" json:"issue_cover_month"` // 13 = label: "No Cover Date Month"
	PublisherName                      string                          `bson:"publisher_name" json:"publisher_name"`
	IsKeyIssue                         bool                            `bson:"is_key_issue" json:"is_key_issue"`
	KeyIssue                           int8                            `bson:"key_issue" json:"key_issue"`
	KeyIssueOther                      string                          `bson:"key_issue_other" json:"key_issue_other"`
	KeyIssueDetail                     string                          `bson:"key_issue_detail" json:"key_issue_detail"`
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

type PCBuilder interface {
	GeneratePDF(dto *PCBuilderRequestDTO) (*PDFBuilderResponseDTO, error)
}

type pcBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewPCBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) PCBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for pc initializing...")
	_, err := os.Stat(cfg.PDFBuilder.PCTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &pcBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.PCTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *pcBuilder) GeneratePDF(r *PCBuilderRequestDTO) (*PDFBuilderResponseDTO, error) {
	var err error
	bdr.Logger.Debug("opening up template file", slog.String("file", bdr.PDFTemplateFilePath))

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: gopdf.Unit_PT, PageSize: gopdf.Rect{W: 792, H: 612}}) // `US Letter` in landscape orientation.
	pdf.AddPage()

	// DEVELOPER NOTE:
	// The `github.com/signintech/gopdf` library needs to access a `tff` file
	// to utilize to render font family in our PDF. Therefore the following set
	// of lines are going to populate the font family we will need to use,
	// err = pdf.AddTTFFont("roboto", "./static/roboto/Roboto-Regular.ttf")
	// if err != nil {
	// 	panic(err) // For developer purpose only.
	// }
	// err = pdf.AddTTFFont("roboto-bold", "./static/roboto/Roboto-Bold.ttf")
	// if err != nil {
	// 	panic(err) // For developer purpose only.
	// }
	err = pdf.AddTTFFont("arial-bold", "./static/arial/ARIALBD.TTF")
	if err != nil {
		panic(err) // For developer purpose only.
	}
	err = pdf.AddTTFFont("arial", "./static/arial/ARIAL.TTF")
	if err != nil {
		panic(err) // For developer purpose only.
	}

	// DEVELOPERS NOTE:
	// We do not want to print the base page template because the organization
	// does not need it. Therefore the following lines of code are commented out.
	// You will need a file called "Pedigree.pdf" if you want to open this
	// template.

	tpl1 := pdf.ImportPage(bdr.PDFTemplateFilePath, 1, "/MediaBox") // Import single page
	pdf.UseImportedTemplate(tpl1, 0, 0, 792, 612)                   // Draw imported template onto page

	// DEVELOPERS NOTE:
	// This is how you generate the content.

	err = bdr.generateContent(r, &pdf)
	if err != nil {
		return nil, err
	}

	////
	//// Generate the file and save it to the file.
	////

	fileName := fmt.Sprintf("%s.pdf", r.CPSRN)
	filePath := fmt.Sprintf("%s/%s", bdr.DataDirectoryPath, fileName)

	err = pdf.WritePdf(filePath)
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

	/////

	// cleanup(bdr.DataDirectoryPath)

	////
	//// Return the generate invoice.
	////

	return &PDFBuilderResponseDTO{
		FileName: fileName,
		FilePath: filePath,
		Content:  bin,
	}, err
}

func (bdr *pcBuilder) generateContent(r *PCBuilderRequestDTO, pdf *gopdf.GoPdf) error {

	////
	//// FIRST BOX CONTENT
	////

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)

	pdf.SetFont("arial", "", 22)
	pdf.SetXY(-338, -32) // x coordinate specification

	// Set the transformation matrix to rotate 180 degrees
	pdf.Rotate(180, 190, 27) // Note: angle=180, x=190, y=27

	pdf.Cell(nil, r.CPSRN) // Print the text to the PDF.

	pdf.SetTextColor(0, 0, 0) // Set font color to black.
	pdf.RotateReset()         // Undue the transformation matrix to rotate 180 degrees

	//
	// LEFT SIDE
	//

	pdf.SetFont("arial-bold", "", 13)

	// ROW 1
	pdf.SetXY(219.2, 163.4)
	pdf.Cell(nil, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(224.6, 183)
	pdf.Cell(nil, r.IssueVol)
	pdf.SetXY(325.3, 183)
	pdf.Cell(nil, r.IssueNo)

	var issueCoverMonthYear string
	if r.IssueCoverMonth < 13 && r.IssueCoverMonth > 0 {
		issueCoverMonthYear = timekit.GetMonthAbbreviationByInt(int(r.IssueCoverMonth))
		if r.IssueCoverYear > 1 {
			if r.IssueCoverYear == 2 {
				issueCoverMonthYear += "1899 or before"
			} else {
				issueCoverMonthYear += fmt.Sprintf(" %v", int(r.IssueCoverYear)) // 1 whitespace is not an accident! They are required for our design.
			}
		} else { // No cover date year.
			issueCoverMonthYear += "-"
		}
	} else {
		issueCoverMonthYear = "- " // No cover year date.
		if r.IssueCoverYear > 1 {
			if r.IssueCoverYear == 2 {
				issueCoverMonthYear += "1899 or before"
			} else {
				issueCoverMonthYear += fmt.Sprintf(" %v", int(r.IssueCoverYear)) // 1 whitespace is not an accident! They are required for our design.
			}
		} else { // No cover date year.
			issueCoverMonthYear += "-"
		}
	}
	pdf.SetXY(443.4, 183)
	pdf.Cell(nil, issueCoverMonthYear)

	pdf.SetFont("arial-bold", "", 13)

	// ROW 3
	pdf.SetXY(282.4, 202.6)
	if r.KeyIssue == 1 {
		pdf.Cell(nil, r.KeyIssueOther)
	} else {
		pdf.Cell(nil, fmt.Sprintf("%v %v", constants.SubmissionKeyIssue[r.KeyIssue], r.KeyIssueDetail))
	}

	// // DEVELOPERS NOTE: THIS IS ONLY HERE FOR REFERENCE PURPOSES (DELETE WHEN READY TO NO LONGER US)
	// pdf.SetTextColor(178, 34, 34)     // TEST  (see: https://www.rapidtables.com/web/color/red-color.html)
	// pdf.SetFont("arial-bold", "", 65) // TEST
	// pdf.SetXY(635.4, 152)
	// pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))
	// pdf.SetTextColor(0, 0, 0) // Set font color to black.

	// ROW 10 - Grading
	switch r.GradingScale {
	case s_d.GradingScaleLetter:
		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("arial-bold", "", 50) // Adjust the new font for the grading scale.

			pdf.SetXY(635.4, 160)
			pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))

			pdf.SetFont("arial-bold", "", 36) // Adjust the new font for the grading scale.
			pdf.SetXY(710, 147)
			pdf.Cell(nil, "+")
		} else {
			pdf.SetFont("arial-bold", "", 65) // Adjust the new font for the grading scale.
			pdf.SetXY(635.4, 152)
			pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))
		}
	case s_d.GradingScaleNumber:
		pdf.SetFont("arial-bold", "", 65) // Adjust the new font for the grading scale.

		if r.OverallNumberGrade >= 10 {
			pdf.SetXY(646, 152)
			pdf.Cell(nil, fmt.Sprintf("%.0f", r.OverallNumberGrade)) // Force 0 decimal places.
		} else {
			pdf.SetXY(640, 152)
			pdf.Cell(nil, fmt.Sprintf("%.1f", r.OverallNumberGrade)) // Force 2 decimal places.
		}
	case s_d.GradingScaleCPSPercentage:
		pdf.SetFont("arial-bold", "", 50) // Adjust the new font for the grading scale.

		pdf.SetXY(635.4, 159)
		pdf.Cell(nil, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	////
	//// SECOND BOX
	////

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetTextColor(178, 34, 34) // Set font color to firebrick red. (see: https://www.rapidtables.com/web/color/red-color.html)

	pdf.SetFont("arial", "", 22)
	pdf.SetXY(-278, -515) // x coordinate specification

	// PREVIOUS VALUES:
	// (1.) -240, -535

	// Set the transformation matrix to rotate 180 degrees
	pdf.Rotate(180, 190, 27) // Note: angle=180, x=190, y=27

	pdf.Cell(nil, r.CPSRN) // Print the text to the PDF.

	pdf.SetTextColor(0, 0, 0) // Set font color to black.
	pdf.RotateReset()         // Undue the transformation matrix to rotate 180 degrees

	//
	// LEFT SIDE
	//

	pdf.SetFont("arial-bold", "", 8)

	// ROW 1
	pdf.SetXY(252.1, 293.1)
	pdf.Cell(nil, fmt.Sprintf("%v", r.SubmissionDate.Day())) // Day
	pdf.SetXY(269.8, 293.1)
	pdf.Cell(nil, fmt.Sprintf("%v", int(r.SubmissionDate.Month()))) // Month (number)
	pdf.SetXY(287.6, 293.1)
	pdf.Cell(nil, fmt.Sprintf("%v", r.SubmissionDate.Year()%100)) // Year (the %100 will cause us to return last two digits in the number, for example 2023 will return 23)

	// ROW 2
	pdf.SetXY(205.9, 304.3)
	pdf.Cell(nil, fmt.Sprintf("%v %v", r.InspectorFirstName, r.InspectorLastName))

	// ROW 3
	pdf.SetXY(129.1, 315.3)
	pdf.Cell(nil, r.InspectorStoreName)

	//
	// RIGHT SIDE
	//

	// ROW 1
	pdf.SetXY(323.8, 293.1)
	pdf.Cell(nil, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(323.5, 303.6)
	pdf.Cell(nil, r.IssueVol)
	pdf.SetXY(369.3, 303.6)
	pdf.Cell(nil, r.IssueNo)

	pdf.SetXY(416.4, 303.6)
	pdf.Cell(nil, issueCoverMonthYear)

	// ROW 3
	pdf.SetXY(339.6, 315.3)
	pdf.Cell(nil, r.PublisherName)

	//
	// RIGHT
	//

	pdf.SetLineWidth(1) // Set shape outline line.

	// ROW 1 - Creases
	switch strings.ToLower(r.CreasesFinding) {
	case "pr":
		pdf.Oval(219.2, 339.1, 219.2+9.2, 339.1+9.2)
	case "fr":
		pdf.Oval(246.5, 339.1, 246.5+9.2, 339.1+9.2)
	case "gd":
		pdf.Oval(271.7, 338.5, 271.7+9.2, 338.5+9.2)
	case "vg":
		pdf.Oval(296.9, 339.1, 296.9+9.2, 339.1+9.2)
	case "fn":
		pdf.Oval(323.6, 339.1, 323.6+9.2, 339.1+9.2)
	case "vf":
		pdf.Oval(348.8, 339.1, 348.8+9.2, 339.1+9.2)
	case "nm":
		pdf.Oval(374, 339.1, 374+9.2, 339.1+9.2)
	default:
		return fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}

	// ROW 2 - Tears
	switch strings.ToLower(r.TearsFinding) {
	case "pr":
		pdf.Oval(219.2, 349.6, 219.2+9.2, 349.6+9.2)
	case "fr":
		pdf.Oval(246.5, 349.6, 246.5+9.2, 349.6+9.2)
	case "gd":
		pdf.Oval(271.7, 349.6, 271.7+9.2, 349.6+9.2)
	case "vg":
		pdf.Oval(296.9, 349.6, 296.9+9.2, 349.6+9.2)
	case "fn":
		pdf.Oval(323.6, 349.6, 323.6+9.2, 349.6+9.2)
	case "vf":
		pdf.Oval(348.8, 349.6, 348.8+9.2, 349.6+9.2)
	case "nm":
		pdf.Oval(374, 349.6, 374+9.2, 349.6+9.2)
	default:
		return errors.New("missing value for tears finding")
	}

	// ROW 3 - Missing Parts
	switch strings.ToLower(r.MissingPartsFinding) {
	case "pr":
		pdf.Oval(219.2, 359.9, 219.2+9.2, 359.9+9.2)
	case "fr":
		pdf.Oval(246.5, 359.9, 246.5+9.2, 359.9+9.2)
	case "gd":
		pdf.Oval(271.7, 359.9, 271.7+9.2, 359.9+9.2)
	case "vg":
		pdf.Oval(296.9, 359.9, 296.9+9.2, 359.9+9.2)
	case "fn":
		pdf.Oval(323.6, 359.9, 323.6+9.2, 359.9+9.2)
	case "vf":
		pdf.Oval(348.8, 359.9, 348.8+9.2, 359.9+9.2)
	case "nm":
		pdf.Oval(374, 359.9, 374+9.2, 359.9+9.2)
	default:
		return errors.New("missing value for missing parts finding")
	}

	// ROW 4 - Stains / Marks / Substances
	switch strings.ToLower(r.StainsFinding) {
	case "pr":
		pdf.Oval(219.2, 369.6, 219.2+9.2, 369.6+9.2)
	case "fr":
		pdf.Oval(246.5, 369.6, 246.5+9.2, 369.6+9.2)
	case "gd":
		pdf.Oval(271.7, 369.6, 271.7+9.2, 369.6+9.2)
	case "vg":
		pdf.Oval(296.9, 369.6, 296.9+9.2, 369.6+9.2)
	case "fn":
		pdf.Oval(323.6, 369.6, 323.6+9.2, 369.6+9.2)
	case "vf":
		pdf.Oval(348.8, 369.6, 348.8+9.2, 369.6+9.2)
	case "nm":
		pdf.Oval(374, 369.6, 374+9.2, 369.6+9.2)
	default:
		return errors.New("missing value for stains finding")
	}

	// ROW 5 - Distortion / Colour
	switch strings.ToLower(r.DistortionFinding) {
	case "pr":
		pdf.Oval(219.2, 380.2, 219.2+9.2, 380.2+9.2)
	case "fr":
		pdf.Oval(246.5, 380.2, 246.5+9.2, 380.2+9.2)
	case "gd":
		pdf.Oval(271.7, 380.2, 271.7+9.2, 380.2+9.2)
	case "vg":
		pdf.Oval(296.9, 380.2, 296.9+9.2, 380.2+9.2)
	case "fn":
		pdf.Oval(323.6, 380.2, 323.6+9.2, 380.2+9.2)
	case "vf":
		pdf.Oval(348.8, 380.2, 348.8+9.2, 380.2+9.2)
	case "nm":
		pdf.Oval(374, 380.2, 374+9.2, 380.2+9.2)
	default:
		return errors.New("missing value for distorion finding")
	}

	// ROW 6 - Paper Quality
	switch strings.ToLower(r.PaperQualityFinding) {
	case "pr":
		pdf.Oval(219.2, 390.4, 219.2+9.2, 390.4+9.2)
	case "fr":
		pdf.Oval(246.5, 390.4, 246.5+9.2, 390.4+9.2)
	case "gd":
		pdf.Oval(271.7, 390.4, 271.7+9.2, 390.4+9.2)
	case "vg":
		pdf.Oval(296.9, 390.4, 296.9+9.2, 390.4+9.2)
	case "fn":
		pdf.Oval(323.6, 390.4, 323.6+9.2, 390.4+9.2)
	case "vf":
		pdf.Oval(348.8, 390.4, 348.8+9.2, 390.4+9.2)
	case "nm":
		pdf.Oval(374, 390.4, 374+9.2, 390.4+9.2)
	default:
		return errors.New("missing value for paper quality finding")
	}

	// ROW 7 - Spine / Staples
	switch strings.ToLower(r.SpineFinding) {
	case "pr":
		pdf.Oval(219.2, 400, 219.2+9.2, 400+9.2)
	case "fr":
		pdf.Oval(246.5, 400, 246.5+9.2, 400+9.2)
	case "gd":
		pdf.Oval(271.7, 400, 271.7+9.2, 400+9.2)
	case "vg":
		pdf.Oval(296.9, 400, 296.9+9.2, 400+9.2)
	case "fn":
		pdf.Oval(323.6, 400, 323.6+9.2, 400+9.2)
	case "vf":
		pdf.Oval(348.8, 400, 348.8+9.2, 400+9.2)
	case "nm":
		pdf.Oval(374, 400, 374+9.2, 400+9.2)
	default:
		return errors.New("missing value for paper quality finding")
	}

	// ROW 8 - Cover (Front & Back)
	switch strings.ToLower(r.CoverFinding) {
	case "pr":
		pdf.Oval(219.2, 410.5, 219.2+9.2, 410.5+9.2)
	case "fr":
		pdf.Oval(246.5, 410.5, 246.5+9.2, 410.5+9.2)
	case "gd":
		pdf.Oval(271.7, 410.5, 271.7+9.2, 410.5+9.2)
	case "vg":
		pdf.Oval(296.9, 410.5, 296.9+9.2, 410.5+9.2)
	case "fn":
		pdf.Oval(323.6, 410.5, 323.6+9.2, 410.5+9.2)
	case "vf":
		pdf.Oval(348.8, 410.5, 348.8+9.2, 410.5+9.2)
	case "nm":
		pdf.Oval(374, 410.5, 374+9.2, 410.5+9.2)
	default:
		return errors.New("missing value cover finding")
	}

	pdf.SetFont("arial-bold", "", 8)

	// ROW 9 - Shows signs of tampering
	if r.ShowsSignsOfTamperingOrRestoration == true {
		pdf.SetXY(212.3, 422.2)
		pdf.Cell(nil, "X")
	} else {
		pdf.SetXY(233.5, 422.2)
		pdf.Cell(nil, "X")
	}

	// // DEVELOPERS NOTE: THIS IS ONLY HERE FOR REFERENCE PURPOSES (DELETE WHEN READY TO NO LONGER US)
	// pdf.SetTextColor(178, 34, 34)     // TEST  (see: https://www.rapidtables.com/web/color/red-color.html)
	// pdf.SetFont("arial-bold", "", 20) // TEST
	// pdf.SetXY(332.1, 437.5)
	// pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))
	// pdf.SetTextColor(0, 0, 0) // Set font color to black.

	// ROW 10 - Grading
	switch r.GradingScale {
	case s_d.GradingScaleLetter:

		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("arial-bold", "", 16)
			pdf.SetXY(333, 440)
			pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))

			pdf.SetFont("arial-bold", "", 12) // Start subscript.
			pdf.SetXY(357, 436)
			pdf.Cell(nil, "+")
		} else {
			pdf.SetFont("arial-bold", "", 20)
			pdf.SetXY(332.1, 437.5)
			pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))
		}
	case s_d.GradingScaleNumber:
		if r.OverallNumberGrade >= 10 {
			pdf.SetFont("arial-bold", "", 20)
			pdf.SetXY(336, 437.5)
			pdf.Cell(nil, fmt.Sprintf("%v", r.OverallNumberGrade))
		} else {
			pdf.SetFont("arial-bold", "", 20)
			pdf.SetXY(334, 437.5)
			pdf.Cell(nil, fmt.Sprintf("%.1f", r.OverallNumberGrade))
		}
	case s_d.GradingScaleCPSPercentage:
		pdf.SetFont("arial-bold", "", 16)
		pdf.SetXY(332, 439)
		pdf.Cell(nil, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}
	pdf.SetFont("arial-bold", "", 20) // Resume the previous font.

	//
	// LEFT
	//

	pdf.SetFont("arial", "", 8) // Developer Note: Really small text!

	if len(r.SpecialNotes) > 200 { // Note: 25 characters per row * 8 rows = 200 characters max.
		r.SpecialNotes = r.SpecialNotes[:200] // Take the first 200 characters.
	}

	// DEVELOPERS NOTE:
	// The 'special notes' is too small to fit the signature and the user
	// inputted `special notes`; therefore, as a result if the system detects
	// any signatures then it'll be the only text printed in the `special notse`.
	// If there are no signatures then the user inputted `special notes` will
	// be printed in the PDF.

	if len(r.Signatures) >= 1 {
		ln1 := fmt.Sprintf("Signature of %v %v authenticated by CPS.", r.Signatures[0].Role, r.Signatures[0].Name)
		r.SpecialNotes = ln1
	}

	specialNotesLines := splitText(r.SpecialNotes, 17)

	// Begin rendering the special note lines onto the screen.

	if specialNote, ok := getElementAtIndex(specialNotesLines, 0); ok { // ROW 1
		pdf.SetXY(404, 342.8+6.75*0)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 1); ok {
		pdf.SetXY(404, 342.8+6.75*1)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 2); ok {
		pdf.SetXY(404, 342.8+6.75*2)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 3); ok {
		pdf.SetXY(404, 342.8+6.75*3)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 4); ok {
		pdf.SetXY(404, 342.8+6.75*4)
		pdf.Cell(nil, specialNote)
	}
	// if specialNote, ok := getElementAtIndex(specialNotesLines, 5); ok {
	// 	pdf.SetXY(404, 342.8+6.75*5)
	// 	pdf.Cell(nil, specialNote)
	// }
	// if specialNote, ok := getElementAtIndex(specialNotesLines, 6); ok {
	// 	pdf.SetXY(404, 342.8+6.75*6)
	// 	pdf.Cell(nil, specialNote)
	// }
	// if specialNote, ok := getElementAtIndex(specialNotesLines, 7); ok {
	// 	pdf.SetXY(404, 342.8+6.75*7)
	// 	pdf.Cell(nil, specialNote)
	// }

	////////////////////////////////////////////////////////////////////////////

	pdf.SetFont("arial", "", 8)

	if len(r.GradingNotes) > 200 {
		r.GradingNotes = r.GradingNotes[:200] // Take the first 200 characters.
	}

	gradingNotesLines := splitText(r.GradingNotes, 17)

	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 0); ok {
		pdf.SetXY(404, 408.4+6.75*0)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 1); ok {
		pdf.SetXY(404, 408.4+6.75*1)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 2); ok {
		pdf.SetXY(404, 408.4+6.75*2)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 3); ok {
		pdf.SetXY(404, 408.4+6.75*3)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 4); ok {
		pdf.SetXY(404, 408.4+6.75*4)
		pdf.Cell(nil, gradingNote)
	}
	// if gradingNote, ok := getElementAtIndex(gradingNotesLines, 5); ok {
	// 	pdf.SetXY(404, 408.4+6.75*5)
	// 	pdf.Cell(nil, gradingNote)
	// }
	// if gradingNote, ok := getElementAtIndex(gradingNotesLines, 6); ok {
	// 	pdf.SetXY(404, 408.4+6.75*6)
	// 	pdf.Cell(nil, gradingNote)
	// }
	// if gradingNote, ok := getElementAtIndex(gradingNotesLines, 7); ok {
	// 	pdf.SetXY(404, 408.4+6.75*7)
	// 	pdf.Cell(nil, gradingNote)
	// }

	return nil
}
