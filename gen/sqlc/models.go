// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"database/sql"
	"time"
)

type Company struct {
	ID        int32
	Name      string
	Address   sql.NullString
	Person    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID        int32
	Name      sql.NullString
	Price     sql.NullInt32
	CompanyID sql.NullInt32
	CreatedAt time.Time
	UpdatedAt time.Time
}
