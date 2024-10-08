package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	usr_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*usr_s.User, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData usr_s.User

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateUpdateRequest(dirtyData *usr_s.User) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
	if dirtyData.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if dirtyData.LastName == "" {
		e["last_name"] = "missing value"
	}
	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Phone == "" {
		e["phone"] = "missing value"
	}
	// if dirtyData.Country == "" {
	// 	e["country"] = "missing value"
	// }
	// if dirtyData.Region == "" {
	// 	e["region"] = "missing value"
	// }
	// if dirtyData.City == "" {
	// 	e["city"] = "missing value"
	// }
	// if dirtyData.PostalCode == "" {
	// 	e["postal_code"] = "missing value"
	// }
	// if dirtyData.AddressLine1 == "" {
	// 	e["address_line1"] = "missing value"
	// }
	if dirtyData.HowDidYouHearAboutUs > 7 || dirtyData.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if dirtyData.HowDidYouHearAboutUs == 1 && dirtyData.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}

	if dirtyData.HowLongCollectingComicBooksForGrading == 0 {
		e["how_long_collecting_comic_books_for_grading"] = "missing value"
	}
	if dirtyData.HasPreviouslySubmittedComicBookForGrading == 0 {
		e["has_previously_submitted_comic_book_for_grading"] = "missing value"
	}
	if dirtyData.HasOwnedGradedComicBooks == 0 {
		e["has_owned_graded_comic_books"] = "missing value"
	}
	if dirtyData.HasRegularComicBookShop == 0 {
		e["has_regular_comic_book_shop"] = "missing value"
	}
	if dirtyData.HasPreviouslyPurchasedFromAuctionSite == 0 {
		e["has_previously_purchased_from_auction_site"] = "missing value"
	}
	if dirtyData.HasPreviouslyPurchasedFromFacebookMarketplace == 0 {
		e["has_previously_purchased_from_facebook_marketplace"] = "missing value"
	}
	if dirtyData.HasRegularlyAttendedComicConsOrCollectibleShows == 0 {
		e["has_regularly_attended_comic_cons_or_collectible_shows"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	data.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	customer, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(customer, w)
}

func MarshalUpdateResponse(res *usr_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
