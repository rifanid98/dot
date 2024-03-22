package model

import (
	"database/sql"
	"time"
)

func GetTime(tm time.Time, valid bool) sql.NullTime {
	return sql.NullTime{
		Time:  tm,
		Valid: valid,
	}
}
