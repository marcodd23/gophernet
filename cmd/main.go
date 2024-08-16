package main

import (
	"context"
	"flag"
	"github.com/marcodd23/go-micro-core/pkg/logmgr"
	"github.com/marcodd23/go-micro-core/pkg/shutdown"
	"github.com/marcodd23/gopernet/internal/config"
	"net/http"
	"sync"
	"time"

	"github.com/marcodd23/gopernet/internal/api"
	"github.com/marcodd23/gopernet/internal/async"
	"github.com/marcodd23/gopernet/internal/repository"
	"github.com/marcodd23/gopernet/internal/services"
)

// ShutdownTimeoutMilli - timeout for cleaning up resources before shutting down the server.
const ShutdownTimeoutMilli = 500

func main() {
	rootCtx := context.Background()

	config := config.LoadConfiguration()

	logmgr.SetupLogger(config)

	// Define a string flag with a name, default value, and usage description.
	dataFile := flag.String("dataFile", "data/state.json", "Path to the initial state file")

	// Parse the command-line flags from os.Args (the arguments passed to the program).
	flag.Parse()

	// Initialize the repository
	memoryRepo := repository.NewMemoryRepository(*dataFile, "data/report.txt")

	// Initialize the service
	gopherNetService := services.NewGopherNetService(memoryRepo)

	// Load the initial state using the repository
	if err := gopherNetService.LoadInitialState(); err != nil {
		logmgr.GetLogger().LogError(rootCtx, "Failed to load initial state", err)
	}

	// Initialize background task manager
	backgroundTasks := async.NewBackgroundTaskManager(gopherNetService)

	// Set up cancelCtx and wait-group for managing goroutines
	cancelCtx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start background tasks
	backgroundTasks.StartBurrowUpdater(cancelCtx, &wg, time.Minute)
	backgroundTasks.StartPeriodicSaver(cancelCtx, &wg, 5*time.Minute)
	backgroundTasks.StartReportGenerator(cancelCtx, &wg, 5*time.Minute)

	// Create the server and define routes
	server := api.NewServer(gopherNetService, config)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logmgr.GetLogger().LogFatal(rootCtx, "Could not listen on :8080 \n", err)
		}
	}()

	shutdown.WaitForShutdown(rootCtx, ShutdownTimeoutMilli, func(timeoutCtx context.Context) {
		logmgr.GetLogger().LogInfo(timeoutCtx, "Shutting down server...")
		cancel() // Signal all goroutines to stop

		// Wait for all goroutines to finish
		wg.Wait()

		// Save state one last time before shutdown
		err := gopherNetService.SaveState()
		if err != nil {
			logmgr.GetLogger().LogInfo(timeoutCtx, "Error saving state during shutdown")
		}

		logmgr.GetLogger().LogInfo(timeoutCtx, "Server shut down cleanly")
	})
}
