package customerrors

import "errors"

var ErrRecordNotFound = errors.New("could not find record")

var ErrRecordExists = errors.New("record already exists")

var ErrDataTooLong = errors.New("data is too long")

var ErrDataInvalid = errors.New("data is invalid")
