package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
)

// NFT structured used to store the metadata for the NFT using the
// `OpenSea` standard via https://docs.opensea.io/docs/metadata-standards.
type NFT struct {
	TenantID              primitive.ObjectID      `bson:"tenant_id" json:"tenant_id"`
	TenantName            string                  `bson:"tenant_name" json:"tenant_name"`
	TenantTimezone        string                  `bson:"tenant_timezone" json:"tenant_timezone"`
	ID                    primitive.ObjectID      `bson:"_id" json:"id"`
	Status                int8                    `bson:"status" json:"status"`
	CreatedAt             time.Time               `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedFromIPAddress  string                  `bson:"created_from_ip_address" json:"created_from_ip_address,omitempty"`
	ModifiedAt            time.Time               `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedFromIPAddress string                  `bson:"modified_from_ip_address" json:"modified_from_ip_address,omitempty"`
	CollectionID          primitive.ObjectID      `bson:"collection_id" json:"collection_id"`
	CollectionName        string                  `bson:"collection_name" json:"collection_name"`
	TokenID               uint64                  `bson:"token_id" json:"token_id"`
	ImageID               primitive.ObjectID      `bson:"image_id" json:"image_id"`
	Image                 string                  `bson:"image" json:"image"` // This is the URL to the image of the item. Can be just about any type of image (including SVGs, which will be cached into PNGs by OpenSea), IPFS or Arweave URLs or paths. We recommend using a minimum 3000 x 3000 image.
	ImageFilename         string                  `bson:"image_filename" json:"image_filename"`
	ImageCID              string                  `bson:"image_cid" json:"image_cid"`
	AnimationID           primitive.ObjectID      `bson:"animation_id" json:"animation_id"`
	AnimationFilename     string                  `bson:"animation_filename" json:"animation_filename"`
	AnimationURL          string                  `bson:"animation_url" json:"animation_url"` // A URL to a multi-media attachment for the item. The file extensions GLTF, GLB, WEBM, MP4, M4V, OGV, and OGG are supported, along with the audio-only extensions MP3, WAV, and OGA. Animation_url also supports HTML pages, allowing you to build rich experiences and interactive NFTs using JavaScript canvas, WebGL, and more. Scripts and relative paths within the HTML page are now supported. However, access to browser extensions is not supported.
	AnimationCID          string                  `bson:"animation_cid" json:"animation_cid"`
	ExternalURL           string                  `bson:"external_url" json:"external_url"`         // This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	Description           string                  `bson:"description" json:"description"`           // A human-readable description of the item. Markdown is supported.
	Name                  string                  `bson:"name" json:"name"`                         // Name of the item.
	Attributes            []*NFTMetadataAttribute `bson:"attributes" json:"attributes"`             // These are the attributes for the item, which will show up on the OpenSea page for the item. (see below)
	BackgroundColor       string                  `bson:"background_color" json:"background_color"` // Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	YoutubeURL            string                  `bson:"youtube_url" json:"youtube_url"`           // A URL to a YouTube video (only used if animation_url is not provided).
	FileCID               string                  `bson:"file_cid" json:"file_cid"`
	FileIPNSPath          string                  `bson:"file_ipns_path" json:"file_ipns_path"` // The path of this metadata file in the IPFS network utilizing IPNS.
	MintedToAddress       string                  `bson:"minted_to_address" json:"minted_to_address"`
}

type NFTMetadataFile struct {
	Image           string                  `bson:"image" json:"image"`
	ExternalURL     string                  `bson:"external_url" json:"external_url"`
	Description     string                  `bson:"description" json:"description"`
	Name            string                  `bson:"name" json:"name"`
	Attributes      []*NFTMetadataAttribute `bson:"attributes" json:"attributes"`
	BackgroundColor string                  `bson:"background_color" json:"background_color"`
	AnimationURL    string                  `bson:"animation_url" json:"animation_url"`
	YoutubeURL      string                  `bson:"youtube_url" json:"youtube_url"`
}

type NFTMetadataAttribute struct {
	DisplayType string `bson:"display_type" json:"display_type"`
	TraitType   string `bson:"trait_type" json:"trait_type"`
	Value       string `bson:"value" json:"value"`
}

type NFTListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	TenantID        primitive.ObjectID
	UserID          primitive.ObjectID
	ExcludeArchived bool
	SearchText      string
}

type NFTListResult struct {
	Results     []*NFT             `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// NFTStorer Interface for tenant.
type NFTStorer interface {
	Create(ctx context.Context, m *NFT) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*NFT, error)
	GetByTokenID(ctx context.Context, tokenID uint64) (*NFT, error)
	UpdateByID(ctx context.Context, m *NFT) error
	ListByFilter(ctx context.Context, m *NFTPaginationListFilter) (*NFTPaginationListResult, error)
	ListByNFTCollectionID(ctx context.Context, nftCollectionID primitive.ObjectID) (*NFTPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *NFTPaginationListFilter) ([]*NFTAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type NFTAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type NFTStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) NFTStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("nfts")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &NFTStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
