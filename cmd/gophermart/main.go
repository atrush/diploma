package main

import (
	"context"
	"github.com/atrush/diploma.git/api"
	"log"
	"os"
	"os/signal"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(server.Run())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}

}
