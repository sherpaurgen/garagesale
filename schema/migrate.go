package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

// migration contain query needed to create table
// Never remove entries from migration SLICE once they have been run in production
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `
		CREATE TABLE products(
			product_id UUID,
			name TEXT,
			cost INT,
			quantity INT,
			date_created TIMESTAMP,
			date_updated TIMESTAMP,
			PRIMARY KEY (product_id)
		);`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
