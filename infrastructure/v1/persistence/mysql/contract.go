package mysql

import (
	"database/sql"
	"dot/core"
	"gorm.io/gorm"
)

//go:generate mockery --name Database --filename database.go --output ./mocks
type Database interface {
	Table(name string, args ...interface{}) Database
	WithContext(ic *core.InternalContext) Database
	Model(value interface{}) Database
	Create(value interface{}) Database
	Where(query interface{}, args ...interface{}) Database
	Save(value interface{}) Database
	Find(dest interface{}, conds ...interface{}) Database
	First(dest interface{}, conds ...interface{}) Database
	Take(dest interface{}, conds ...interface{}) Database
	Update(column string, value interface{}) Database
	Updates(values interface{}) Database
	Omit(columns ...string) Database
	Delete(value interface{}, conds ...interface{}) Database
	Exec(sql string, values ...interface{}) Database
	Raw(sql string, values ...interface{}) Database
	Scan(dest interface{}) Database
	ScanRows(rows *sql.Rows, dest interface{}) error
	Limit(limit int) Database
	Offset(limit int) Database
	Count(count *int64) Database
	Begin(opts ...*sql.TxOptions) Database
	Rollback() Database
	Commit() Database
	Error() error
	RowsAffected() int64
	DB() *gorm.DB
}

//go:generate mockery --name Rows --filename rows.go --output ./mocks
type Rows interface {
	Columns() ([]string, error)
	ColumnTypes() ([]*sql.ColumnType, error)
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close() error
}
