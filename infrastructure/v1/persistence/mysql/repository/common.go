package repository

import (
	"dot/core"
	"dot/pkg/util"
	"gorm.io/gorm"
)

var log = util.NewLogger()

func db(ic *core.InternalContext, db *gorm.DB) *gorm.DB {
	ctxData := ic.GetData()
	if ctxData != nil {
		_db := ctxData["db"]
		if _db != nil {
			return _db.(*gorm.DB) // implements context.Context
		}
		return db
	}
	return db
}
