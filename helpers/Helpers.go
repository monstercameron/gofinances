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

// ThDate converts an integer date to a string and appends the correct English ordinal indicator.
func ThDate(date int) string {
	// Handle special cases: 11th, 12th, and 13th are exceptions in English
	if date%100 == 11 || date%100 == 12 || date%100 == 13 {
		// Use "th" for these special cases
		return fmt.Sprintf("%dth", date)
	}

	// Calculate the last digit of the date for determining the ordinal indicator
	lastDigit := date % 10

	// Determine and append the correct ordinal indicator based on the last digit
	switch lastDigit {
	case 1:
		// Use "st" for dates ending in 1 (except those ending in 11)
		return fmt.Sprintf("%dst", date)
	case 2:
		// Use "nd" for dates ending in 2 (except those ending in 12)
		return fmt.Sprintf("%dnd", date)
	case 3:
		// Use "rd" for dates ending in 3 (except those ending in 13)
		return fmt.Sprintf("%drd", date)
	default:
		// Use "th" for all other dates
		return fmt.Sprintf("%dth", date)
	}
}
