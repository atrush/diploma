package main

import (
	"context"
	"github.com/atrush/diploma.git/api"
	"github.com/atrush/diploma.git/pkg"
	"github.com/atrush/diploma.git/storage/psql"
	"log"
	"os"
	"os/signal"
)

func main() {

	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	//fmt.Printf("config dsn: %v" \n, cfg.DatabaseDSN)

	db, err := psql.NewStorage(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := api.NewServer(cfg, db)
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
