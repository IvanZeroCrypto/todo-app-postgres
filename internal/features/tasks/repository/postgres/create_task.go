package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gitverse.ru/IvanVerse/go-todoApp/internal/core/domain"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	INSERT INTO todoapp.tasks (title,description,completed,created_at,completed_at,author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id,version,title,description,completed,created_at,completed_at,author_user_id;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return domain.Task{}, fmt.Errorf("%v: user with id=%d:%w", err, task.AuthorUserID, core_errors.ErrNotFound)
			}
		}
		return domain.Task{}, fmt.Errorf("scan error:%w", err)
	}
	taskDomain := taskDomainFromModel(taskModel)
	return taskDomain, nil
}
