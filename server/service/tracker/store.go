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
		`INSERT INTO goals (title, description, priority, status, owner_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, description, priority, status, owner_id, created_at`,
		payload.Title,
		payload.Description,
		normalizePriority(payload.Priority),
		normalizeGoalStatus(payload.Status),
		ownerID,
	)
	return scanRowIntoGoal(row)
}

func (s *Store) UpdateGoal(goalID, ownerID int, payload types.CreateGoalPayload) (*types.Goal, error) {
	row := s.db.QueryRow(
		`UPDATE goals
		 SET title = $1,
		     description = $2,
		     priority = $3,
		     status = $4
		 WHERE id = $5 AND owner_id = $6
		 RETURNING id, title, description, priority, status, owner_id, created_at`,
		payload.Title,
		payload.Description,
		normalizePriority(payload.Priority),
		normalizeGoalStatus(payload.Status),
		goalID,
		ownerID,
	)

	goal, err := scanRowIntoGoal(row)
	if err == sql.ErrNoRows {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *Store) DeleteGoal(goalID, ownerID int) error {
	result, err := s.db.Exec(
		`DELETE FROM goals
		 WHERE id = $1 AND owner_id = $2`,
		goalID,
		ownerID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrForbidden
	}
	return nil
}

func (s *Store) GetGoalsByOwner(_ int) ([]*types.GoalWithTasks, error) {
	rows, err := s.db.Query(
		`SELECT
			g.id,
			g.title,
			g.description,
			g.priority,
			g.status,
			g.owner_id,
			g.created_at,
			TRIM(CONCAT(owner_u.first_name, ' ', owner_u.last_name)) AS owner_name,
			t.id,
			t.goal_id,
			t.title,
			t.description,
			t.priority,
			t.is_completed,
			t.assignee_id,
			t.created_by,
			t.created_at,
			TRIM(CONCAT(assignee_u.first_name, ' ', assignee_u.last_name)) AS assignee_name,
			TRIM(CONCAT(creator_u.first_name, ' ', creator_u.last_name)) AS creator_name
		FROM goals g
		JOIN users owner_u ON owner_u.id = g.owner_id
		LEFT JOIN tasks t ON t.goal_id = g.id
		LEFT JOIN users assignee_u ON assignee_u.id = t.assignee_id
		LEFT JOIN users creator_u ON creator_u.id = t.created_by
		ORDER BY
			CASE WHEN g.status = 'achieved' THEN 1 ELSE 0 END,
			CASE g.priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
			g.created_at DESC,
			CASE WHEN t.is_completed THEN 1 ELSE 0 END,
			CASE t.priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 WHEN 'low' THEN 2 ELSE 3 END,
			t.created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	goals := make([]*types.GoalWithTasks, 0)
	goalByID := make(map[int]*types.GoalWithTasks)

	for rows.Next() {
		var (
			goal             types.Goal
			goalOwnerName    string
			taskID           sql.NullInt64
			taskGoalID       sql.NullInt64
			taskTitle        sql.NullString
			taskDesc         sql.NullString
			taskPriority     sql.NullString
			taskIsCompleted  sql.NullBool
			assigneeID       sql.NullInt64
			createdBy        sql.NullInt64
			taskAt           sql.NullTime
			taskAssigneeName sql.NullString
			taskCreatorName  sql.NullString
		)

		if err := rows.Scan(
			&goal.ID,
			&goal.Title,
			&goal.Description,
			&goal.Priority,
			&goal.Status,
			&goal.OwnerID,
			&goal.CreatedAt,
			&goalOwnerName,
			&taskID,
			&taskGoalID,
			&taskTitle,
			&taskDesc,
			&taskPriority,
			&taskIsCompleted,
			&assigneeID,
			&createdBy,
			&taskAt,
			&taskAssigneeName,
			&taskCreatorName,
		); err != nil {
			return nil, err
		}

		goal.OwnerName = goalOwnerName

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
				ID:            int(taskID.Int64),
				GoalID:        int(taskGoalID.Int64),
				GoalTitle:     current.Title,
				Title:         taskTitle.String,
				Description:   taskDesc.String,
				Priority:      normalizePriority(taskPriority.String),
				IsCompleted:   taskIsCompleted.Valid && taskIsCompleted.Bool,
				CreatedBy:     int(createdBy.Int64),
				CreatedByName: taskCreatorName.String,
				CreatedAt:     taskAt.Time,
			}
			if assigneeID.Valid {
				value := int(assigneeID.Int64)
				task.AssigneeID = &value
			}
			if taskAssigneeName.Valid {
				task.AssigneeName = taskAssigneeName.String
			}
			current.Tasks = append(current.Tasks, task)
		}
	}

	return goals, rows.Err()
}

func (s *Store) GetGoalWithTasks(goalID, _ int) (*types.GoalWithTasks, error) {
	rows, err := s.db.Query(
		`SELECT
			g.id,
			g.title,
			g.description,
			g.priority,
			g.status,
			g.owner_id,
			g.created_at,
			TRIM(CONCAT(owner_u.first_name, ' ', owner_u.last_name)) AS owner_name,
			t.id,
			t.goal_id,
			t.title,
			t.description,
			t.priority,
			t.is_completed,
			t.assignee_id,
			t.created_by,
			t.created_at,
			TRIM(CONCAT(assignee_u.first_name, ' ', assignee_u.last_name)) AS assignee_name,
			TRIM(CONCAT(creator_u.first_name, ' ', creator_u.last_name)) AS creator_name
		FROM goals g
		JOIN users owner_u ON owner_u.id = g.owner_id
		LEFT JOIN tasks t ON t.goal_id = g.id
		LEFT JOIN users assignee_u ON assignee_u.id = t.assignee_id
		LEFT JOIN users creator_u ON creator_u.id = t.created_by
		WHERE g.id = $1
		ORDER BY
			CASE WHEN t.is_completed THEN 1 ELSE 0 END,
			CASE t.priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 WHEN 'low' THEN 2 ELSE 3 END,
			t.created_at ASC`,
		goalID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		goalFound bool
		goal      types.Goal
		goalModel *types.GoalWithTasks
	)

	for rows.Next() {
		var (
			currentGoal      types.Goal
			goalOwnerName    string
			taskID           sql.NullInt64
			taskGoalID       sql.NullInt64
			taskTitle        sql.NullString
			taskDesc         sql.NullString
			taskPriority     sql.NullString
			taskIsCompleted  sql.NullBool
			assigneeID       sql.NullInt64
			createdBy        sql.NullInt64
			taskAt           sql.NullTime
			taskAssigneeName sql.NullString
			taskCreatorName  sql.NullString
		)

		if err := rows.Scan(
			&currentGoal.ID,
			&currentGoal.Title,
			&currentGoal.Description,
			&currentGoal.Priority,
			&currentGoal.Status,
			&currentGoal.OwnerID,
			&currentGoal.CreatedAt,
			&goalOwnerName,
			&taskID,
			&taskGoalID,
			&taskTitle,
			&taskDesc,
			&taskPriority,
			&taskIsCompleted,
			&assigneeID,
			&createdBy,
			&taskAt,
			&taskAssigneeName,
			&taskCreatorName,
		); err != nil {
			return nil, err
		}

		if !goalFound {
			currentGoal.OwnerName = goalOwnerName
			goal = currentGoal
			goalModel = &types.GoalWithTasks{
				Goal:  goal,
				Tasks: []*types.Task{},
			}
			goalFound = true
		}

		if taskID.Valid {
			task := &types.Task{
				ID:            int(taskID.Int64),
				GoalID:        int(taskGoalID.Int64),
				GoalTitle:     goal.Title,
				Title:         taskTitle.String,
				Description:   taskDesc.String,
				Priority:      normalizePriority(taskPriority.String),
				IsCompleted:   taskIsCompleted.Valid && taskIsCompleted.Bool,
				CreatedBy:     int(createdBy.Int64),
				CreatedByName: taskCreatorName.String,
				CreatedAt:     taskAt.Time,
			}
			if assigneeID.Valid {
				value := int(assigneeID.Int64)
				task.AssigneeID = &value
			}
			if taskAssigneeName.Valid {
				task.AssigneeName = taskAssigneeName.String
			}
			goalModel.Tasks = append(goalModel.Tasks, task)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !goalFound {
		return nil, ErrForbidden
	}
	return goalModel, nil
}

func (s *Store) GetUsersWithCurrentTasks() ([]*types.UserTasksBoard, error) {
	rows, err := s.db.Query(
		`SELECT
			u.id,
			TRIM(CONCAT(u.first_name, ' ', u.last_name)) AS user_name,
			u.email,
			t.id,
			t.goal_id,
			t.title,
			t.description,
			t.priority,
			t.is_completed,
			t.assignee_id,
			t.created_by,
			t.created_at,
			g.title AS goal_title,
			TRIM(CONCAT(assignee_u.first_name, ' ', assignee_u.last_name)) AS assignee_name,
			TRIM(CONCAT(creator_u.first_name, ' ', creator_u.last_name)) AS creator_name
		FROM users u
		LEFT JOIN tasks t
			ON t.assignee_id = u.id
			AND t.is_completed = FALSE
		LEFT JOIN goals g ON g.id = t.goal_id
		LEFT JOIN users assignee_u ON assignee_u.id = t.assignee_id
		LEFT JOIN users creator_u ON creator_u.id = t.created_by
		ORDER BY
			u.first_name,
			u.last_name,
			u.id,
			CASE t.priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 WHEN 'low' THEN 2 ELSE 3 END,
			t.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*types.UserTasksBoard, 0)
	boardByUserID := make(map[int]*types.UserTasksBoard)

	for rows.Next() {
		var (
			userID    int
			userName  string
			userEmail string

			taskID           sql.NullInt64
			taskGoalID       sql.NullInt64
			taskTitle        sql.NullString
			taskDesc         sql.NullString
			taskPriority     sql.NullString
			taskIsCompleted  sql.NullBool
			assigneeID       sql.NullInt64
			createdBy        sql.NullInt64
			taskAt           sql.NullTime
			taskGoalTitle    sql.NullString
			taskAssigneeName sql.NullString
			taskCreatorName  sql.NullString
		)

		if err := rows.Scan(
			&userID,
			&userName,
			&userEmail,
			&taskID,
			&taskGoalID,
			&taskTitle,
			&taskDesc,
			&taskPriority,
			&taskIsCompleted,
			&assigneeID,
			&createdBy,
			&taskAt,
			&taskGoalTitle,
			&taskAssigneeName,
			&taskCreatorName,
		); err != nil {
			return nil, err
		}

		board, exists := boardByUserID[userID]
		if !exists {
			board = &types.UserTasksBoard{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
				Tasks: []*types.Task{},
			}
			boardByUserID[userID] = board
			boards = append(boards, board)
		}

		if taskID.Valid {
			task := &types.Task{
				ID:          int(taskID.Int64),
				GoalID:      int(taskGoalID.Int64),
				GoalTitle:   taskGoalTitle.String,
				Title:       taskTitle.String,
				Description: taskDesc.String,
				Priority:    normalizePriority(taskPriority.String),
				IsCompleted: taskIsCompleted.Valid && taskIsCompleted.Bool,
				CreatedBy:   int(createdBy.Int64),
				CreatedAt:   taskAt.Time,
			}
			if assigneeID.Valid {
				value := int(assigneeID.Int64)
				task.AssigneeID = &value
			}
			if taskAssigneeName.Valid {
				task.AssigneeName = taskAssigneeName.String
			}
			if taskCreatorName.Valid {
				task.CreatedByName = taskCreatorName.String
			}
			board.Tasks = append(board.Tasks, task)
		}
	}

	return boards, rows.Err()
}

func (s *Store) CreateTask(goalID, creatorID int, payload types.CreateTaskPayload) (*types.Task, error) {
	row := s.db.QueryRow(
		`INSERT INTO tasks (goal_id, title, description, priority, assignee_id, created_by)
		 SELECT g.id, $2, $3, $4, $5, $6
		 FROM goals g
		 WHERE g.id = $1 AND g.owner_id = $6
		 RETURNING id, goal_id, title, description, priority, is_completed, assignee_id, created_by, created_at`,
		goalID,
		payload.Title,
		payload.Description,
		normalizePriority(payload.Priority),
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

func (s *Store) UpdateTask(taskID, requesterID int, payload types.UpdateTaskPayload) (*types.Task, error) {
	row := s.db.QueryRow(
		`UPDATE tasks t
		 SET goal_id = $1,
		     title = $2,
		     description = $3,
		     priority = $4,
		     is_completed = $5,
		     assignee_id = $6
		 FROM goals current_goal, goals new_goal
		 WHERE t.id = $7
		   AND t.goal_id = current_goal.id
		   AND current_goal.owner_id = $8
		   AND new_goal.id = $1
		   AND new_goal.owner_id = $8
		 RETURNING t.id, t.goal_id, t.title, t.description, t.priority, t.is_completed, t.assignee_id, t.created_by, t.created_at`,
		payload.GoalID,
		payload.Title,
		payload.Description,
		normalizePriority(payload.Priority),
		payload.IsCompleted,
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

func (s *Store) DeleteTask(taskID, requesterID int) error {
	result, err := s.db.Exec(
		`DELETE FROM tasks t
		 USING goals g
		 WHERE t.goal_id = g.id
		   AND t.id = $1
		   AND g.owner_id = $2`,
		taskID,
		requesterID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrForbidden
	}
	return nil
}

func (s *Store) AssignTask(taskID, requesterID int, payload types.AssignTaskPayload) (*types.Task, error) {
	row := s.db.QueryRow(
		`UPDATE tasks t
		 SET assignee_id = $1
		 FROM goals g
		 WHERE t.goal_id = g.id AND t.id = $2 AND g.owner_id = $3
		 RETURNING t.id, t.goal_id, t.title, t.description, t.priority, t.is_completed, t.assignee_id, t.created_by, t.created_at`,
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
		`SELECT
			t.id,
			t.goal_id,
			t.title,
			t.description,
			t.priority,
			t.is_completed,
			t.assignee_id,
			t.created_by,
			t.created_at,
			g.title AS goal_title,
			TRIM(CONCAT(assignee_u.first_name, ' ', assignee_u.last_name)) AS assignee_name,
			TRIM(CONCAT(creator_u.first_name, ' ', creator_u.last_name)) AS creator_name
		 FROM tasks t
		 JOIN goals g ON g.id = t.goal_id
		 LEFT JOIN users assignee_u ON assignee_u.id = t.assignee_id
		 LEFT JOIN users creator_u ON creator_u.id = t.created_by
		 WHERE t.assignee_id = $1
		 ORDER BY
			CASE WHEN t.is_completed THEN 1 ELSE 0 END,
			CASE t.priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
			t.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task, err := scanRowIntoTaskWithLookups(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *Store) ListUsers() ([]*types.UserLookup, error) {
	rows, err := s.db.Query(
		`SELECT id, TRIM(CONCAT(first_name, ' ', last_name)) AS full_name
		 FROM users
		 ORDER BY first_name, last_name, id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*types.UserLookup, 0)
	for rows.Next() {
		user := new(types.UserLookup)
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRowIntoGoal(row rowScanner) (*types.Goal, error) {
	g := new(types.Goal)
	if err := row.Scan(&g.ID, &g.Title, &g.Description, &g.Priority, &g.Status, &g.OwnerID, &g.CreatedAt); err != nil {
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
		&task.Priority,
		&task.IsCompleted,
		&assigneeID,
		&task.CreatedBy,
		&task.CreatedAt,
	); err != nil {
		return nil, err
	}
	task.Priority = normalizePriority(task.Priority)
	if assigneeID.Valid {
		value := int(assigneeID.Int64)
		task.AssigneeID = &value
	}
	return task, nil
}

func scanRowIntoTaskWithLookups(row rowScanner) (*types.Task, error) {
	task := new(types.Task)
	var assigneeID sql.NullInt64
	var assigneeName sql.NullString
	var creatorName sql.NullString
	if err := row.Scan(
		&task.ID,
		&task.GoalID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.IsCompleted,
		&assigneeID,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.GoalTitle,
		&assigneeName,
		&creatorName,
	); err != nil {
		return nil, err
	}
	task.Priority = normalizePriority(task.Priority)
	if assigneeID.Valid {
		value := int(assigneeID.Int64)
		task.AssigneeID = &value
	}
	if assigneeName.Valid {
		task.AssigneeName = assigneeName.String
	}
	if creatorName.Valid {
		task.CreatedByName = creatorName.String
	}
	return task, nil
}

func normalizePriority(priority string) string {
	switch priority {
	case "high", "medium", "low":
		return priority
	default:
		return "medium"
	}
}

func normalizeGoalStatus(status string) string {
	switch status {
	case "todo", "in_progress", "achieved":
		return status
	default:
		return "todo"
	}
}
