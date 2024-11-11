package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	_ "time/tzdata"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/tenant/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*sub_s.Tenant, error) {
	// Initialize our array which will tenant all the results from the remote server.
	var requestData sub_s.Tenant

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *sub_s.Tenant) error {
	e := make(map[string]string)

	// if dirtyData.ServiceType == 0 {
	// 	e["service_type"] = "missing value"
	// }
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.WebsiteURL == "" {
		e["website_url"] = "missing value"
	}
	if dirtyData.EstimatedSubmissionsPerMonth == 0 {
		e["estimated_submissions_per_month"] = "missing value"
	}
	if dirtyData.HasOtherGradingService == 0 {
		e["has_other_grading_service"] = "missing value"
	} else {
		// if dirtyData.OtherGradingServiceName == "" {
		// 	e["other_grading_service_name"] = "missing value"
		// }
	}
	if dirtyData.RequestWelcomePackage == 0 {
		e["request_welcome_package"] = "missing value"
	}
	if dirtyData.HowLongTenantOperating == 0 {
		e["how_long_tenant_operating"] = "missing value"
	}
	// if dirtyData.GradingComicsExperience == "" {
	// 	e["grading_comics_experience"] = "missing value"
	// }
	if dirtyData.RetailPartnershipReason == "" {
		e["retail_partnership_reason"] = "missing value"
	}
	if dirtyData.CPS_IPFSSTOREPartnershipReason == "" {
		e["cps-ipfsstore_partnership_reason"] = "missing value"
	}
	if dirtyData.Level == 0 {
		e["level"] = "missing value"
	}
	if dirtyData.Timezone == "" {
		e["timezone"] = "missing value"
	} else {
		// Confirm the timezone is one that exists.
		location, err := time.LoadLocation(dirtyData.Timezone)
		if err != nil || location == nil {
			e["timezone"] = "unsupported value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	log.Println("data:", data)

	data, err = h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(data, w)
}

func MarshalCreateResponse(res *sub_s.Tenant, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
