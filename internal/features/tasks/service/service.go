// Package tasks_service — бизнес-логика задач.
package tasks_service

import (
	"context"

	"github.com/google/uuid"
	"todoapp/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	SaveTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id uuid.UUID,
	) error

	UpdateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
