package datastore

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SortOrderAscending  = 1
	SortOrderDescending = -1
)

type PinObjectPaginationListFilter struct {
	// Pagination related.
	Cursor    string
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	TenantID  primitive.ObjectID
	ProjectID primitive.ObjectID
}

// PinObjectPaginationListResult represents the paginated list results for
// the associate records.
type PinObjectPaginationListResult struct {
	Results     []*PinObject `json:"results"`
	NextCursor  string       `json:"next_cursor"`
	HasNextPage bool         `json:"has_next_page"`
}

// newPaginationFilter will create the mongodb filter to apply the cursor or
// or ignore it depending if a cursor was specified in the filter.
func (impl PinObjectStorerImpl) newPaginationFilter(f *PinObjectPaginationListFilter) (bson.M, error) {
	if len(f.Cursor) > 0 {
		// STEP 1: Decode the cursor which is encoded in a base64 format.
		decodedCursor, err := base64.RawStdEncoding.DecodeString(f.Cursor)
		if err != nil {
			return bson.M{}, fmt.Errorf("Failed to decode string: %v", err)
		}

		// STEP 2: Pick the specific cursor to build or else error.
		switch f.SortField {
		case "created", "modified_by_user_id":
			// STEP 3: Build for `time` field.
			return impl.newPaginationFilterBasedOnTime(f, string(decodedCursor))
		default:
			return nil, fmt.Errorf("unsupported sort field for `%v`, only supported fields are `created` and `modified_at`", f.SortField)
		}
	}
	return bson.M{}, nil
}

func (impl PinObjectStorerImpl) newPaginationFilterBasedOnTime(f *PinObjectPaginationListFilter, decodedCursor string) (bson.M, error) {
	// Extract our cursor into two parts which we need to use.
	arr := strings.Split(decodedCursor, "|")
	if len(arr) < 1 {
		return nil, fmt.Errorf("cursor is corrupted for the value `%v`", decodedCursor)
	}

	// The first part will contain the name we left off at. The second part will
	// be last ID we left off at.
	timeStr := arr[0]
	lastID, err := primitive.ObjectIDFromHex(arr[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to convert into mongodb object id: %v, from the decoded cursor of: %v", err, decodedCursor)
	}

	time, err := timekit.ParseJavaScriptTimeString(timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse javascript time: `%v`", err)
	}

	switch f.SortOrder {
	case SortOrderAscending:
		filter := bson.M{}
		filter["$or"] = []bson.M{
			bson.M{f.SortField: bson.M{"$gt": time}},
			bson.M{f.SortField: time, "_id": bson.M{"$gt": lastID}},
		}
		return filter, nil
	case SortOrderDescending:
		filter := bson.M{}
		filter["$or"] = []bson.M{
			bson.M{f.SortField: bson.M{"$lt": time}},
			bson.M{f.SortField: time, "_id": bson.M{"$lt": lastID}},
		}
		return filter, nil
	default:
		return nil, fmt.Errorf("unsupported sort order for `%v`, only supported values are `1` or `-1`", f.SortOrder)
	}
}

// newPaginatorOptions will generate the mongodb options which will support the
// paginator in ordering the data to work.
func (impl PinObjectStorerImpl) newPaginationOptions(f *PinObjectPaginationListFilter) (*options.FindOptions, error) {
	options := options.Find().SetLimit(f.PageSize)

	// DEVELOPERS NOTE:
	// We want to be able to return a list without sorting so we will need to
	// run the following code.
	if f.SortField != "" {
		options = options.
			SetSort(bson.D{
				{f.SortField, f.SortOrder},
				{"_id", f.SortOrder}, // Include _id in sorting for consistency
			})
	}

	return options, nil
}

// newPaginatorNextCursor will return the base64 encoded next cursor which works
// with our paginator.
func (impl PinObjectStorerImpl) newPaginatorNextCursor(f *PinObjectPaginationListFilter, results []*PinObject) (string, error) {
	var lastDatum *PinObject

	// Remove the extra document from the current page
	results = results[:len(results)]

	// Get the last document's _id as the next cursor
	lastDatum = results[len(results)-1]

	// Variable used to store the next cursor.
	var nextCursor string

	switch f.SortField {
	case "created":
		time := lastDatum.Created.UnixMilli()
		nextCursor = fmt.Sprintf("%v|%v", time, lastDatum.ID.Hex())
		break
	case "modified_at":
		time := lastDatum.ModifiedAt.UnixMilli()
		nextCursor = fmt.Sprintf("%v|%v", time, lastDatum.ID.Hex())
		break
	default:
		return "", fmt.Errorf("unsupported sort field in options for `%v`, only supported fields are `created` and `modified_at`", f.SortField)
	}

	// Encode to base64 without the `=` symbol that would corrupt when we
	// use the http url argument. Special thanks to:
	// https://www.golinuxcloud.com/golang-base64-encode/
	encoded := base64.RawStdEncoding.EncodeToString([]byte(nextCursor))

	return encoded, nil
}
