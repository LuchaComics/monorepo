package controller

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/pdfbuilder"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ComicSubmissionControllerImpl) generateLabelPDF(sessCtx mongo.SessionContext, m *s_d.ComicSubmission) (*pdfbuilder.PDFBuilderResponseDTO, error) {
	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[m.PublisherName]
	if m.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = m.PublisherNameOther
	}

	switch m.ServiceType {
	case s_d.ServiceTypePreScreening:
		return c.generateFindingsFormPDF(sessCtx, m)
	case s_d.ServiceTypePedigree:
		c.Logger.Debug("beginning to generate `pedigree` pdf")
		// The next following lines of code will create the PDF file gnerator
		// request to be submitted into our PDF file generator to generate the data.
		r := &pdfbuilder.PCBuilderRequestDTO{
			CPSRN:                              m.CPSRN,
			Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SubmissionDate:                     time.Now(),
			SeriesTitle:                        m.SeriesTitle,
			IssueVol:                           m.IssueVol,
			IssueNo:                            m.IssueNo,
			IssueCoverYear:                     m.IssueCoverYear,
			IssueCoverMonth:                    m.IssueCoverMonth,
			PublisherName:                      publisherNameDisplay,
			SpecialNotes:                       m.SpecialNotes,
			GradingNotes:                       m.GradingNotes,
			CreasesFinding:                     m.CreasesFinding,
			TearsFinding:                       m.TearsFinding,
			MissingPartsFinding:                m.MissingPartsFinding,
			StainsFinding:                      m.StainsFinding,
			DistortionFinding:                  m.DistortionFinding,
			PaperQualityFinding:                m.PaperQualityFinding,
			SpineFinding:                       m.SpineFinding,
			CoverFinding:                       m.CoverFinding,
			ShowsSignsOfTamperingOrRestoration: m.ShowsSignsOfTamperingOrRestoration == 1,
			GradingScale:                       m.GradingScale,
			OverallLetterGrade:                 m.OverallLetterGrade,
			OverallNumberGrade:                 m.OverallNumberGrade,
			CpsPercentageGrade:                 m.CpsPercentageGrade,
			InspectorFirstName:                 m.InspectorFirstName,
			InspectorLastName:                  m.InspectorLastName,
			InspectorStoreName:                 m.StoreName,
			Signatures:                         m.Signatures,
			PrimaryLabelDetails:                m.PrimaryLabelDetails,
			PrimaryLabelDetailsOther:           m.PrimaryLabelDetailsOther,
			IsOverallLetterGradeNearMintPlus:   m.IsOverallLetterGradeNearMintPlus,
			KeyIssue:                           m.KeyIssue,
			KeyIssueOther:                      m.KeyIssueOther,
			KeyIssueDetail:                     m.KeyIssueDetail,
		}
		pdfResponse, err := c.PCBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
		c.Logger.Debug("finished generate `pedigree` pdf")
		return pdfResponse, nil
	case s_d.ServiceTypeCPSCapsule:
		c.Logger.Debug("beginning to generate `cc` pdf")
		r := &pdfbuilder.CCBuilderRequestDTO{
			CPSRN:                            m.CPSRN,
			Filename:                         fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SeriesTitle:                      m.SeriesTitle,
			IssueVol:                         m.IssueVol,
			IssueNo:                          m.IssueNo,
			IssueCoverYear:                   m.IssueCoverYear,
			IssueCoverMonth:                  m.IssueCoverMonth,
			PublisherName:                    publisherNameDisplay,
			SpecialNotes:                     m.SpecialNotes,
			PrimaryLabelDetails:              m.PrimaryLabelDetails,
			PrimaryLabelDetailsOther:         m.PrimaryLabelDetailsOther,
			GradingScale:                     m.GradingScale,
			OverallLetterGrade:               m.OverallLetterGrade,
			IsOverallLetterGradeNearMintPlus: m.IsOverallLetterGradeNearMintPlus,
			OverallNumberGrade:               m.OverallNumberGrade,
			CpsPercentageGrade:               m.CpsPercentageGrade,
			InspectorFirstName:               m.InspectorFirstName,
			InspectorLastName:                m.InspectorLastName,
			InspectorStoreName:               m.StoreName,
			Signatures:                       m.Signatures,
			// KeyIssue:                         m.KeyIssue,
			// KeyIssueOther:                    m.KeyIssueOther,
			// KeyIssueDetail:                   m.KeyIssueDetail,
		}
		pdfResponse, err := c.CCBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate cc` pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
		c.Logger.Debug("finished generate `cc` pdf")
		return pdfResponse, nil
	case s_d.ServiceTypeCPSCapsuleSignatureCollection:
		c.Logger.Debug("beginning to generate `ccsc` pdf")
		r := &pdfbuilder.CCSCBuilderRequestDTO{
			CPSRN:                            m.CPSRN,
			Filename:                         fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SeriesTitle:                      m.SeriesTitle,
			IssueVol:                         m.IssueVol,
			IssueNo:                          m.IssueNo,
			IssueCoverYear:                   m.IssueCoverYear,
			IssueCoverMonth:                  m.IssueCoverMonth,
			PublisherName:                    publisherNameDisplay,
			SpecialNotes:                     m.SpecialNotes,
			PrimaryLabelDetails:              m.PrimaryLabelDetails,
			PrimaryLabelDetailsOther:         m.PrimaryLabelDetailsOther,
			GradingScale:                     m.GradingScale,
			OverallLetterGrade:               m.OverallLetterGrade,
			IsOverallLetterGradeNearMintPlus: m.IsOverallLetterGradeNearMintPlus,
			OverallNumberGrade:               m.OverallNumberGrade,
			CpsPercentageGrade:               m.CpsPercentageGrade,
			InspectorFirstName:               m.InspectorFirstName,
			InspectorLastName:                m.InspectorLastName,
			InspectorStoreName:               m.StoreName,
			Signatures:                       m.Signatures,
			// KeyIssue:                         m.KeyIssue,
			// KeyIssueOther:                    m.KeyIssueOther,
			// KeyIssueDetail:                   m.KeyIssueDetail,
		}
		pdfResponse, err := c.CCSCBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate `ccsc` pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
		c.Logger.Debug("finished generate `ccsc` pdf")
		return pdfResponse, nil
	case s_d.ServiceTypeCPSCapsuleIndieMintGem:
		c.Logger.Debug("beginning to generate `ccimg` pdf")
		r := &pdfbuilder.CCIMGBuilderRequestDTO{
			CPSRN:                            m.CPSRN,
			Filename:                         fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SeriesTitle:                      m.SeriesTitle,
			IssueVol:                         m.IssueVol,
			IssueNo:                          m.IssueNo,
			IssueCoverYear:                   m.IssueCoverYear,
			IssueCoverMonth:                  m.IssueCoverMonth,
			PublisherName:                    publisherNameDisplay,
			SpecialNotes:                     m.SpecialNotes,
			PrimaryLabelDetails:              m.PrimaryLabelDetails,
			PrimaryLabelDetailsOther:         m.PrimaryLabelDetailsOther,
			GradingScale:                     m.GradingScale,
			OverallLetterGrade:               m.OverallLetterGrade,
			IsOverallLetterGradeNearMintPlus: m.IsOverallLetterGradeNearMintPlus,
			OverallNumberGrade:               m.OverallNumberGrade,
			CpsPercentageGrade:               m.CpsPercentageGrade,
			InspectorFirstName:               m.InspectorFirstName,
			InspectorLastName:                m.InspectorLastName,
			InspectorStoreName:               m.StoreName,
			Signatures:                       m.Signatures,
			// KeyIssue:                         m.KeyIssue,
			// KeyIssueOther:                    m.KeyIssueOther,
			// KeyIssueDetail:                   m.KeyIssueDetail,
		}
		pdfResponse, err := c.CCIMGBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
		c.Logger.Debug("finished generate `ccimg` pdf")
		return pdfResponse, nil
	case s_d.ServiceTypeCPSCapsuleYouGrade:
		c.Logger.Debug("beginning to generate `ccug` pdf")
		r := &pdfbuilder.CCUGBuilderRequestDTO{
			CPSRN:                            m.CPSRN,
			Filename:                         fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SeriesTitle:                      m.SeriesTitle,
			IssueVol:                         m.IssueVol,
			IssueNo:                          m.IssueNo,
			IssueCoverYear:                   m.IssueCoverYear,
			IssueCoverMonth:                  m.IssueCoverMonth,
			PublisherName:                    publisherNameDisplay,
			SpecialNotes:                     m.SpecialNotes,
			PrimaryLabelDetails:              m.PrimaryLabelDetails,
			PrimaryLabelDetailsOther:         m.PrimaryLabelDetailsOther,
			GradingScale:                     m.GradingScale,
			OverallLetterGrade:               m.OverallLetterGrade,
			IsOverallLetterGradeNearMintPlus: m.IsOverallLetterGradeNearMintPlus,
			OverallNumberGrade:               m.OverallNumberGrade,
			CpsPercentageGrade:               m.CpsPercentageGrade,
			InspectorFirstName:               m.InspectorFirstName,
			InspectorLastName:                m.InspectorLastName,
			InspectorStoreName:               m.StoreName,
			Signatures:                       m.Signatures,
			// KeyIssue:                         m.KeyIssue,
			// KeyIssueOther:                    m.KeyIssueOther,
			// KeyIssueDetail:                   m.KeyIssueDetail,
		}
		pdfResponse, err := c.CCUGBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
		c.Logger.Debug("finished generate `ccug` pdf")
		return pdfResponse, nil
	default:
		return nil, fmt.Errorf("unsupported service-type via: %v", m.ServiceType)
	}
}

func (c *ComicSubmissionControllerImpl) generateFindingsFormPDF(sessCtx mongo.SessionContext, m *s_d.ComicSubmission) (*pdfbuilder.PDFBuilderResponseDTO, error) {
	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[m.PublisherName]
	if m.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = m.PublisherNameOther
	}

	c.Logger.Debug("beginning to generate `pre-screening` pdf")
	// The next following lines of code will create the PDF file gnerator
	// request to be submitted into our PDF file generator to generate the data.
	r := &pdfbuilder.CBFFBuilderRequestDTO{
		CPSRN:                              m.CPSRN,
		Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
		SubmissionDate:                     time.Now(),
		SeriesTitle:                        m.SeriesTitle,
		IssueVol:                           m.IssueVol,
		IssueNo:                            m.IssueNo,
		IssueCoverYear:                     m.IssueCoverYear,
		IssueCoverMonth:                    m.IssueCoverMonth,
		PublisherName:                      publisherNameDisplay,
		IsKeyIssue:                         m.IsKeyIssue,
		KeyIssue:                           m.KeyIssue,
		KeyIssueOther:                      m.KeyIssueOther,
		KeyIssueDetail:                     m.KeyIssueDetail,
		SpecialNotes:                       m.SpecialNotes,
		GradingNotes:                       m.GradingNotes,
		CreasesFinding:                     m.CreasesFinding,
		TearsFinding:                       m.TearsFinding,
		MissingPartsFinding:                m.MissingPartsFinding,
		StainsFinding:                      m.StainsFinding,
		DistortionFinding:                  m.DistortionFinding,
		PaperQualityFinding:                m.PaperQualityFinding,
		SpineFinding:                       m.SpineFinding,
		CoverFinding:                       m.CoverFinding,
		ShowsSignsOfTamperingOrRestoration: m.ShowsSignsOfTamperingOrRestoration == 1,
		GradingScale:                       m.GradingScale,
		OverallLetterGrade:                 m.OverallLetterGrade,
		IsOverallLetterGradeNearMintPlus:   m.IsOverallLetterGradeNearMintPlus,
		OverallNumberGrade:                 m.OverallNumberGrade,
		CpsPercentageGrade:                 m.CpsPercentageGrade,
		InspectorFirstName:                 m.InspectorFirstName,
		InspectorLastName:                  m.InspectorLastName,
		InspectorStoreName:                 m.StoreName,
		Signatures:                         m.Signatures,
		PrimaryLabelDetails:                m.PrimaryLabelDetails,
		PrimaryLabelDetailsOther:           m.PrimaryLabelDetailsOther,
	}
	pdfResponse, err := c.CBFFBuilder.GeneratePDF(r)
	if err != nil {
		c.Logger.Error("generate pdf error", slog.Any("error", err))
		return nil, err
	}
	if pdfResponse == nil {
		c.Logger.Error("generate pdf error does not return a response")
		return nil, errors.New("no response from pdf generator")
	}
	c.Logger.Debug("finished generate `pre-screening` pdf")
	return pdfResponse, nil
}

func (c *ComicSubmissionControllerImpl) generateAndUploadFindingsFormPDF(sessCtx mongo.SessionContext, m *s_d.ComicSubmission) (string, string, time.Time, error) {
	pdfResponse, err := c.generateFindingsFormPDF(sessCtx, m)
	if err != nil {
		c.Logger.Error("pdf generation error", slog.Any("error", err))
		return "", "", time.Now(), err
	}
	if pdfResponse == nil {
		c.Logger.Error("findings form pdf generation produced nothing")
		return "", "", time.Now(), fmt.Errorf("findings form pdf did not generate")
	}

	// The next few lines will upload our PDF to our remote storage. Once the
	// file is saved remotely, we will have a connection to it through a "key"
	// unique reference to the uploaded file.
	path := fmt.Sprintf("uploads/%v", pdfResponse.FileName)

	// Append `_findings_form.pdf` to name.
	path = strings.Replace(path, ".pdf", "_findings_form.pdf", 1)

	c.Logger.Debug("S3 will upload...",
		slog.String("path", path))

	err = c.S3.UploadContent(sessCtx, path, pdfResponse.Content)
	if err != nil {
		c.Logger.Error("s3 upload error", slog.Any("error", err))
		return "", "", time.Now(), err
	}

	c.Logger.Debug("S3 uploaded with success",
		slog.String("path", path))

	// The following will generate a pre-signed URL so user can download the file.
	expiryDate := time.Now().Add(time.Minute * 15)
	downloadableURL, err := c.S3.GetDownloadablePresignedURL(sessCtx, path, time.Minute*15)
	if err != nil {
		c.Logger.Error("s3 presign error", slog.Any("error", err))
		return "", "", time.Now(), err
	}

	// Removing local file from the directory and don't do anything if we have errors.
	if err := os.Remove(pdfResponse.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	return path, downloadableURL, expiryDate, nil
}

func (c *ComicSubmissionControllerImpl) generateAndUploadLabelPDF(sessCtx mongo.SessionContext, m *s_d.ComicSubmission) (string, string, time.Time, error) {
	pdfResponse, err := c.generateLabelPDF(sessCtx, m)
	if err != nil {
		c.Logger.Error("pdf generation error", slog.Any("error", err))
		return "", "", time.Now(), err
	}
	if pdfResponse == nil {
		c.Logger.Error("findings form pdf generation produced nothing")
		return "", "", time.Now(), fmt.Errorf("findings form pdf did not generate")
	}

	// The next few lines will upload our PDF to our remote storage. Once the
	// file is saved remotely, we will have a connection to it through a "key"
	// unique reference to the uploaded file.
	path := fmt.Sprintf("uploads/%v", pdfResponse.FileName)

	c.Logger.Debug("S3 will upload...",
		slog.String("path", path))

	err = c.S3.UploadContent(sessCtx, path, pdfResponse.Content)
	if err != nil {
		c.Logger.Error("s3 upload error", slog.Any("error", err))
		return "", "", time.Now(), err
	}

	c.Logger.Debug("S3 uploaded with success",
		slog.String("path", path))

	// The following will generate a pre-signed URL so user can download the file.
	expiryDate := time.Now().Add(time.Minute * 15)
	downloadableURL, err := c.S3.GetDownloadablePresignedURL(sessCtx, path, time.Minute*15)
	if err != nil {
		c.Logger.Error("s3 presign error", slog.Any("error", err))
		return "", "", time.Now(), err
	}

	// Removing local file from the directory and don't do anything if we have errors.
	if err := os.Remove(pdfResponse.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	return path, downloadableURL, expiryDate, nil
}
