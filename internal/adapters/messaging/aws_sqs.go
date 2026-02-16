package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/golang-auth/internal/core/ports"
)

type SQSAdapter struct {
	client   *sqs.Client
	queueURL string
	logger   ports.Logger
}

func NewSQSAdapter(client *sqs.Client, queueURL string, logger ports.Logger) ports.EventPublisher {
	return &SQSAdapter{
		client:   client,
		queueURL: queueURL,
		logger:   logger,
	}
}

func (a *SQSAdapter) PublishUserRegistered(ctx context.Context, email string, token string) error {
	payload := map[string]string{
		"email": email,
		"token": token,
		"event": "user.registered",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SQS message: %w", err)
	}
	_ = body
	a.logger.Info("SQS Message: ", token)
	// _, err = a.client.SendMessage(ctx, &sqs.SendMessageInput{
	// 	QueueUrl:    aws.String(a.queueURL),
	// 	MessageBody: aws.String(string(body)),
	// 	// Optional: Add MessageAttributes for filtering
	// })

	// if err != nil {
	// 	return fmt.Errorf("failed to send message to SQS: %w", err)
	// }

	return nil
}
