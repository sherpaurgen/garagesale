package schema

import (
	"github.com/jmoiron/sqlx"
)

const seeds = `
INSERT INTO products (product_id,name,cost,quantity,date_created,date_updated) VALUES ('01caf1ef-964a-41cd-804e-56c352f3733e','Comicbook',44,55,'2020-12-12 00:00:12.000001+00','2020-12-12 00:03:12.000001+00'),('c53f9e14-18be-11ed-861d-0242ac120002','Marvel',24,15,'2020-12-12 00:00:12.000001+00','2020-12-12 00:03:12.000001+00') ON CONFLICT DO NOTHING;
`

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
