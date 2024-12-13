package s3 // Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	cloudinterface "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud"
	c "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type s3Adapter struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	Logger        *slog.Logger
	BucketName    string
}

const maxRetries = 10
const retryDelay = 15 * time.Second

// connectToS3 connects to S3 and returns the S3 client.
func connectToS3(appConf *c.Configuration, logger *slog.Logger) (*s3.Client, error) {
	// Step 1: Initialize the custom `endpoint` we will connect to.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: appConf.AWS.Endpoint,
		}, nil
	})

	// Step 2: Configure.
	sdkConfig, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion(appConf.AWS.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appConf.AWS.AccessKey, appConf.AWS.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	// Step 3: Load up S3 instance.
	s3Client := s3.NewFromConfig(sdkConfig)

	// For debugging purposes only.
	logger.Debug("S3 connected to remote service")

	return s3Client, nil
}

// NewCloudStorage connects to a specific S3 bucket instance and returns a connected
// instance structure.
func NewCloudStorage(appConf *c.Configuration, logger *slog.Logger) cloudinterface.CloudStorage {
	// DEVELOPERS NOTE:
	// How can I use the AWS SDK v2 for Go with DigitalOcean Spaces? via https://stackoverflow.com/a/74284205
	logger.Debug("s3 initializing...")

	// Retry logic
	var err error
	var s3Client *s3.Client
	for i := 1; i <= maxRetries; i++ {
		s3Client, err = connectToS3(appConf, logger)
		if err == nil {
			break
		}

		logger.Warn(fmt.Sprintf("Failed to connect to S3 (attempt %d/%d): %v", i, maxRetries, err))
		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Fatal("Failed to connect to S3 after retries")
	}

	// Create our storage handler.
	s3Storage := &s3Adapter{
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
		Logger:        logger,
		BucketName:    appConf.AWS.BucketName,
	}

	// STEP 4: Connect to the s3 bucket instance and confirm that bucket exists.
	doesExist, err := s3Storage.BucketExists(context.TODO(), appConf.AWS.BucketName)
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}
	if !doesExist {
		log.Fatal("bucket name does not exist") // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	// For debugging purposes only.
	logger.Debug("s3 initialized")

	// Return our s3 storage handler.
	return s3Storage
}

func (s *s3Adapter) UploadContent(ctx context.Context, objectKey string, content []byte) error {
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Adapter) UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error {
	// Create the S3 upload input parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	}

	// Perform the file upload to S3
	_, err := s.S3Client.PutObject(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Adapter) UploadContentFromBytes(ctx context.Context, objectKey string, content []byte) error {
	// Create the S3 upload input parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(content), // Convert content []byte to io.Reader using bytes.NewReader
	}

	// Perform the file upload to S3
	_, err := s.S3Client.PutObject(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Adapter) BucketExists(ctx context.Context, bucketName string) (bool, error) {
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

func (s *s3Adapter) GetDownloadablePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	// DEVELOPERS NOTE:
	// AWS S3 Bucket — presigned URL APIs with Go (2022) via https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57

	presignedUrl, err := s.PresignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket:                     aws.String(s.BucketName),
			Key:                        aws.String(key),
			ResponseContentDisposition: aws.String("attachment"), // This field allows the file to download it directly from your browser
		},
		s3.WithPresignExpires(duration))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}

func (s *s3Adapter) GetPresignedURL(ctx context.Context, objectKey string, duration time.Duration) (string, error) {
	// DEVELOPERS NOTE:
	// AWS S3 Bucket — presigned URL APIs with Go (2022) via https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57

	presignedUrl, err := s.PresignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(duration))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}

func (s *s3Adapter) GetContentByKey(ctx context.Context, objectKey string) ([]byte, error) {
	// Create the S3 GetObject input parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}

	// Retrieve the object from S3
	resp, err := s.S3Client.GetObject(ctx, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Read the content from the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *s3Adapter) GetMultipartFileByKey(ctx context.Context, objectKey string) (multipart.File, error) {
	// Create the S3 GetObject input parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}

	// Retrieve the object from S3
	resp, err := s.S3Client.GetObject(ctx, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Read the content into a buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, err
	}

	// Create a multipart.File from the buffer
	file := &multipartFile{
		Reader: bytes.NewReader(buf.Bytes()),
	}

	return file, nil
}

// Implement a custom type that satisfies the multipart.File interface
type multipartFile struct {
	*bytes.Reader
}

func (f *multipartFile) Close() error {
	// For this custom implementation, there's nothing to close
	return nil
}

func (s *s3Adapter) DeleteByKeys(ctx context.Context, objectKeys []string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s.S3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.BucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", s.BucketName, err)
	}
	return err
}
