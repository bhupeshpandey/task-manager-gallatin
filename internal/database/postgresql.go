package database

import (
	"database/sql"
	"fmt"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	_ "github.com/lib/pq"
	"strings"
)

type postgresRepository struct {
	db           *sql.DB
	defaultLimit int
}

func NewPostgresRepository(db *sql.DB) models.TaskRepository {
	return &postgresRepository{db: db, defaultLimit: 200}
}

func (r *postgresRepository) CreateTask(task *models.Task) error {
	//query := `INSERT INTO tasks (id, title, description, created_at, updated_at)
	//		  VALUES ($1, $3, $4, $5, $6)`
	//if task.ParentID != "" {
	query := `INSERT INTO tasks (id, parent_id, title, description, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, task.Id, task.ParentID, task.Title, task.Description, task.CreatedAt, task.UpdatedAt)
	return err
	//}
	//_, err := r.db.Exec(query, task.Id, task.ParentID, task.Title, task.Description, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *postgresRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, task.Title, task.Description, task.UpdatedAt, task.Id)
	return err
}

func (r *postgresRepository) GetTaskByID(id string) (*models.Task, error) {
	task := &models.Task{}
	query := `SELECT id, parent_id, title, description, created_at, updated_at FROM tasks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&task.Id, &task.ParentID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt)
	return task, err
}

func (r *postgresRepository) DeleteTask(id string) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postgresRepository) ListTasks(request *models.ListTasksRequest) ([]*models.Task, error) {
	var args []interface{}
	var conditions []string
	var queryBuilder strings.Builder

	// Start the base query
	queryBuilder.WriteString(`SELECT id, parent_id, title, description, created_at, updated_at FROM tasks`)

	// Add conditions based on the presence of startTime and endTime
	if request.StartTime != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, request.StartTime)
	}

	if request.EndTime != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", len(args)+1))
		args = append(args, request.EndTime)
	}

	// Add the WHERE clause if there are any conditions
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	// Add ORDER BY clause
	queryBuilder.WriteString(" ORDER BY created_at DESC")

	// Add LIMIT and OFFSET if provided
	if request.Limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)+1))
		args = append(args, request.Limit)
	}

	if request.Offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)+1))
		args = append(args, request.Offset)
	}

	// Execute the query
	query := queryBuilder.String()
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}

		if err = rows.Scan(&task.Id, &task.ParentID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		} else {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}
