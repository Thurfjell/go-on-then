package main

import (
	"context"
	"goonthen/internal/data"
	"goonthen/internal/server"
	"goonthen/internal/server/game"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gameService, err := game.NewService(data.NewInMemoryState())

	if err != nil {
		log.Panic(err.Error())
	}

	server := server.New(server.WithRoutes(gameService.Routes()))

	go func() {
		log.Printf("Starting server at %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalln("Oh no server shutdown errored\n", err)
		}
	}()

	// Tidy tidy tidy
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Stopping server. Bye!")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shut down error: %v", err)
	}

}
