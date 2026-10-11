package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	"todoapp/internal/core/domain"
	core_errors "todoapp/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) (domain.File, error) {
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.File{}, fmt.Errorf(
				"file: %s: %w",
				filePath,
				core_errors.ErrNotFound,
			)
		}

		return domain.File{}, fmt.Errorf(
			"get file: %s: %w",
			filePath,
			err,
		)
	}

	file := domain.NewFile(buffer)

	return file, nil
}
