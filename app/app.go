package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gookit/color"

	o "github.com/denniswon/tcex/app/order"
	"github.com/denniswon/tcex/app/rest"
)

// Run - Application to be invoked from main runner using this function
func Run(configFile string) {

	ctx, cancel := context.WithCancel(context.Background())
	orderQueue, publishQueue, _redis := bootstrap(configFile)

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

	go orderQueue.Start(ctx)
	defer orderQueue.Stop()

	go o.ProcessOrderReplays(publishQueue, _redis)

	// Starting http server on main thread
	rest.RunHTTPServer(orderQueue, _redis)
}
