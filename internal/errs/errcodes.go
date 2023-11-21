package errs

import "errors"

var ErrUniqueLinkCode = errors.New("unique link violation")
var ErrEmptyURL = errors.New("url is empty")
var ErrKeyNotFound = errors.New("key not found")
