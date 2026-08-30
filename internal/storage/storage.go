package storage

import (
	"TaskManager/internal/models"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskStorage interface {
	GetById(id int) (*models.Task, error)
	GetAll() ([]*models.Task, error)
	Create(t *models.Task) (int, error)
	Update(t *models.Task) error
	Delete(id int) error
}

var ErrNotFound = errors.New("not found")

type PostgresTaskStorage struct {
	mu sync.RWMutex
	db *pgxpool.Pool
}

func NewPostgresTaskStorage(db *pgxpool.Pool) *PostgresTaskStorage {
	return &PostgresTaskStorage{
		db: db,
	}
}

func (s *PostgresTaskStorage) GetById(id int) (*models.Task, error) {
	query := `
		SELECT id, title, description
		FROM tasks
		WHERE id =$1
	`

	task := &models.Task{}
	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("task %v: %w", id, ErrNotFound)
		}
		return nil, err
	}

	return task.Clone(), nil
}

func (s *PostgresTaskStorage) GetAll() ([]*models.Task, error) {
	query := `
		SELECT id, title, description
		FROM tasks
	`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0)

	for rows.Next() {
		task := &models.Task{}

		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *PostgresTaskStorage) Create(t *models.Task) (int, error) {
	query := `
		INSERT INTO tasks (title, description)
		VALUES ($1,$2)
		RETURNING id`

	err := s.db.QueryRow(
		context.Background(),
		query,
		t.Title,
		t.Description,
	).Scan(&t.Id)
	if err != nil {
		return 0, err
	}

	return t.Id, err
}

func (s *PostgresTaskStorage) Update(t *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2
		WHERE id = $3`

	result, err := s.db.Exec(
		context.Background(),
		query,
		t.Title,
		t.Description,
		t.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected != 1 {
		return fmt.Errorf("task %v: %w", t.Id, ErrNotFound)
	}

	return nil
}

func (s *PostgresTaskStorage) Delete(id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	`
	result, err := s.db.Exec(
		context.Background(),
		query,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected != 1 {
		return fmt.Errorf("task %v: %w", id, ErrNotFound)
	}

	return nil
}
