package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sherpaurgen/garagesale/internal/product"
)

type ProductService struct {
	db *sqlx.DB
}

func main() {

	ps := ProductService{db: db}

	api := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(ps.List),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
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
		log.Fatalf("Error while listening and starting http server: %v", err)

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
			log.Fatalf("main: could not stop server gracefully Error: %v", err)
		}

	}

}

func (p *ProductService) List(w http.ResponseWriter, req *http.Request) {

	list, err := product.List(p.db)
	const q = `SELECT * FROM products`
	err = p.db.Select(&list, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Marshalling error:", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Print("error writing ", err)
	}
}
