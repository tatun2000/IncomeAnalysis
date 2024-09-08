package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const cfpPath = "./config"

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	app, gracefulShutdown, err := InjectAppGod(ctx, cfpPath)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		app.telegramService.Run(ctx)
	}()

	<-ctx.Done()
	gracefulShutdown()
}
