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

	"github.com/signintech/gopdf"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// Pre-Screening Service.

type PDFBuilderResponseDTO struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Content  []byte `json:"content"`
}

type CBFFBuilderRequestDTO struct {
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

type CBFFBuilder interface {
	GeneratePDF(dto *CBFFBuilderRequestDTO) (*PDFBuilderResponseDTO, error)
}

type cbffBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewCBFFBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) CBFFBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for cbff initializing...")
	_, err := os.Stat(cfg.PDFBuilder.CBFFTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &cbffBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.CBFFTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *cbffBuilder) GeneratePDF(r *CBFFBuilderRequestDTO) (*PDFBuilderResponseDTO, error) {
	var err error
	bdr.Logger.Debug("opening up template file", slog.String("file", bdr.PDFTemplateFilePath))

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: gopdf.Unit_PT, PageSize: gopdf.Rect{W: 841.89, H: 595.28}}) //595.28, 841.89 = A4
	pdf.AddPage()

	// DEVELOPER NOTE:
	// The `github.com/signintech/gopdf` library needs to access a `tff` file
	// to utilize to render font family in our PDF. Therefore the following set
	// of lines are going to populate the font family we will need to use,
	err = pdf.AddTTFFont("roboto", "./static/roboto/Roboto-Regular.ttf")
	if err != nil {
		panic(err) // For developer purpose only.
	}
	err = pdf.AddTTFFont("roboto-bold", "./static/roboto/Roboto-Bold.ttf")
	if err != nil {
		panic(err) // For developer purpose only.
	}

	//// Import page 1
	tpl1 := pdf.ImportPage(bdr.PDFTemplateFilePath, 1, "/MediaBox")
	pdf.UseImportedTemplate(tpl1, 0, 0, 841.89, 595.28) // Draw imported template onto page

	//// Generate page one contents.

	err = bdr.generateContent(r, &pdf)
	if err != nil {
		return nil, err
	}

	//// Add Page two

	pdf.AddPage()
	tpl2 := pdf.ImportPage(bdr.PDFTemplateFilePath, 2, "/MediaBox")
	pdf.UseImportedTemplate(tpl2, 0, 0, 841.89, 595.28) // Draw imported template onto page

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

func (bdr *cbffBuilder) generateContent(r *CBFFBuilderRequestDTO, pdf *gopdf.GoPdf) error {

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetFont("roboto", "", 20)
	pdf.SetXY(25, 30) // x coordinate specification
	pdf.Cell(nil, r.CPSRN)

	//
	// LEFT SIDE
	//

	pdf.SetFont("roboto-bold", "", 14)

	// ROW 1
	pdf.SetXY(320, 83)
	pdf.Cell(nil, fmt.Sprintf("%v", r.SubmissionDate.Day())) // Day
	pdf.SetXY(350, 83)
	pdf.Cell(nil, fmt.Sprintf("%v", int(r.SubmissionDate.Month()))) // Month (number)
	pdf.SetXY(380, 83)
	pdf.Cell(nil, fmt.Sprintf("%v", r.SubmissionDate.Year())) // Day

	// ROW 2
	pdf.SetXY(225, 110)
	pdf.Cell(nil, r.InspectorFirstName)
	pdf.SetXY(325, 110)
	pdf.Cell(nil, r.InspectorLastName)

	// ROW 3
	pdf.SetXY(60, 138)
	pdf.Cell(nil, r.InspectorStoreName)

	// RIGHT SIDE
	//

	// ROW 1
	pdf.SetXY(465, 83)
	pdf.Cell(nil, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(463, 110)
	pdf.Cell(nil, r.IssueVol)
	pdf.SetXY(555, 110)
	pdf.Cell(nil, r.IssueNo)
	if r.IssueCoverMonth < 13 && r.IssueCoverMonth > 0 {
		pdf.SetXY(652, 110)
		pdf.Cell(nil, fmt.Sprintf("%v", time.Month(int(r.IssueCoverMonth))))
	} else {
		pdf.SetXY(652, 110)
		pdf.Cell(nil, "-") // No cover year date.
	}
	if r.IssueCoverYear > 1 {
		pdf.SetXY(720, 110)
		if r.IssueCoverYear == 2 {
			pdf.Cell(nil, "1899 or before")
		} else {
			pdf.Cell(nil, fmt.Sprintf("%v", int(r.IssueCoverYear)))
		}
	} else { // No cover date year.
		pdf.SetXY(720, 110)
		pdf.Cell(nil, "-")
	}

	// ROW 3
	pdf.SetXY(500, 138)
	pdf.Cell(nil, r.PublisherName)

	//
	// RIGHT
	//

	// pdf.SetFont("Helvetica", "B", 14) // This controls the next text.

	// ROW 1 - Creases
	switch strings.ToLower(r.CreasesFinding) {
	case "pr":
		pdf.SetXY(253, 198)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 198)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 198)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 198)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 198)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 198)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 198)
		pdf.Cell(nil, "NM")
	default:
		return fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}

	// ROW 2 - Tears
	switch strings.ToLower(r.TearsFinding) {
	case "pr":
		pdf.SetXY(253, 221)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 221)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 221)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 221)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 221)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 221)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 221)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for tears finding")
	}

	// ROW 3 - Missing Parts
	switch strings.ToLower(r.MissingPartsFinding) {
	case "pr":
		pdf.SetXY(253, 245)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 245)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 245)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 245)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 245)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 245)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 245)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for missing parts finding")
	}

	// ROW 4 - Stains / Marks / Substances
	switch strings.ToLower(r.StainsFinding) {
	case "pr":
		pdf.SetXY(253, 269)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 269)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 269)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 269)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 269)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 269)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 269)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for stains finding")
	}

	// ROW 5 - Distortion / Colour
	switch strings.ToLower(r.DistortionFinding) {
	case "pr":
		pdf.SetXY(253, 293)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 293)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 293)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 293)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 293)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 293)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 293)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for distorion finding")
	}

	// ROW 6 - Paper Quality
	switch strings.ToLower(r.PaperQualityFinding) {
	case "pr":
		pdf.SetXY(253, 317)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 317)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 317)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 317)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 317)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 317)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 317)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for paper quality finding")
	}

	// ROW 7 - Spine / Staples
	switch strings.ToLower(r.SpineFinding) {
	case "pr":
		pdf.SetXY(253, 341)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 341)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 341)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 341)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 341)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 341)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 341)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value for paper quality finding")
	}

	// ROW 8 - Cover (Front & Back)
	switch strings.ToLower(r.CoverFinding) {
	case "pr":
		pdf.SetXY(253, 365)
		pdf.Cell(nil, "PR")
	case "fr":
		pdf.SetXY(310, 365)
		pdf.Cell(nil, "FR")
	case "gd":
		pdf.SetXY(361, 365)
		pdf.Cell(nil, "GD")
	case "vg":
		pdf.SetXY(412, 365)
		pdf.Cell(nil, "VG")
	case "fn":
		pdf.SetXY(468, 365)
		pdf.Cell(nil, "FN")
	case "vf":
		pdf.SetXY(523, 365)
		pdf.Cell(nil, "VF")
	case "nm":
		pdf.SetXY(572, 365)
		pdf.Cell(nil, "NM")
	default:
		return errors.New("missing value cover finding")
	}

	// ROW 9 - Shows signs of temp
	if r.ShowsSignsOfTamperingOrRestoration == true {
		pdf.SetXY(236, 390)
		pdf.Cell(nil, "X")
	} else {
		pdf.SetXY(281, 390)
		pdf.Cell(nil, "X")
	}

	pdf.SetFont("roboto-bold", "", 40)

	// For debugging purposes only.
	bdr.Logger.Debug("row 10",
		slog.Int64("GradingScale", int64(r.GradingScale)),
	)

	// ROW 10 - Grading
	switch r.GradingScale {
	case s_d.GradingScaleLetter:
		pdf.SetXY(490, 431)
		pdf.Cell(nil, strings.ToUpper(r.OverallLetterGrade))

		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			// pdf.SetFont("Helvetica", "B", 30) // Start subscript.
			pdf.SetXY(554, 420)
			pdf.Cell(nil, "+")
			// pdf.SetFont("Helvetica", "B", 40) // Resume the previous font.
		}
	case s_d.GradingScaleNumber:
		if r.OverallNumberGrade >= 10 {
			pdf.SetXY(500, 431)
			pdf.Cell(nil, fmt.Sprintf("%.0f", r.OverallNumberGrade)) // Force 0 decimal places.
		} else {
			pdf.SetXY(500, 431)
			pdf.Cell(nil, fmt.Sprintf("%.1f", r.OverallNumberGrade)) // Force 2 decimal places.
		}
	case s_d.GradingScaleCPSPercentage:
		pdf.SetXY(490, 431)
		pdf.Cell(nil, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	//
	// LEFT
	//

	pdf.SetFont("roboto-bold", "", 7)

	if len(r.SpecialNotes) > 638 {
		return errors.New("special notes length over 455")
	}

	// Handle `key issue`.
	if r.IsKeyIssue {
		if r.KeyIssue == s_d.KeyIssueOther {
			r.SpecialNotes = r.KeyIssueOther + " " + r.SpecialNotes
		} else {
			r.SpecialNotes = s_d.KeyIssueMap[int(r.KeyIssue)] + " " + r.KeyIssueDetail + ". " + r.SpecialNotes
		}
	}

	specialNotesLines := splitText(r.SpecialNotes, 50)

	if specialNote, ok := getElementAtIndex(specialNotesLines, 0); ok { // ROW 1
		pdf.SetXY(636, 198+8*0)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 1); ok {
		pdf.SetXY(636, 198+8*1)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 2); ok {
		pdf.SetXY(636, 198+8*2)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 3); ok {
		pdf.SetXY(636, 198+8*3)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 4); ok {
		pdf.SetXY(636, 198+8*4)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 5); ok {
		pdf.SetXY(636, 198+8*5)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 6); ok {
		pdf.SetXY(636, 198+8*6)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 7); ok {
		pdf.SetXY(636, 198+8*7)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 8); ok {
		pdf.SetXY(636, 198+8*8)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 9); ok {
		pdf.SetXY(636, 198+8*9)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 10); ok {
		pdf.SetXY(636, 198+8*10)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 11); ok {
		pdf.SetXY(636, 198+8*11)
		pdf.Cell(nil, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 12); ok { // ROW 13 - MAXIMUM
		pdf.SetXY(636, 198+8*12)
		pdf.Cell(nil, specialNote)
	}

	////////////////////////////////////////////////////////////////////////////

	if len(r.GradingNotes) > 638 {
		return errors.New("grading notes length over 638")
	}

	gradingNotesLines := splitText(r.GradingNotes, 50)

	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 0); ok { // ROW 1
		pdf.SetXY(636, 356+8*0)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 1); ok {
		pdf.SetXY(636, 356+8*1)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 2); ok {
		pdf.SetXY(636, 356+8*2)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 3); ok {
		pdf.SetXY(636, 356+8*3)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 4); ok {
		pdf.SetXY(636, 356+8*4)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 5); ok {
		pdf.SetXY(636, 356+8*5)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 6); ok {
		pdf.SetXY(636, 356+8*6)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 7); ok {
		pdf.SetXY(636, 356+8*7)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 8); ok {
		pdf.SetXY(636, 356+8*8)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 9); ok {
		pdf.SetXY(636, 356+8*9)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 10); ok {
		pdf.SetXY(636, 356+8*10)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 11); ok {
		pdf.SetXY(636, 356+8*11)
		pdf.Cell(nil, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 12); ok {
		pdf.SetXY(636, 356+8*12)
		pdf.Cell(nil, gradingNote)
	}
	return nil
}
