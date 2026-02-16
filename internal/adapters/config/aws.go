package config

import (
	"errors"

	"github.com/spf13/viper"
)

type AWSConfig struct {
	Region   string
	QueueURL string
}

func NewAWSConfig() (*AWSConfig, error) {
	region := viper.GetString("aws.region")
	queueURL := viper.GetString("aws.sqs.queue_url")

	if region == "" || queueURL == "" {
		return nil, errors.New("AWS region and SQS queue URL must be provided in config")
	}

	return &AWSConfig{
		Region:   region,
		QueueURL: queueURL,
	}, nil
}
