package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func GracefulShutdown(ctx context.Context, wg *sync.WaitGroup, srv *http.Server, srvErrCh <-chan error) {
	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown error: %v\n", err)
		}

		wg.Wait()
		return

	case err := <-srvErrCh:
		if err == nil {
			log.Println("Server stopped without error")
			wg.Wait()
			return
		}

		if err == http.ErrServerClosed {
			log.Println("Server closed")
			wg.Wait()
			return
		}

		log.Fatalf("Server error: %v\n", err)
	}
}
