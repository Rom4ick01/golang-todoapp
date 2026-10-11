// Package tasks_postgres_repository — хранение задач в PostgreSQL (одна операция — один файл).
package tasks_postgres_repository

import core_postgres_pool "todoapp/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
