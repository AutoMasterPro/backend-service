package s3

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Minio  *minio.Client
	Bucket string
	Region string
}

func New(endpoint, accessKey, secretKey, bucket, region string, useSSL bool) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket %s created\n", bucket)
	}

	return &Client{Minio: minioClient, Bucket: bucket, Region: region}, nil
}

func (c *Client) Upload(ctx context.Context, objectName string, fileData []byte, contentType string) error {
	reader := NewByteReader(fileData)

	_, err := c.Minio.PutObject(ctx, c.Bucket, objectName, reader, int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) GetURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	u, err := c.Minio.PresignedGetObject(ctx, c.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.Minio.RemoveObject(ctx, c.Bucket, objectName, minio.RemoveObjectOptions{})
}
