package configurations

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/workers"
	"go.uber.org/fx"
)

func HookQueueClient(lifecycle fx.Lifecycle, client *asynq.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return client.Close()
		},
	})
}

func HookQueueServer(lifecycle fx.Lifecycle, server *asynq.Server, mux *asynq.ServeMux) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Run(mux); err != nil {
					panic(err)
				}
				fmt.Println("MuxServer is running")
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown()
			return nil
		},
	})
}

func HookStartWorker(lifecycle fx.Lifecycle, d *workers.DiscountWorker) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				d.Start()
			}()
			return nil
		},
	})
}
