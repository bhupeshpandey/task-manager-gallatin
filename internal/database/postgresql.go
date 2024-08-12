package database

import (
	"database/sql"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) models.TaskRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (id, parent_id, title, description, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, task.ID, task.ParentID, task.Title, task.Description, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *postgresRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, task.Title, task.Description, task.UpdatedAt, task.ID)
	return err
}

func (r *postgresRepository) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	task := &models.Task{}
	query := `SELECT id, parent_id, title, description, created_at, updated_at FROM tasks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.ParentID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt)
	return task, err
}

func (r *postgresRepository) DeleteTask(id uuid.UUID) error {
	query := "DELETE FROM `tasks` WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postgresRepository) ListTasks() ([]*models.Task, error) {
	query := "SELECT id, parent_id, title, description, created_at, updated_at FROM tasks"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(&task.ID, &task.ParentID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
