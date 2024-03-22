package mysql

import (
	"database/sql"
	"gorm.io/gorm"
)

type rowsImpl struct {
	gorm.Rows
}

func NewRows(rows gorm.Rows) *rowsImpl {
	return &rowsImpl{rows}
}

func (r *rowsImpl) Columns() ([]string, error) {
	return r.Rows.Columns()
}

func (r *rowsImpl) Next() bool {
	return r.Rows.Next()
}

func (r *rowsImpl) ColumnTypes() ([]*sql.ColumnType, error) {
	return r.Rows.ColumnTypes()
}

func (r *rowsImpl) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r *rowsImpl) Err() error {
	return r.Rows.Err()
}

func (r *rowsImpl) Close() error {
	return r.Rows.Close()
}
