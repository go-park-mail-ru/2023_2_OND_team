package repository

import "errors"

const (
	TimeFormat = "02.01.2006"
)

var (
	ErrMethodUnimplemented = errors.New("unimplemented")
	ErrNoData              = errors.New("got no data from repository layer")
)
