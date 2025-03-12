package errorx

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrDuplicateEntryCode = 1062
)

func mySQLErrCode(err error) int {
	var mysqlErr *mysql.MySQLError
	ok := errors.As(err, &mysqlErr)
	if !ok {
		return 0
	}
	return int(mysqlErr.Number)
}

func IsMySQLDuplicateEntry(err error) bool {
	return mySQLErrCode(err) == ErrDuplicateEntryCode
}
