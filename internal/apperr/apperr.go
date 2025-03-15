package apperr

import "github.com/pkg/errors"


var (
	ErrNoData = errors.New("no data found")
	ErrNoDataCache = errors.New("no data found in cache")
)