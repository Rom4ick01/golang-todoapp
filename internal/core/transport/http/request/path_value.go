package core_http_request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_errors "todoapp/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.UUID{}, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf(
			"path value='%s' by key='%s' not a valid uuid: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
