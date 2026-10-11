// Package core_errors — общие ошибки приложения. HTTP-слой переводит их в статусы (см. HTTPResponseHandler.ErrorResponse).
package core_errors

import "errors"

var (
	// HTTP 404
	ErrNotFound = errors.New("not found")

	// HTTP 400
	ErrInvalidArgument = errors.New("invalid argument")

	// HTTP 409
	ErrConflict = errors.New("conflict")
)
