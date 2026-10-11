// Package statistics_postgres_repository — чтение задач для статистики.
// Сам расчёт — в domain.CreateStatistics.
package statistics_postgres_repository

import core_postgres_pool "todoapp/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
