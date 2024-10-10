package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/logger"
	handler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
)

// Command line argument flags
var (
	flagProofOfAuthorityAccountAddress string
	flagProofOfAuthorityWalletPassword string
	flagMintToAddress                  string
	flagMintImage                      string
	flagMintExternalURL                string
	flagMintDescription                string
	flagMintName                       string
	flagMintBackgroundColor            string
	flagMintAnimationURL               string
	flagMintYoutubeURL                 string
	flagMintMetadataURI                string

// TODO: IMPL.
// Attributes []*NFTMetadataAttribute `bson:"attributes" json:"attributes"`
//
//	type NFTMetadataAttribute struct {
//		DisplayType string `bson:"display_type" json:"display_type"`
//		TraitType   string `bson:"trait_type" json:"trait_type"`
//		Value       string `bson:"value" json:"value"`
//	}
)

func httpJsonApiMintNFTCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mint",
		Short: "(PoA & Authority only) Creates a new non-fungible token in our blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			doMintNFT()
		},
	}

	cmd.Flags().StringVar(&flagProofOfAuthorityAccountAddress, "poa-address", "", "(Required for `PoA` consensus protocol) The address of the authority's account")
	cmd.MarkFlagRequired("poa-address")
	cmd.Flags().StringVar(&flagProofOfAuthorityWalletPassword, "poa-password", "", "(Required for `PoA` consensus protocol) The password in the authority's wallet")
	cmd.MarkFlagRequired("poa-password")
	cmd.Flags().StringVar(&flagMintToAddress, "recipient-address", "", "The address of the account whom will receive this NFT")
	cmd.MarkFlagRequired("recipient-address")

	// Fields for inputting the NFT
	cmd.Flags().StringVar(&flagMintImage, "image", "", "This is the URL to the image of the item. Can be just about any type of image (including SVGs, which will be cached into PNGs by OpenSea), IPFS or Arweave URLs or paths. We recommend using a minimum 3000 x 3000 image.")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringVar(&flagMintExternalURL, "external-url", "", "This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.")
	cmd.MarkFlagRequired("external-url")
	cmd.Flags().StringVar(&flagMintDescription, "description", "", "A human-readable description of the item. Markdown is supported.")
	cmd.MarkFlagRequired("description")
	cmd.Flags().StringVar(&flagMintName, "name", "", "Name of the item.")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flagMintBackgroundColor, "background-color", "", "Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.")
	cmd.MarkFlagRequired("background-color")
	cmd.Flags().StringVar(&flagMintAnimationURL, "animation-url", "", "Animation_url also supports HTML pages, allowing you to build rich experiences and interactive NFTs using JavaScript canvas, WebGL, and more. Scripts and relative paths within the HTML page are now supported. However, access to browser extensions is not supported.")
	cmd.MarkFlagRequired("animation-url")
	cmd.Flags().StringVar(&flagMintYoutubeURL, "youtube-url", "", "(Optional) A URL to a YouTube video (only used if animation_url is not provided).")
	cmd.Flags().StringVar(&flagMintMetadataURI, "metadata-uri", "", "The location of this tokens metadata file.")
	cmd.MarkFlagRequired("metadata-uri")

	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "(Optional) The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "(Optional) The HTTP JSON API server's ip-address")

	return cmd
}

func doMintNFT() {
	//
	// STEP 1:
	// Get our project dependencies in order.
	//
	logger := logger.NewProvider()

	//
	// STEP 2:
	// Create our request payload.
	//

	httpEndpoint := fmt.Sprintf("http://%s:%d%s", flagListenHTTPIP, flagListenHTTPPort, nftsURL)

	metadata := handler.MintNFTRequestIDO{
		ProofOfAuthorityAccountAddress: flagProofOfAuthorityAccountAddress,
		ProofOfAuthorityWalletPassword: flagProofOfAuthorityWalletPassword,
		To:                             flagMintToAddress,
		Image:                          flagMintImage,
		ExternalURL:                    flagMintExternalURL,
		Description:                    flagMintDescription,
		Name:                           flagMintName,
		Attributes: []*handler.NFTMetadataAttributeRequestIDO{
			&handler.NFTMetadataAttributeRequestIDO{ //TODO: IMPL.
				DisplayType: "string",
				TraitType:   "color",
				Value:       "red",
			},
		},
		BackgroundColor: flagMintBackgroundColor,
		AnimationURL:    flagMintAnimationURL,
		YoutubeURL:      flagMintYoutubeURL,
		MetadataURI:     flagMintMetadataURI,
	}
	logger.Debug("Creating new NFT in blockchain",
		slog.Any("node-url", httpEndpoint),
		slog.Any("image", flagMintImage),
		slog.Any("external_url", flagMintExternalURL),
		slog.Any("description", flagMintDescription),
		slog.Any("name", flagMintName),
		slog.Any("background_color", flagMintBackgroundColor),
		slog.Any("animation_url", flagMintAnimationURL),
		slog.Any("youtube_url", flagMintYoutubeURL),
		slog.Any("metadata_uri", flagMintMetadataURI),
	)

	//
	// STEP 3
	// Convert request to binary and submit to running HTTP JSON API.
	//

	reqBytes, err := json.Marshal(&metadata)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	if reqBytes == nil {
		log.Fatal("nothing marshalled")
	}

	//
	// STEP 4
	// Setup our HTTP client for sending
	//

	r, err := http.NewRequest("POST", httpEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("failed to setup post request: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "POST"))

	//
	// STEP 5
	// Wait for the submission to finish sending and then get the resposne.
	//

	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		//
		// STEP 6
		// If successful then return our results.
		//
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

		// Try to decode the response as a string first
		var jsonStr string
		err := json.NewDecoder(teeReader).Decode(&jsonStr)
		if err != nil {
			logger.Error("decoding string error",
				slog.Any("err", err),
				slog.String("json", rawJSON.String()),
			)
			return
		}

		// Now try to decode the string into a map
		err = json.Unmarshal([]byte(jsonStr), &e)
		if err != nil {
			logger.Error("decoding map error",
				slog.Any("err", err),
				slog.String("json", jsonStr),
			)
			return
		}

		logger.Debug("Parsed error response",
			slog.Any("errors", e),
		)
		return
	}
	logger.Debug("NFT mint request submitted successful")
}
