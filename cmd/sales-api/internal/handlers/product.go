package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
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
	DB *sqlx.DB
}

func (p *ProductService) List(w http.ResponseWriter, req *http.Request) {
	list := []Product{} //initialize empty list other wise it will return nil which is fine but difficult for frontend to parse json, doing this will return empty list []
	const q = `SELECT * FROM products`
	err := p.DB.Select(&list, q)
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
