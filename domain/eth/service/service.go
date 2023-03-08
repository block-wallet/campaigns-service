package service

import "context"

type EventsService interface {
	Watch(ctx context.Context) error
}
