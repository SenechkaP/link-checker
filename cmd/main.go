package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	linkHandler "github.com/SenechkaP/link-checker/internal/handler/link"
	"github.com/SenechkaP/link-checker/internal/server"
	linkService "github.com/SenechkaP/link-checker/internal/service/link"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	router := chi.NewRouter()

	wg := &sync.WaitGroup{}
	service := linkService.NewLinkService(wg)
	handler := linkHandler.NewLinkHandler(service)

	router.Route("/links", func(r chi.Router) {
		r.Get("/statuses", handler.GetStatusesHandler)
		r.Get("/pdf", handler.GetStatusesByNumsHandler)
	})

	srv := &http.Server{Addr: ":8080", Handler: router}

	srvErrCh := make(chan error, 1)

	go func() {
		log.Printf("Server is launched on: %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			srvErrCh <- err
			return
		}
		srvErrCh <- nil
	}()

	server.GracefulShutdown(ctx, wg, srv, srvErrCh)

	log.Println("service stopped")
}
