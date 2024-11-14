package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

// FileBaseRepoConfigurationProvider is an interface for configuration providers
// that provide all needed settings to connect to an FileBase node either remote
// or a local FileBase node.
type FileBaseRepoConfigurationProvider interface {
	GetApiVersion() string
	GetAccessKeyId() string
	GetSecretAccessKey() string
	GetEndpoint() string
	GetRegion() string
	GetS3ForcePathStyle() string
}

// FileBaseRepoConfigurationProviderImpl is a struct that implements
// FileBaseRepoConfigurationProvider for storing FileBase connection details.
type FileBaseRepoConfigurationProviderImpl struct {
	apiVersion       string
	accessKeyId      string
	secretAccessKey  string
	endpoint         string
	region           string
	s3ForcePathStyle string
}

// NewFileBaseRepoConfigurationProvider constructs a new configuration provider
// for FileBase connections.
func NewFileBaseRepoConfigurationProvider(
	apiVersion string,
	accessKeyId string,
	secretAccessKey string,
	endpoint string,
	region string,
	s3ForcePathStyle string,
) FileBaseRepoConfigurationProvider {
	// Defensive code: Enforce parameters.
	if apiVersion == "" {
		log.Fatal("Missing `apiVersion` parameter.")
	}
	if accessKeyId == "" {
		log.Fatal("Missing `accessKeyId` parameter.")
	}
	if secretAccessKey == "" {
		log.Fatal("Missing `secretAccessKey` parameter.")
	}
	if endpoint == "" {
		log.Fatal("Missing `endpoint` parameter.")
	}
	if endpoint == "" {
		log.Fatal("Missing `endpoint` parameter.")
	}
	if s3ForcePathStyle == "" {
		log.Fatal("Missing `s3ForcePathStyle` parameter.")
	}
	return &FileBaseRepoConfigurationProviderImpl{
		apiVersion:       apiVersion,
		accessKeyId:      accessKeyId,
		secretAccessKey:  secretAccessKey,
		endpoint:         endpoint,
		region:           region,
		s3ForcePathStyle: s3ForcePathStyle,
	}
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetApiVersion() string {
	return impl.apiVersion
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetAccessKeyId() string {
	return impl.accessKeyId
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetSecretAccessKey() string {
	return impl.secretAccessKey
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetEndpoint() string {
	return impl.endpoint
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetRegion() string {
	return impl.region
}

func (impl *FileBaseRepoConfigurationProviderImpl) GetS3ForcePathStyle() string {
	return impl.s3ForcePathStyle
}

type FileBaseRepo struct {
	config        FileBaseRepoConfigurationProvider // Holds FileBase connection configuration
	logger        *slog.Logger
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	BucketName    string
}

const maxRetries = 10
const retryDelay = 15 * time.Second

// connectToS3 connects to S3 and returns the S3 client.
func connectToS3(cfg FileBaseRepoConfigurationProvider, logger *slog.Logger) (*s3.Client, error) {
	// Step 1: Initialize the custom `endpoint` we will connect to.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: cfg.GetEndpoint(),
		}, nil
	})

	// Step 2: Configure.
	sdkConfig, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion(cfg.GetRegion()),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.GetAccessKeyId(), cfg.GetSecretAccessKey(), "")),
	)
	if err != nil {
		return nil, err
	}

	// Step 3: Load up S3 instance.
	s3Client := s3.NewFromConfig(sdkConfig)

	// For debugging purposes only.
	logger.Debug("Connected to FileBase.")

	return s3Client, nil
}

// NewFileBaseRepo returns a new FileBaseNode instance
func NewFileBaseRepo(cfg FileBaseRepoConfigurationProvider, logger *slog.Logger) domain.FileBaseRepository {
	logger.Debug("FileBase initializing...")

	// Retry logic
	var err error
	var s3Client *s3.Client
	for i := 1; i <= maxRetries; i++ {
		s3Client, err = connectToS3(cfg, logger)
		if err == nil {
			break
		}

		logger.Warn(fmt.Sprintf("Failed to connect to FileBase (attempt %d/%d): %v", i, maxRetries, err))
		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Fatal("Failed to connect to FileBase after retries")
	}

	fileBaseRepo := &FileBaseRepo{
		config:        cfg,
		logger:        logger,
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
		BucketName:    "comiccoin",
	}

	// STEP 4: Connect to the s3 bucket instance and confirm that bucket exists.
	doesExist, err := fileBaseRepo.BucketExists(context.TODO(), "comiccoin")
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}
	if !doesExist {
		log.Fatal("bucket name does not exist") // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	// For debugging purposes only.
	logger.Debug("FileBase ready for usage.")

	return fileBaseRepo
}

func (s *FileBaseRepo) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	// Note: https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html#actions

	_, err := s.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	}

	return exists, err
}

// // ID returns the FileBase node's identity information
//
//	func (r *FileBaseRepo) ID() (peer.ID, error) {
//		keyAPI := r.api.Key()
//		if keyAPI == nil {
//			return "", fmt.Errorf("Failed getting key: %v", "does not exist")
//		}
//		selfKey, err := keyAPI.Self(context.Background())
//		if err != nil {
//			return "", fmt.Errorf("Failed getting self: %v", err)
//		}
//		if selfKey == nil {
//			return "", fmt.Errorf("Failed getting self: %v", "does not exist")
//		}
//		return selfKey.ID(), nil
//	}
func (r *FileBaseRepo) AddViaFilePath(fullFilePath string, shouldPin bool) (string, error) {
	return "", nil
}

// 	unixfs := r.api.Unixfs()
// 	if unixfs == nil {
// 		return "", fmt.Errorf("Failed getting unix fs: %v", "does not exist")
// 	}
//
// 	// Open the file
// 	file, err := os.Open(fullFilePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()
//
// 	// Get the file stat
// 	stat, err := file.Stat()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Create a reader file node
// 	node, err := files.NewReaderPathFile(fullFilePath, file, stat)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// We want to use the newest `CidVersion` in our update.
// 	opts := func(settings *options.UnixfsAddSettings) error {
// 		settings.CidVersion = 1
// 		settings.Pin = shouldPin
// 		return nil
// 	}
//
// 	pathRes, err := unixfs.Add(context.Background(), node, opts)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return strings.Replace(pathRes.String(), "/ipfs/", "", -1), nil
// }
//
// func (r *FileBaseRepo) AddViaFileContent(fileContent []byte, shouldPin bool) (string, error) {
// 	unixfs := r.api.Unixfs()
// 	if unixfs == nil {
// 		return "", fmt.Errorf("Failed getting unix fs: %v", "does not exist")
// 	}
//
// 	file := bytes.NewReader(fileContent)
//
// 	// Create a reader file node
// 	node := files.NewReaderFile(file)
//
// 	// We want to use the newest `CidVersion` in our update.
// 	opts := func(settings *options.UnixfsAddSettings) error {
// 		settings.CidVersion = 1
// 		settings.Pin = shouldPin
// 		return nil
// 	}
//
// 	pathRes, err := unixfs.Add(context.Background(), node, opts)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return strings.Replace(pathRes.String(), "/ipfs/", "", -1), nil
// }
//
// func (r *FileBaseRepo) AddViaFile(file *os.File, shouldPin bool) (string, error) {
// 	unixfs := r.api.Unixfs()
// 	if unixfs == nil {
// 		return "", fmt.Errorf("Failed getting unix fs: %v", "does not exist")
// 	}
//
// 	// Get the file stat
// 	stat, err := file.Stat()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Create a reader file node
// 	node := files.NewReaderStatFile(file, stat)
//
// 	// We want to use the newest `CidVersion` in our update.
// 	opts := func(settings *options.UnixfsAddSettings) error {
// 		settings.CidVersion = 1
// 		settings.Pin = shouldPin
// 		return nil
// 	}
//
// 	pathRes, err := unixfs.Add(context.Background(), node, opts)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return strings.Replace(pathRes.String(), "/ipfs/", "", -1), nil
// }
//
// func (r *FileBaseRepo) AddViaReaderFile(node files.File, shouldPin bool) (string, error) {
// 	unixfs := r.api.Unixfs()
// 	if unixfs == nil {
// 		return "", fmt.Errorf("Failed getting unix fs: %v", "does not exist")
// 	}
//
// 	// We want to use the newest `CidVersion` in our update.
// 	opts := func(settings *options.UnixfsAddSettings) error {
// 		settings.CidVersion = 1
// 		settings.Pin = shouldPin
// 		return nil
// 	}
//
// 	pathRes, err := unixfs.Add(context.Background(), node, opts)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return strings.Replace(pathRes.String(), "/ipfs/", "", -1), nil
// }
//
// func (impl *FileBaseRepo) Pin(cidString string) error {
// 	impl.logger.Debug("pinning content to FileBase", slog.String("cid", cidString))
//
// 	cid, err := cid.Decode(cidString)
// 	if err != nil {
// 		impl.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
// 		return fmt.Errorf("failed to decode CID: %v", err)
// 	}
//
// 	// Convert the CID to a path.Path
// 	ipfsPath := path.FromCid(cid)
//
// 	// Attempt to pin the content to the FileBase node using the CID
// 	if err := impl.api.Pin().Add(context.Background(), ipfsPath); err != nil {
// 		impl.logger.Error("failed to pin content to FileBase", slog.String("cid", cidString), slog.Any("error", err))
// 		return fmt.Errorf("failed to pin content to FileBase: %v", err)
// 	}
// 	return nil
// }
//
// func (r *FileBaseRepo) PinAddViaFilePath(fullFilePath string) (string, error) {
// 	fileCID, err := r.AddViaFilePath(fullFilePath, false)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if err := r.Pin(fileCID); err != nil {
// 		return "", err
// 	}
//
// 	return fileCID, nil
// }
//
// // Cat retrieves the contents of a file from FileBase
// func (s *FileBaseRepo) Get(ctx context.Context, cidString string) ([]byte, string, error) {
// 	s.logger.Debug("fetching content from FileBase", slog.String("cid", cidString))
//
// 	cid, err := cid.Decode(cidString)
// 	if err != nil {
// 		s.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
// 		return nil, "", fmt.Errorf("failed to decode CID: %v", err)
// 	}
//
// 	// Convert the CID to a path.Path
// 	ipfsPath := path.FromCid(cid)
//
// 	// Add a timeout to prevent hanging requests.
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	// Attempt to get the file from FileBase using the path
// 	fileNode, err := s.api.Unixfs().Get(ctx, ipfsPath)
// 	if err != nil {
// 		s.logger.Debug("Failed fetching from remote FileBase...",
// 			slog.String("cid", cidString),
// 			slog.Any("error", err))
// 		return s.getViaHTTPPublicGateway(ctx, cidString)
//
// 	}
//
// 	// Convert the file node to a reader
// 	fileReader := files.ToFile(fileNode)
// 	if fileReader == nil {
// 		s.logger.Error("failed to convert FileBase node to file reader", slog.String("cid", cidString))
// 		return nil, "", fmt.Errorf("failed to convert FileBase node to file reader")
// 	}
//
// 	// Read the content from the file reader
// 	content, err := io.ReadAll(fileReader)
// 	if err != nil {
// 		s.logger.Error("failed to read content from FileBase", slog.String("cid", cidString), slog.Any("error", err))
// 		return nil, "", fmt.Errorf("failed to read content from FileBase: %v", err)
// 	}
//
// 	return content, http.DetectContentType(content), nil
// }
//
// func (s *FileBaseRepo) getViaHTTPPublicGateway(ctx context.Context, cidString string) ([]byte, string, error) {
// 	uri := fmt.Sprintf("%v/ipfs/%v", s.config.GetPublicFileBaseGatewayAddress(), cidString)
//
// 	s.logger.Debug("Fetching from public FileBase gateway... Please wait...",
// 		slog.String("cid", cidString))
//
// 	resp, err := http.Get(uri)
// 	if err != nil {
// 		s.logger.Error("Failed fetching metadata uri via http.",
// 			slog.Any("error", err))
// 		return nil, "", err
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		err := fmt.Errorf("status code error: %d", resp.StatusCode)
// 		s.logger.Error("Status code error",
// 			slog.Any("error", err))
// 		return nil, "", err
// 	}
//
// 	// Get the content type from the response header
// 	contentType := resp.Header.Get("Content-Type")
// 	if contentType == "" {
// 		err := fmt.Errorf("Content type not specified in response header")
// 		s.logger.Error("Content-Type error",
// 			slog.Any("error", err))
// 		return nil, "", err
// 	}
//
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		s.logger.Error("Failed read all.",
// 			slog.Any("error", err))
// 		return nil, "", err
// 	}
//
// 	s.logger.Debug("Successfully fetching from public FileBase gateway.",
// 		slog.String("cid", cidString))
//
// 	return body, contentType, nil
// }
