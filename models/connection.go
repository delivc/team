package models

import (
	"github.com/delivc/team/storage"
	"github.com/gobuffalo/pop/v5"
)

// Pagination model
type Pagination struct {
	Page    uint64
	PerPage uint64
	Count   uint64
}

// Offset for pagination
func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.PerPage
}

// SortDirection holds Ascending or Descending cosnt
type SortDirection string

// Ascending sort direction
const Ascending SortDirection = "ASC"

// Descending sortdirection
const Descending SortDirection = "DESC"

// CreatedAt is a constant! ;O
const CreatedAt = "created_at"

// SortParams ?field,field,field
type SortParams struct {
	Fields []SortField
}

// SortField sort by what
type SortField struct {
	Name string
	Dir  SortDirection
}

// TruncateAll truncates all models
func TruncateAll(conn *storage.Connection) error {
	return conn.Transaction(func(tx *storage.Connection) error {
		return tx.RawQuery("TRUNCATE " + (&pop.Model{Value: Account{}}).TableName()).Exec()
	})
}
