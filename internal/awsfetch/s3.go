package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Bucket represents an S3 bucket.
type S3Bucket struct {
	Name     string `json:"Name"`
	Location string `json:"Location"`
	// Additional fields (e.g., Tags, Policy) can be added as needed.
}

// FetchS3Buckets retrieves all S3 buckets.
func FetchS3Buckets(ctx context.Context, cfg aws.Config) ([]S3Bucket, error) {
	client := s3.NewFromConfig(cfg)
	out, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching S3 buckets: %w", err)
	}
	var buckets []S3Bucket
	for _, b := range out.Buckets {
		bucket := S3Bucket{
			Name:     *b.Name,
			Location: "us-east-1", // Default; ideally, call GetBucketLocation
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}
