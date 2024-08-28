package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/denniswon/tcex/app/config"
	"github.com/gookit/color"

	o "github.com/denniswon/tcex/app/order"
	"github.com/denniswon/tcex/app/rest"
)

// Run - Application to be invoked from main runner using this function
func Run(configFile string) {

	ctx, cancel := context.WithCancel(context.Background())
	requestQueue, replayQueue, _redis := bootstrap(configFile)

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", cfg.GetUploadDirName())

	if err != nil {
		log.Print(color.Red.Sprintf("[!] Failed to create temporary directory for uploads : %s", err.Error()))
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	// Attempting to listen to Ctrl+C signal
	// and when received gracefully shutting down the service
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT)

	// All resources being used gets cleaned up
	// when we're returning from this function scope
	go func() {

		<-interruptChan

		// This call should be received in all places
		// where root context is passed along
		//
		// @note This can ( needs to ) be improved
		cancel()

		if err := _redis.Close(); err != nil {
			log.Print(color.Red.Sprintf("[!] Failed to close connection to Redis : %s", err.Error()))
			return
		}

		// Stopping process
		log.Print(color.Magenta.Sprintf("\n[+] Gracefully shut down the service"))
		os.Exit(0)

	}()

	go o.ProcessOrderReplays(ctx, requestQueue, replayQueue, _redis)

	// Starting http server on main thread
	rest.RunHTTPServer(requestQueue, _redis, tempDir)
}
