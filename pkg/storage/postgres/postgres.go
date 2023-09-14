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

// Функция возвращает задачи по id или по автору
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

// Возвращает все задачи из БД
func (s *Storage) ReturnAllTasks() ([]Task, error) {
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
	ORDER BY id;
	`)

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

// Возвращает задачи по имени метки
func (s *Storage) ReturnTasksOnLabels(labelName string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT
	tasks.id,
	tasks.opened,
	tasks.closed,
	tasks.author_id,
	tasks.assigned_id,
	tasks.title,
	tasks.content
	FROM tasks
	JOIN tasks_labels ON tasks.id = tasks_labels.task_id
	JOIN labels ON labels.id = tasks_labels.labels_id
	WHERE labels.name = $1
	ORDER BY tasks.id;
	`,
		labelName,
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

// Обновляет задачу по id
func (s *Storage) UpdateTask(taskID int, t Task) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE tasks SET title = $1, content = $2
	WHERE id = $3;
	`,
		t.Title, t.Content, taskID,
	)

	return err
}

func (s *Storage) DeleteTask(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks_labels
	WHERE tasks_id = $1;
	DELETE FROM tasks 
	WHERE id = $1;
	`,
		taskID,
	)
	return err
}
