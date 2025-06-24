package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gcConfig "github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
)

// S3Client wraps AWS S3 client
type S3Client struct {
	Client *s3.Client
	Bucket string
}

// NewS3Client sets up AWS S3 client using values from config.yaml
func NewS3Client(ctx context.Context) (*S3Client, error) {
	region := gcConfig.GetString(ctx, "s3.region")
	endpoint := gcConfig.GetString(ctx, "s3.endpoint")
	bucket := gcConfig.GetString(ctx, "s3.bucket")

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:               endpoint,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)),
	)
	if err != nil {
		log.Errorf(" AWS config load failed: %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Needed for LocalStack
	})

	return &S3Client{
		Client: s3Client,
		Bucket: bucket,
	}, nil
}
