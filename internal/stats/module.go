package stats

import (
	"context"
	"sync"

	"github.com/tnfy-link/client-go/queue"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"stats",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("stats")
	}),
	fx.Provide(newRepository, fx.Private),
	fx.Provide(queue.NewStatsQueue),
	fx.Provide(NewService),
	fx.Invoke(func(lc fx.Lifecycle, s *Service, log *zap.Logger) {
		wg := sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				wg.Add(1)
				go func() {
					defer wg.Done()

					log.Info("stats processing started")
					err := s.Process(ctx)
					if err != nil && err != context.Canceled {
						log.Error("failed to process stats", zap.Error(err))
					}
					log.Info("stats processing stopped")
				}()
				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				wg.Wait()
				return nil
			},
		})
	}),
)
