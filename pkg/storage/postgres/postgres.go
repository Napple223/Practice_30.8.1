package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// абстракция на БД
type Storage struct {
	db *pgxpool.Pool
}

// функция - конструктор для подключения к БД
func New(constr string) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), constr)

	if err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}
	return &s, nil
}

// структура задачи
type Task struct {
	ID         int
	Opened     int
	Closed     int
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Функция возвращает все имеющиеся в БД задачи
func (s *Storage) ReturnTasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT
		id,
		opened,
		closed,
		author_id,
		assigned_id,
		title,
		content
	FROM tasks
	WHERE
		($1=0 OR id=$1) AND
		($2=0 OR author_id=$2)
	ORDER BY id;
	`,
		taskID, authorID,
	)

	if err != nil {
		return nil, err
	}

	var tasks []Task

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

// функция дообавляет новую задачу и возвращает ее id
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(),
		`
		INSERT INTO tasks (title, content) VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}
