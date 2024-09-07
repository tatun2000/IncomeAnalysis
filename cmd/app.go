package main

import (
	"context"

	"telegrammbot.core/internal/domains/telegram"
)

type tracerDummy struct{}

type app struct {
	tracer tracerDummy // only for wire injection

	telegramService *telegram.Service
}

func newTracer(ctx context.Context) (tracerDummy, func(), error) {
	return struct{}{}, func() {
	}, nil
}
