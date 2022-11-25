package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.rebrainme.com/golang_users_repos/2184/final/internal"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("fatal error: %v\n", err)
		}
	}()

	var paths []string
	if len(os.Args) > 3 {
		panic("wrong count arguments")
	}
	for i := 1; i < len(os.Args); i++ {
		paths = append(paths, os.Args[i])
	}

	startApp(paths[0], paths[1])
}

func startApp(origin, copy string) {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	defer func() {
		cancel()
	}()

	cron := internal.NewCron(500 * time.Millisecond)
	cron.Start(ctx, origin, copy)

	select {
	case err := <-errChan:
		log.Println(err)
	case <-stopChan:
		log.Println("stop app")
	}
}
