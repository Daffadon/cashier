package dto

import "errors"

var (
	ErrBadrequest = errors.New("Missing or invalid request")
	ErrToSaveFile = errors.New("Something wrong when saving file")
)
