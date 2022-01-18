package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Server interface {
	// Run starts the server and stops it when context is cancelled.
	// This is a blocking call until the sever has stopped gracefully.
	Run() error
}

type server struct {
	ctx        context.Context
	httpServer *http.Server
}

func NewServer(ctx context.Context, addr string) Server {
	s := &server{
		ctx,
		&http.Server{
			Addr: addr,
		},
	}
	return s
}

func (s *server) Run() error {
	fmt.Println("Starting HTTP server...")

	g, gCtx := errgroup.WithContext(s.ctx)
	g.Go(func() error {
		err := s.httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})
	g.Go(func() error {
		<-gCtx.Done()
		err := s.httpServer.Shutdown(context.Background())
		fmt.Println("Stopped HTTP server")
		return err
	})

	return g.Wait()
}
