package repository

import (
	"errors"
)

// for backward compatibility
var (
	ErrMethodUnimplemented = errors.New("unimplemented")
	ErrNoData              = errors.New("got no data from repository layer")
	ErrNoDataAffected      = errors.New("no repository data affected by affecting query")
)
