package utils

import (
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SharedState = struct {
	Pool *pgxpool.Pool
}

// Get name of any function `skip` amount away from this function's position
func GetFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}

	funcPtr := runtime.FuncForPC(pc)
	if funcPtr == nil {
		return ""
	}

	return funcPtr.Name()
}
