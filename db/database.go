package db

import "database/sql"

type DbInstance struct {
	*sql.DB
}
