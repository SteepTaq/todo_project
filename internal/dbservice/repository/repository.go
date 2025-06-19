package repository

import (
	"database/sql"

	"github.com/SteepTaq/todo_project/internal/dbservice/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTodo(id string) (*model.Todo, error) {
	row := r.db.QueryRow("SELECT * FROM todos WHERE id = $1", id)
	var todo model.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
	return &todo, err
}

func (r *PostgresRepository) CreateTodo(todo *model.Todo) error {
	_, err := r.db.Exec("INSERT INTO todos (id, title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", todo.ID, todo.Title, todo.Description, todo.Status, todo.CreatedAt, todo.UpdatedAt)
	return err
}

func (r *PostgresRepository) UpdateTodo(todo *model.Todo) error {
	_, err := r.db.Exec("UPDATE todos SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5", todo.Title, todo.Description, todo.Status, todo.UpdatedAt, todo.ID)
	return err
}

func (r *PostgresRepository) DeleteTodo(id string) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) GetAllTodos() ([]*model.Todo, error) {
	rows, err := r.db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

