package async_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/marcodd23/gopernet/internal/async"
)

// MockGopherService already defined above

func TestBackgroundTaskManager_StartBurrowUpdater(t *testing.T) {
	mockService := new(MockGopherService)
	mockService.On("UpdateBurrows").Return()

	taskManager := async.NewBackgroundTaskManager(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	interval := 10 * time.Millisecond

	// Start the BurrowUpdater
	taskManager.StartBurrowUpdater(ctx, &wg, interval)

	// Wait for the task to run at least once
	time.Sleep(25 * time.Millisecond)

	// Cancel the context to stop the goroutine
	cancel()

	// Wait for the goroutine to finish
	wg.Wait()

	// Verify that the UpdateBurrows method was called
	mockService.AssertCalled(t, "UpdateBurrows")
}

func TestBackgroundTaskManager_StartPeriodicSaver(t *testing.T) {
	mockService := new(MockGopherService)
	mockService.On("SaveState").Return(nil)

	taskManager := async.NewBackgroundTaskManager(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	interval := 10 * time.Millisecond

	// Start the PeriodicSaver
	taskManager.StartPeriodicSaver(ctx, &wg, interval)

	// Wait for the task to run at least once
	time.Sleep(25 * time.Millisecond)

	// Cancel the context to stop the goroutine
	cancel()

	// Wait for the goroutine to finish
	wg.Wait()

	// Verify that the SaveState method was called
	mockService.AssertCalled(t, "SaveState")
}

func TestBackgroundTaskManager_StartReportGenerator(t *testing.T) {
	mockService := new(MockGopherService)
	mockService.On("SaveReport").Return(nil)

	taskManager := async.NewBackgroundTaskManager(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	interval := 10 * time.Millisecond

	// Start the ReportGenerator
	taskManager.StartReportGenerator(ctx, &wg, interval)

	// Wait for the task to run at least once
	time.Sleep(25 * time.Millisecond)

	// Cancel the context to stop the goroutine
	cancel()

	// Wait for the goroutine to finish
	wg.Wait()

	// Verify that the SaveReport method was called
	mockService.AssertCalled(t, "SaveReport")
}
