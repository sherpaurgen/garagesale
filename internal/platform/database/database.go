package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register postgres driver
)

type DBconfig struct {
	Host       string `json:"Host"`
	Port       int    `json:"Port"`
	User       string `json:"User"`
	Pass       string `json:"Pass"`
	DBname     string `json:"DBname"`
	DisableTLS bool   `json:"DisableTLS,omitempty"`
}

func OpenDb(dbc DBconfig) (*sqlx.DB, error) {
	//Host, DBname, User, Pass string, Port int, DisableTLS bool

	q := url.Values{}
	DisableTLS := true
	//hardcoded since db is local docker
	q.Set("sslmode", "require")
	if DisableTLS {
		q.Set("sslmode", "disable")
	}

	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbc.User, dbc.Pass),
		Host:     dbc.Host,
		Path:     dbc.DBname,
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
	//returns an   *sqlx.DB.
}
