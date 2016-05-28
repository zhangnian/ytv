package utils

import (
	"database/sql"
	"time"
)

func DefaultString(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

func CurTimeStr() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)

	return tm.Format("2006-01-02 15:04:05")
}
