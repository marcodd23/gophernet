package async

import (
	"context"
	"github.com/marcodd23/go-micro-core/pkg/logmgr"
	"log"
	"sync"
	"time"

	"github.com/marcodd23/gopernet/internal/services"
)

type BackgroundTaskManager struct {
	service services.GopherService
}

func NewBackgroundTaskManager(service services.GopherService) *BackgroundTaskManager {
	return &BackgroundTaskManager{
		service: service,
	}
}

func (b *BackgroundTaskManager) StartBurrowUpdater(cancellableCtx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	wg.Add(1)
	ticker := time.NewTicker(interval)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				logmgr.GetLogger().LogDebug(cancellableCtx, "updating burrows ....")
				b.service.UpdateBurrows()
			case <-cancellableCtx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (b *BackgroundTaskManager) StartPeriodicSaver(cancellableCtx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	wg.Add(1)
	ticker := time.NewTicker(interval)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				logmgr.GetLogger().LogDebug(cancellableCtx, "saving state ....")
				err := b.service.SaveState()
				if err != nil {
					log.Printf("Error saving state: %v", err)
				}
			case <-cancellableCtx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (b *BackgroundTaskManager) StartReportGenerator(cancellableCtx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	wg.Add(1)
	ticker := time.NewTicker(interval)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				logmgr.GetLogger().LogDebug(cancellableCtx, "saving report ....")
				err := b.service.SaveReport()
				if err != nil {
					log.Printf("Error generating report: %v", err)
				}
			case <-cancellableCtx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
