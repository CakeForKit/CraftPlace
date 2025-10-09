package dberrors

import "errors"

var (
	ErrQueryBuilds  = errors.New("query build failed")
	ErrQueryExec    = errors.New("query execution failed")
	ErrExpectedOne  = errors.New("expected one")
	ErrRowsAffected = errors.New("no rows affected")
	ErrOpenConnect  = errors.New("open connect failed")
	ErrPing         = errors.New("ping failed")
)
