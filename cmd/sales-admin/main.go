// this is for seed / migrating database
package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sherpaurgen/garagesale/schema"
)

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

}
