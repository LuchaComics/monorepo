package domain

import "github.com/ethereum/go-ethereum/common"

// NFTTransaction represents a transaction specifically related to an NFT.
type NFTTransaction struct {
	ChainID     uint16          `json:"chain_id"`     // Ethereum: The chain id that is listed in the genesis file.
	TokenID     uint64          `json:"token_id"`     // Unique identifier for the NFT.
	From        *common.Address `json:"from"`         // Account sending the NFT.
	To          *common.Address `json:"to"`           // Account receiving the NFT.
	Metadata    *NFTMetadata    `json:"metadata"`     // Metadata of the NFT.
	MetedataURI string          `json:"metadata_uri"` // URI pointing to NFT metadata file.
	TimeStamp   uint64          `json:"timestamp"`    // Timestamp of the NFT transaction.
}

// NFTMetadata structured used to store the metadata for the NFT using the
// `OpenSea` standard via https://docs.opensea.io/docs/metadata-standards.
type NFTMetadata struct {
	// This is the URL to the image of the item. Can be just about any type of image (including SVGs, which will be cached into PNGs by OpenSea), IPFS or Arweave URLs or paths. We recommend using a minimum 3000 x 3000 image.
	Image string `bson:"image" json:"image"`

	// This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	ExternalURL string `bson:"external_url" json:"external_url"`

	// A human-readable description of the item. Markdown is supported.
	Description string `bson:"description" json:"description"`

	// Name of the item.
	Name string `bson:"name" json:"name"`

	// These are the attributes for the item, which will show up on the OpenSea page for the item. (see below
	Attributes []*NFTMetadataAttribute `bson:"attributes" json:"attributes"`

	// Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	BackgroundColor string `bson:"background_color" json:"background_color"`

	// A URL to a multi-media attachment for the item. The file extensions GLTF, GLB, WEBM, MP4, M4V, OGV, and OGG are supported, along with the audio-only extensions MP3, WAV, and OGA.
	//
	// Animation_url also supports HTML pages, allowing you to build rich experiences and interactive NFTs using JavaScript canvas, WebGL, and more. Scripts and relative paths within the HTML page are now supported. However, access to browser extensions is not supported.
	AnimationURL string `bson:"animation_url" json:"animation_url"`

	// A URL to a YouTube video (only used if animation_url is not provided).
	YoutubeURL string `bson:"youtube_url" json:"youtube_url"`
}

type NFTMetadataAttribute struct {
	DisplayType string `bson:"display_type" json:"display_type"`
	TraitType   string `bson:"trait_type" json:"trait_type"`
	Value       string `bson:"value" json:"value"`
}
