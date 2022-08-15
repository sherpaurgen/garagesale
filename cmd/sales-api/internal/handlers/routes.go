package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sherpaurgen/garagesale/internal/platform/web"
)

// API construct handler that knows all api routes
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)
	p := ProductService{DB: db, Log: logger}
	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.GetProduct)

	return app
}
