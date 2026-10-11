package statistics_service

import (
	"context"
	"fmt"
	"time"

	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be after from:%w", core_errors.ErrInvalidArgument)
		}
	}
	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository:%w", err)
	}
	statistics := calcStatistics(tasks)
	return statistics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}
	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletionDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}
	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100
	var tasksAverageCompletionTime *time.Duration

	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}
	return domain.Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         &tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
