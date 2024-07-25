package httptransport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	sub_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/controller"
	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) unmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_c.ComicSubmissionUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_c.ComicSubmissionUpdateRequestIDO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it. Useful for diagnosing problems with inputed JSON structure.

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(teeReader).Decode(&requestData) // [1]
	if err != nil {
		h.Logger.Error("decoding error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateUpdateRequest(&requestData); err != nil {
		h.Logger.Warn("validation error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return nil, err
	}

	return &requestData, nil
}

func ValidateUpdateRequest(dirtyData *sub_c.ComicSubmissionUpdateRequestIDO) error {
	e := make(map[string]string)

	// if dirtyData.ServiceType == 0 {
	// 	e["service_type"] = "missing value"
	// }
	if dirtyData.SeriesTitle == "" {
		e["series_title"] = "missing value"
	}
	if dirtyData.IssueVol == "" {
		e["issue_vol"] = "missing value"
	}
	if dirtyData.IssueNo == "" {
		e["issue_no"] = "missing value"
	}
	if dirtyData.IssueCoverYear <= 0 {
		e["issue_cover_year"] = "missing value"
	}
	if dirtyData.IssueCoverMonth <= 0 {
		e["issue_cover_month"] = "missing value"
	}
	if dirtyData.PublisherName < 1 || dirtyData.PublisherName > 27 {
		e["publisher_name"] = "missing choice"
	} else if dirtyData.PublisherName == 1 && dirtyData.PublisherNameOther == "" {
		e["publisher_name_other"] = "missing choice"
	}
	if dirtyData.CreasesFinding == "" {
		e["creases_finding"] = "missing choice"
	}
	if dirtyData.TearsFinding == "" {
		e["tears_finding"] = "missing choice"
	}
	if dirtyData.MissingPartsFinding == "" {
		e["missing_parts_finding"] = "missing choice"
	}
	if dirtyData.StainsFinding == "" {
		e["stains_finding"] = "missing choice"
	}
	if dirtyData.DistortionFinding == "" {
		e["distortion_finding"] = "missing choice"
	}
	if dirtyData.PaperQualityFinding == "" {
		e["paper_quality_finding"] = "missing choice"
	}
	if dirtyData.SpineFinding == "" {
		e["spine_finding"] = "missing choice"
	}
	if dirtyData.CoverFinding == "" {
		e["cover_finding"] = "missing choice"
	}
	if dirtyData.GradingScale <= 0 || dirtyData.GradingScale > 3 {
		e["grading_scale"] = "missing choice"
	} else {
		if dirtyData.OverallLetterGrade == "" && dirtyData.GradingScale == sub_s.GradingScaleLetter {
			e["overall_letter_grade"] = "missing value"
		}
		if dirtyData.OverallNumberGrade <= 0 && dirtyData.OverallNumberGrade > 10 && dirtyData.GradingScale == sub_s.GradingScaleNumber {
			e["overall_number_grade"] = "missing value"
		}
		if dirtyData.CpsPercentageGrade < 5 && dirtyData.CpsPercentageGrade > 100 && dirtyData.GradingScale == sub_s.GradingScaleCPSPercentage {
			e["cps_percentage_grade"] = "missing value"
		}
	}
	if dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.YesItShowsSignsOfTamperingOrRestoration && dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.NoItDoesNotShowsSignsOfTamperingOrRestoration {
		e["shows_signs_of_tampering_or_restoration"] = "missing value"
	}

	// Process optional validation for `Special Notes` based on PDF requirements (see `adapter/pdfbuilder` package).
	if dirtyData.SpecialNotes != "" {
		if dirtyData.ServiceType == sub_s.ServiceTypePreScreening || dirtyData.ServiceType == sub_s.ServiceTypePedigree {
			if len(dirtyData.SpecialNotes) > 638 {
				e["special_notes"] = "over 638 characters"
			}
		}
		if dirtyData.ServiceType == sub_s.ServiceTypeCPSCapsule {
			if len(dirtyData.SpecialNotes) > 172 {
				e["special_notes"] = "over 172 characters"
			}
		}
		if dirtyData.ServiceType == sub_s.ServiceTypeCPSCapsuleYouGrade {
			if len(dirtyData.SpecialNotes) > 100 {
				e["special_notes"] = "over 100 characters"
			}
		}
		if dirtyData.ServiceType == sub_s.ServiceTypeCPSCapsuleSignatureCollection {
			if len(dirtyData.SpecialNotes) > 43 {
				e["special_notes"] = "over 43 characters"
			}
		}
		if dirtyData.ServiceType == sub_s.ServiceTypeCPSCapsuleIndieMintGem {
			if len(dirtyData.SpecialNotes) > 26 {
				e["special_notes"] = "over 26 characters"
			}
		}
	}

	// Process optional validation for `Grading Notes`.
	if dirtyData.GradingNotes != "" && len(dirtyData.GradingNotes) > 638 {
		e["grading_notes"] = "over 638 characters"
	}
	if dirtyData.Status == 0 {
		e["status"] = "missing choice"
	}
	if dirtyData.ServiceType == 0 {
		e["service_type"] = "missing choice"
	}
	if dirtyData.StoreID.IsZero() {
		e["store_id"] = "missing choice"
	}
	if dirtyData.Printing == 0 {
		e["printing"] = "missing choice"
	}
	if dirtyData.PrimaryLabelDetailsOther == "" && dirtyData.PrimaryLabelDetails == 1 {
		e["primary_label_details_other"] = "missing choice"
	}

	// Special case: Signatures only supported in certain cases.
	if len(dirtyData.Signatures) > 0 {
		if dirtyData.ServiceType != sub_s.ServiceTypeCPSCapsuleIndieMintGem && dirtyData.ServiceType != sub_s.ServiceTypeCPSCapsuleSignatureCollection && dirtyData.ServiceType != sub_s.ServiceTypePedigree {
			e["service_type"] = "cannot have signatures"
		}
	}
	if dirtyData.IsKeyIssue {
		if dirtyData.KeyIssue == 0 {
			e["key_issue"] = "missing choice"
		} else {
			// OTHER option.
			if dirtyData.KeyIssue == 1 && dirtyData.KeyIssueOther == "" {
				e["key_issue_other"] = "missing choice"
			}
			// Non-OTHER option.
			if dirtyData.KeyIssue != 1 {
				if dirtyData.KeyIssueDetail == "" {
					e["KeyIssueDetail"] = "missing choice"
				}
			}
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	d, err := h.unmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	d.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	submission, err := h.Controller.UpdateByID(ctx, d)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(submission, w)
}

func MarshalUpdateResponse(res *sub_s.ComicSubmission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
