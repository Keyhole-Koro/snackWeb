package db

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Client *dynamodb.Client
var TableName string

func InitDB() error {
	ctx := context.TODO()

	// Load AWS Config
	// Use custom endpoint resolver if AWS_ENDPOINT_URL is set (for LocalStack)
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return err
	}

	Client = dynamodb.NewFromConfig(cfg)

	TableName = os.Getenv("DYNAMODB_TABLE")
	if TableName == "" {
		TableName = "SnackTable"
		log.Println("DYNAMODB_TABLE not set, defaulting to SnackTable")
	}

	log.Println("DynamoDB Client initialized. Table:", TableName)
	return nil
}

func CloseDB() {
	// AWS SDK v2 clients do not need explicit closing
}
