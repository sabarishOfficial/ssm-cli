package awsconfig

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)
func SetupAWSConfig(region string) (aws.Config, error) {
	if region == "" {
		region = "us-east-1"
	}
	fmt.Println("AWS Region:", region)

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return aws.Config{}, err
	}

	fmt.Println("AWS configuration loaded successfully.")

	return cfg, nil
}