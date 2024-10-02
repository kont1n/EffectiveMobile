package api

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func StartHttpServer(wg *sync.WaitGroup, address string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		defer wg.Done()

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe() error: %v", err)
		}
	}()

	return srv
}
