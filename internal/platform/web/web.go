package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	mux *chi.Mux
	log *log.Logger
}

func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

// handle connects method and url pattern to particular application handler
func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
