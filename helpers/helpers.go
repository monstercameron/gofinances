package helpers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// SetupSignalHandling configures handling of interrupt signals
func SetupSignalHandling(mux *http.ServeMux) (*http.Server, <-chan struct{}) {
	done := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create an http.Server instance using the given ServeMux
	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	go func() {
		<-signalChan
		fmt.Println("Signal received, starting graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("HTTP server Shutdown: %v\n", err)
		}

		close(done)
	}()

	return server, done
}

// ExtractSegmentFromPath extracts a segment from the URL path at a given index.
func ExtractSegmentFromPath(path string, index int) (string, error) {
	pathSegments := strings.Split(path, "/")
	// Check if the requested index is within the range of available segments.
	if index > 0 && index < len(pathSegments) {
		return pathSegments[index], nil
	}
	return "", fmt.Errorf("segment not found at index %d", index)
}
