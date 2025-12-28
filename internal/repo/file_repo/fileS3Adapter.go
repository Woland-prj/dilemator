package file_repo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/file_errors"
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
)

var _ dilemma_service.FileRepositoryPort = (*FileS3RepositoryAdapter)(nil)

// FileS3RepositoryAdapter — FileRepositoryPort implementation for AWS S3 / Selectel / Compatible storage.
type FileS3RepositoryAdapter struct {
	client       *s3.Client
	presigner    *s3.PresignClient
	bucketDomain string
	bucketName   string
	endpoint     string
	region       string
	presignHours int
}

// NewFileS3Repository creates a new S3 repository.
func NewFileS3Repository(
	accessKey, secretKey, bucketName, region, endpoint, bucketDomain string,
	presignHours int,
) (*FileS3RepositoryAdapter, error) {
	const op = "repo - s3 - NewFileS3Repository"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	presigner := s3.NewPresignClient(client)

	return &FileS3RepositoryAdapter{
		client:       client,
		presigner:    presigner,
		bucketName:   bucketName,
		endpoint:     endpoint,
		bucketDomain: bucketDomain,
		region:       region,
		presignHours: presignHours,
	}, nil
}

func (r *FileS3RepositoryAdapter) Save(ctx context.Context, file []byte, contentType string) (string, error) {
	const op = "repo - s3 - Save"

	key := uuid.New().String()
	contentLength := int64(len(file))

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &r.bucketName,
		Key:           &key,
		Body:          bytes.NewReader(file),
		ContentType:   &contentType,
		ContentLength: &contentLength,
	})
	if err != nil {
		return "", berrors.Wrap(op, "Failed to uploadfile: ", err)
	}

	return key, nil
}

func (r *FileS3RepositoryAdapter) DeleteByKey(ctx context.Context, key string) error {
	const op = "repo - s3 - DeleteByKey"

	if err := r.validateExists(ctx, key); err != nil {
		return berrors.FromErr(op, err)
	}

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	})
	if err != nil {
		return berrors.FromErr(op, file_errors.ErrRemovingFailed)
	}

	return nil
}

func (r *FileS3RepositoryAdapter) GetDownloadLink(ctx context.Context, key string) (string, error) {
	const op = "repo - s3 - GetDownloadLink"
	if err := r.validateExists(ctx, key); err != nil {
		return "", berrors.FromErr(op, err)
	}

	// Format: https://bucket_domain/key
	u := fmt.Sprintf("%s/%s", r.bucketDomain, url.PathEscape(key))

	return u, nil
}

func (r *FileS3RepositoryAdapter) GetAuthorizedDownloadLink(ctx context.Context, key string) (string, error) {
	const op = "repo - s3 - GetAuthorizedDownloadLink"

	if err := r.validateExists(ctx, key); err != nil {
		return "", berrors.FromErr(op, err)
	}

	req := s3.GetObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	}

	presigned, err := r.presigner.
		PresignGetObject(ctx, &req, s3.WithPresignExpires(time.Duration(r.presignHours)*time.Hour))
	if err != nil {
		return "", berrors.InternalFromErr(op, err)
	}

	return presigned.URL, nil
}

func (r *FileS3RepositoryAdapter) validateExists(ctx context.Context, key string) error {
	const op = "repo - s3 - validateExists"

	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	})
	if err != nil {
		var nsk *types.NotFound
		if errors.As(err, &nsk) {
			return berrors.FromErr(op, file_errors.ErrFileNotFound)
		}

		return berrors.InternalFromErr(op, err)
	}

	return nil
}
