package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

func run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	// opening a DB
	pg := NewPG()
	// closing DB must happen after all services are stopped, hence defer
	defer func() {
		if err := pg.Close(); err != nil {
			fmt.Printf("DB closed with error: %v\n", err)
		}
	}()
	if err := pg.Open(); err != nil {
		return err
	}

	// spin off a http server
	s := NewServer(ctx, ":8000")
	g.Go(s.Run)

	// spin off a couple of long-running services
	w1 := NewService(ctx, "Sender", 3*time.Second, pg)
	w2 := NewService(ctx, "Receiver", 2*time.Second, pg)
	g.Go(w1.Run)
	g.Go(w2.Run)

	// this blocks the main thread until everything but DB is stopped
	return g.Wait()
}

func main() {
	// global ctx shall be propagated to all child components
	ctx, cancel := context.WithCancel(context.Background())
	go ShutdownHandler(cancel)

	// blocks main thread until ctx is cancelled and all components are stopped gracefully
	if err := run(ctx); err != nil {
		fmt.Printf("Finished with error: %v\n", err)
	}

	fmt.Println("Exiting main()")
}
