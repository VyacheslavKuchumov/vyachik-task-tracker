package tracker

import (
	"VyacheslavKuchumov/test-backend/types"
	"database/sql"
	"errors"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrForbidden = errors.New("forbidden")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateGoal(ownerID int, payload types.CreateGoalPayload) (*types.Goal, error) {
	row := s.db.QueryRow(
		`INSERT INTO goals (title, description, owner_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, title, description, owner_id, created_at`,
		payload.Title,
		payload.Description,
		ownerID,
	)
	return scanRowIntoGoal(row)
}

func (s *Store) GetGoalsByOwner(ownerID int) ([]*types.GoalWithTasks, error) {
	rows, err := s.db.Query(
		`SELECT
			g.id, g.title, g.description, g.owner_id, g.created_at,
			t.id, t.goal_id, t.title, t.description, t.status, t.assignee_id, t.created_by, t.created_at
		FROM goals g
		LEFT JOIN tasks t ON t.goal_id = g.id
		WHERE g.owner_id = $1
		ORDER BY g.created_at DESC, t.created_at ASC`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	goals := make([]*types.GoalWithTasks, 0)
	goalByID := make(map[int]*types.GoalWithTasks)

	for rows.Next() {
		var (
			goal       types.Goal
			taskID     sql.NullInt64
			taskGoalID sql.NullInt64
			taskTitle  sql.NullString
			taskDesc   sql.NullString
			taskStatus sql.NullString
			assigneeID sql.NullInt64
			createdBy  sql.NullInt64
			taskAt     sql.NullTime
		)

		if err := rows.Scan(
			&goal.ID,
			&goal.Title,
			&goal.Description,
			&goal.OwnerID,
			&goal.CreatedAt,
			&taskID,
			&taskGoalID,
			&taskTitle,
			&taskDesc,
			&taskStatus,
			&assigneeID,
			&createdBy,
			&taskAt,
		); err != nil {
			return nil, err
		}

		current, exists := goalByID[goal.ID]
		if !exists {
			current = &types.GoalWithTasks{
				Goal:  goal,
				Tasks: []*types.Task{},
			}
			goalByID[goal.ID] = current
			goals = append(goals, current)
		}

		if taskID.Valid {
			task := &types.Task{
				ID:          int(taskID.Int64),
				GoalID:      int(taskGoalID.Int64),
				Title:       taskTitle.String,
				Description: taskDesc.String,
				Status:      taskStatus.String,
				CreatedBy:   int(createdBy.Int64),
				CreatedAt:   taskAt.Time,
			}
			if assigneeID.Valid {
				value := int(assigneeID.Int64)
				task.AssigneeID = &value
			}
			current.Tasks = append(current.Tasks, task)
		}
	}

	return goals, rows.Err()
}

func (s *Store) CreateTask(goalID, creatorID int, payload types.CreateTaskPayload) (*types.Task, error) {
	row := s.db.QueryRow(
		`INSERT INTO tasks (goal_id, title, description, assignee_id, created_by)
		 SELECT g.id, $2, $3, $4, $5
		 FROM goals g
		 WHERE g.id = $1 AND g.owner_id = $5
		 RETURNING id, goal_id, title, description, status, assignee_id, created_by, created_at`,
		goalID,
		payload.Title,
		payload.Description,
		payload.AssigneeID,
		creatorID,
	)
	task, err := scanRowIntoTask(row)
	if err == sql.ErrNoRows {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) AssignTask(taskID, requesterID int, payload types.AssignTaskPayload) (*types.Task, error) {
	row := s.db.QueryRow(
		`UPDATE tasks t
		 SET assignee_id = $1
		 FROM goals g
		 WHERE t.goal_id = g.id AND t.id = $2 AND g.owner_id = $3
		 RETURNING t.id, t.goal_id, t.title, t.description, t.status, t.assignee_id, t.created_by, t.created_at`,
		payload.AssigneeID,
		taskID,
		requesterID,
	)

	task, err := scanRowIntoTask(row)
	if err == sql.ErrNoRows {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) GetAssignedTasks(userID int) ([]*types.Task, error) {
	rows, err := s.db.Query(
		`SELECT id, goal_id, title, description, status, assignee_id, created_by, created_at
		 FROM tasks
		 WHERE assignee_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task, err := scanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRowIntoGoal(row rowScanner) (*types.Goal, error) {
	g := new(types.Goal)
	if err := row.Scan(&g.ID, &g.Title, &g.Description, &g.OwnerID, &g.CreatedAt); err != nil {
		return nil, err
	}
	return g, nil
}

func scanRowIntoTask(row rowScanner) (*types.Task, error) {
	task := new(types.Task)
	var assigneeID sql.NullInt64
	if err := row.Scan(
		&task.ID,
		&task.GoalID,
		&task.Title,
		&task.Description,
		&task.Status,
		&assigneeID,
		&task.CreatedBy,
		&task.CreatedAt,
	); err != nil {
		return nil, err
	}
	if assigneeID.Valid {
		value := int(assigneeID.Int64)
		task.AssigneeID = &value
	}
	return task, nil
}
