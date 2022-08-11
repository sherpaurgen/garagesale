package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sherpaurgen/garagesale/schema"
)

type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type ProductService struct {
	db *sqlx.DB
}

func openDb() (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
	//returns an   *sqlx.DB.
}

func main() {
	//connection initialization
	db, err := openDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Error in migration:", err)
		}
		log.Println("Migration Complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Error applying seed data:", err)
		}
		log.Println("Seed applied")
		return
	}
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
	list := []Product{} //initialize empty list other wise it will return nil which is fine but difficult for frontend to parse json, doing this will return empty list []
	const q = `SELECT * FROM products`
	err := p.db.Select(&list, q)
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
