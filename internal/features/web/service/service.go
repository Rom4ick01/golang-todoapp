// Package web_service — отдача веб-страниц приложения.
package web_service

import "todoapp/internal/core/domain"

type WebService struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) (domain.File, error)
}

func NewWebService(
	webRepository WebRepository,
) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}
