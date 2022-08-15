package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sherpaurgen/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sherpaurgen/garagesale/internal/platform/conf"
	"github.com/sherpaurgen/garagesale/internal/platform/database"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	//connection initialization
	db, err := database.OpenDb(conf.GetDbConfig())
	if err != nil {
		return errors.Wrap(err, "DB connection issue")
	}
	defer db.Close()
	log := log.New(os.Stdout, "SalesAPI:", log.LstdFlags)
	//logging based on method signature eg. salesapi
	// Goal is to pass logger to the things that need it here it ispassed to handler as struct

	ps := handlers.ProductService{DB: db, Log: log}

	var webconf conf.Webconfig
	webconf = conf.GetWebConfig()
	api := &http.Server{
		Addr:           webconf.Addr,
		Handler:        http.HandlerFunc(ps.List),
		ReadTimeout:    time.Duration(webconf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(webconf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	go func() {
		log.Printf("Main api listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	//above is not blocking operation although waiting on shutdown channel
	select {
	case err := <-serverErrors:
		//log.Fatalf("Error while listening and starting http server: %v", err)
		return errors.Wrap(err, "Issue listening and starting http server")

	case <-shutdown:
		log.Println("main: Starting shutdown")
		const timeout = 5 * time.Second
		// Context - it is used for Cancellation and propagation, the context.Background() gives empty context
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := api.Shutdown(ctx)
		/*Shutdown gracefully shuts down the server without interrupting any active connections.*/
		if err != nil {
			log.Printf("main: Graceful shutdown didnot complete in %v:%v", timeout, err)
			err = api.Close()
			//Close() immediately closes all active net.Listeners and any connections in state StateNew, StateActive, or StateIdle. For a graceful shutdown, use Shutdown.
		}
		if err != nil {
			//log.Fatalf("main: could not stop server gracefully Error: %v", err)
			return errors.Wrap(err, "Main could not stop server gracefully")
		}

	} //
	return nil //since log fatal is taken care of
}
