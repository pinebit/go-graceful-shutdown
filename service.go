package main

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	Run() error
}

type service struct {
	ctx      context.Context
	name     string
	interval time.Duration
	pg       PG
}

func NewService(ctx context.Context, name string, interval time.Duration, pg PG) Service {
	return &service{
		ctx,
		name,
		interval,
		pg,
	}
}

func (s *service) Run() error {
	var i int
	for {
		select {
		case <-s.ctx.Done():
			fmt.Printf("Stopped service: %s\n", s.name)
			return nil
		default:
		}

		if err := s.pg.Insert(s.ctx, i); err != nil {
			return err
		}
		i++

		fmt.Printf("Running service: %s\n", s.name)
		time.Sleep(s.interval)
	}
}
