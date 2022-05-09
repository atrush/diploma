package main

import (
	"context"
	"errors"
	"github.com/atrush/diploma.git/api"
	"github.com/atrush/diploma.git/pkg"
	prov_accual "github.com/atrush/diploma.git/provider/accrual"
	"github.com/atrush/diploma.git/services/accrual"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/atrush/diploma.git/services/order"
	"github.com/atrush/diploma.git/storage/psql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	//fmt.Printf("config dsn: %v" \n, cfg.DatabaseDSN)

	//cont, err := testtools.NewPostgreSQLContainer(context.Background())
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//tt := time.Minute
	//defer cont.Stop(context.Background(), &tt)

	db, err := psql.NewTestStorage()
	if err != nil {
		log.Fatal(err.Error())
	}

	jwtAuth, err := auth.NewAuth(db)
	if err != nil {
		log.Fatalf("error starting auth service:%v", err.Error())
	}

	svcOrder, err := order.NewOrder(db)
	if err != nil {
		log.Fatalf("error starting order service:%v", err.Error())
	}

	accProv, err := prov_accual.NewAccrual("http://localhost:8090/api/orders")
	if err != nil {
		log.Fatalf("error starting accrual provider:%v", err.Error())
	}

	svcAccrual := accrual.NewAccrualService(svcOrder, accProv)
	if err != nil {
		log.Fatalf("error starting accrual service:%v", err.Error())
	}

	accrualCtx, accrualClose := context.WithCancel(context.Background())
	svcAccrual.Run(accrualCtx)
	if err != nil {
		log.Fatalf("error starting accrual worker:%v", err.Error())
	}

	server, err := api.NewServer(cfg, jwtAuth, svcOrder)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := server.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accrualClose()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}

}
