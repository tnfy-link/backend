package stats

import (
	"context"
	"time"

	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/client-go/api"
	"github.com/tnfy-link/client-go/queue"
	"go.uber.org/zap"
)

type Service struct {
	stats *repository

	links *links.Service

	queue *queue.StatsQueue

	log *zap.Logger
}

func (s *Service) Get(ctx context.Context, id string) (api.Stats, error) {
	return s.stats.Get(ctx, id)
}

func (s *Service) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		event, err := s.queue.Dequeue(ctx)
		if err != nil {
			if err != queue.ErrEmptyQueue {
				s.log.Error("failed to dequeue event", zap.Error(err))
			}

			continue
		}

		s.log.Info("processing event", zap.Any("event", event))

		if err := s.processEvent(event); err != nil {
			s.log.Error("failed to process event", zap.Error(err))
			continue
		}
	}
}

func (s *Service) processEvent(event queue.StatsIncrEvent) error {
	subCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	link, err := s.links.Get(subCtx, event.LinkID)
	if err != nil {
		return err
	}

	return s.stats.Incr(subCtx, link, event.Labels)
}

func NewService(stats *repository, links *links.Service, queue *queue.StatsQueue, log *zap.Logger) *Service {
	if stats == nil {
		panic("stats repository is required")
	}
	if links == nil {
		panic("links service is required")
	}
	if queue == nil {
		panic("queue is required")
	}
	if log == nil {
		panic("logger is required")
	}

	return &Service{
		stats: stats,

		links: links,

		queue: queue,

		log: log,
	}
}
