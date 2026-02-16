package ports

import (
	"context"
)

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, email string, token string) error
}
