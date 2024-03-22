package mysql

import (
	"database/sql"
	"dot/core"
	"gorm.io/gorm"
)

type databaseImpl struct {
	db *gorm.DB
}

func NewDatabase(database *gorm.DB) *databaseImpl {
	return &databaseImpl{database}
}

func (database *databaseImpl) Table(name string, args ...interface{}) Database {
	db := database.db.Table(name, args...)
	database.db = db
	return database
}

func (database *databaseImpl) WithContext(ic *core.InternalContext) Database {
	db := database.db.WithContext(ic.ToContext())
	database.db = db
	return database
}

func (database *databaseImpl) Model(value interface{}) Database {
	db := database.db.Model(value)
	database.db = db
	return database
}

func (database *databaseImpl) Create(value interface{}) Database {
	db := database.db.Create(value)
	database.db = db
	return database
}

func (database *databaseImpl) Where(query interface{}, args ...interface{}) Database {
	db := database.db.Where(query, args...)
	database.db = db
	return database
}

func (database *databaseImpl) Save(value interface{}) Database {
	db := database.db.Save(value)
	database.db = db
	return database
}

func (database *databaseImpl) Find(dest interface{}, conds ...interface{}) Database {
	db := database.db.Find(dest, conds...)
	database.db = db
	return database
}

func (database *databaseImpl) First(dest interface{}, conds ...interface{}) Database {
	db := database.db.First(dest, conds...)
	database.db = db
	return database
}

func (database *databaseImpl) Take(dest interface{}, conds ...interface{}) Database {
	db := database.db.Take(dest, conds...)
	database.db = db
	return database
}

func (database *databaseImpl) Update(column string, value interface{}) Database {
	db := database.db.Update(column, value)
	database.db = db
	return database
}

func (database *databaseImpl) Updates(values interface{}) Database {
	db := database.db.Updates(values)
	database.db = db
	return database
}

func (database *databaseImpl) Omit(columns ...string) Database {
	db := database.db.Omit(columns...)
	database.db = db
	return database
}

func (database *databaseImpl) Delete(value interface{}, conds ...interface{}) Database {
	db := database.db.Delete(value, conds...)
	database.db = db
	return database
}

func (database *databaseImpl) Exec(sql string, values ...interface{}) Database {
	db := database.db.Exec(sql, values...)
	database.db = db
	return database
}

func (database *databaseImpl) Raw(sql string, values ...interface{}) Database {
	db := database.db.Raw(sql, values...)
	database.db = db
	return database
}

func (database *databaseImpl) Scan(dest interface{}) Database {
	db := database.db.Scan(dest)
	database.db = db
	return database
}

func (database *databaseImpl) ScanRows(rows *sql.Rows, dest interface{}) error {
	return database.db.ScanRows(rows, dest)
}

func (database *databaseImpl) Limit(limit int) Database {
	db := database.db.Limit(limit)
	database.db = db
	return database
}

func (database *databaseImpl) Offset(limit int) Database {
	db := database.db.Offset(limit)
	database.db = db
	return database
}

func (database *databaseImpl) Count(count *int64) Database {
	db := database.db.Count(count)
	database.db = db
	return database
}

// ------------- TRANSACTION -------------

func (database *databaseImpl) Begin(opts ...*sql.TxOptions) Database {
	//return NewDatabase(tx)
	//tx := db.db.Begin(opts...)
	db := database.db.Begin(opts...)
	database.db = db
	return database
}

func (database *databaseImpl) Rollback() Database {
	db := database.db.Rollback()
	database.db = db
	return database
}

func (database *databaseImpl) Commit() Database {
	db := database.db.Commit()
	database.db = db
	return database
}

// ------------- GETTER -------------

func (database *databaseImpl) Error() error {
	return database.db.Error
}

func (database *databaseImpl) RowsAffected() int64 {
	return database.db.RowsAffected
}

func (database *databaseImpl) DB() *gorm.DB {
	return database.db
}
