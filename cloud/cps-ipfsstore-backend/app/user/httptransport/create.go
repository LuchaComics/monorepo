package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	usr_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/controller"
	usr_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*usr_c.UserCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData usr_c.UserCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("user | UnmarshalCreateRequest | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *usr_c.UserCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.TenantID.IsZero() {
		e["tenant_id"] = "missing value"
	}
	if dirtyData.Role == 0 {
		e["role"] = "missing value"
	}
	if dirtyData.Status == 0 {
		e["status"] = "missing value"
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

	// The following logic will enforce shipping address input validation.
	if dirtyData.HasShippingAddress == true {
		if dirtyData.ShippingName == "" {
			e["shipping_name"] = "missing value"
		}
		if dirtyData.ShippingPhone == "" {
			e["shipping_phone"] = "missing value"
		}
		if dirtyData.ShippingCountry == "" {
			e["shipping_country"] = "missing value"
		}
		if dirtyData.ShippingRegion == "" {
			e["shipping_region"] = "missing value"
		}
		if dirtyData.ShippingCity == "" {
			e["shipping_city"] = "missing value"
		}
		if dirtyData.ShippingPostalCode == "" {
			e["shipping_postal_code"] = "missing value"
		}
		if dirtyData.ShippingAddressLine1 == "" {
			e["shipping_address_line1"] = "missing value"
		}
	}
	if dirtyData.Role == usr_s.UserRoleCustomer {
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

	user, err := h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(user, w)
}

func MarshalCreateResponse(res *usr_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
