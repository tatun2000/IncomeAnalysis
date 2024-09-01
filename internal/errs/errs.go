package errs

import "errors"

var (
	ErrInvalidCategoryType         = errors.New("invalid category type")
	ErrInvalidDay                  = errors.New("invalid day")
	ErrInvalidMonth                = errors.New("invalid month")
	ErrUnknownRequestType          = errors.New("unknown request type")
	ErrInvalidRequestMessageFormat = errors.New("invalid request message format")
)
