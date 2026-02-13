package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type GoalTaskStore interface {
	CreateGoal(ownerID int, payload CreateGoalPayload) (*Goal, error)
	UpdateGoal(goalID, ownerID int, payload CreateGoalPayload) (*Goal, error)
	GetGoalsByOwner(ownerID int) ([]*GoalWithTasks, error)
	CreateTask(goalID, creatorID int, payload CreateTaskPayload) (*Task, error)
	UpdateTask(taskID, requesterID int, payload UpdateTaskPayload) (*Task, error)
	AssignTask(taskID, requesterID int, payload AssignTaskPayload) (*Task, error)
	GetAssignedTasks(userID int) ([]*Task, error)
	ListUsers() ([]*UserLookup, error)
}

type Goal struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     int       `json:"ownerId"`
	OwnerName   string    `json:"ownerName,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GoalWithTasks struct {
	Goal
	Tasks []*Task `json:"tasks"`
}

type CreateGoalPayload struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=3,max=2000"`
}

type Task struct {
	ID            int       `json:"id"`
	GoalID        int       `json:"goalId"`
	GoalTitle     string    `json:"goalTitle,omitempty"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	AssigneeID    *int      `json:"assigneeId,omitempty"`
	AssigneeName  string    `json:"assigneeName,omitempty"`
	CreatedBy     int       `json:"createdBy"`
	CreatedByName string    `json:"createdByName,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateTaskPayload struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=3,max=2000"`
	AssigneeID  *int   `json:"assigneeId,omitempty"`
}

type UpdateTaskPayload struct {
	GoalID      int    `json:"goalId" validate:"required,min=1"`
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=3,max=2000"`
	Status      string `json:"status" validate:"required,oneof=todo in_progress done"`
	AssigneeID  *int   `json:"assigneeId,omitempty"`
}

type AssignTaskPayload struct {
	AssigneeID *int `json:"assigneeId"`
}

type UserLookup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
