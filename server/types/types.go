package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	UpdateUserProfile(userID int, payload UpdateProfilePayload) (*User, error)
	UpdateUserPassword(userID int, hashedPassword string) error
	ListUsers() ([]*UserLookup, error)
}

type GoalTaskStore interface {
	CreateGoal(ownerID int, payload CreateGoalPayload) (*Goal, error)
	UpdateGoal(goalID, ownerID int, payload CreateGoalPayload) (*Goal, error)
	DeleteGoal(goalID, ownerID int) error
	GetGoalsByOwner(ownerID int) ([]*GoalWithTasks, error)
	GetGoalWithTasks(goalID, ownerID int) (*GoalWithTasks, error)
	GetUsersWithCurrentTasks() ([]*UserTasksBoard, error)
	CreateTask(goalID, creatorID int, payload CreateTaskPayload) (*Task, error)
	UpdateTask(taskID, requesterID int, payload UpdateTaskPayload) (*Task, error)
	DeleteTask(taskID, requesterID int) error
	AssignTask(taskID, requesterID int, payload AssignTaskPayload) (*Task, error)
	GetAssignedTasks(userID int) ([]*Task, error)
	ListUsers() ([]*UserLookup, error)
}

type Goal struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
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
	Description string `json:"description" validate:"max=2000"`
	Priority    string `json:"priority" validate:"required,oneof=high medium low"`
	Status      string `json:"status" validate:"required,oneof=todo in_progress achieved"`
}

type Task struct {
	ID            int       `json:"id"`
	GoalID        int       `json:"goalId"`
	GoalTitle     string    `json:"goalTitle,omitempty"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Priority      string    `json:"priority"`
	IsCompleted   bool      `json:"isCompleted"`
	AssigneeID    *int      `json:"assigneeId,omitempty"`
	AssigneeName  string    `json:"assigneeName,omitempty"`
	CreatedBy     int       `json:"createdBy"`
	CreatedByName string    `json:"createdByName,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateTaskPayload struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=2000"`
	Priority    string `json:"priority" validate:"required,oneof=high medium low"`
	AssigneeID  *int   `json:"assigneeId,omitempty"`
}

type UpdateTaskPayload struct {
	GoalID      int    `json:"goalId" validate:"required,min=1"`
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=2000"`
	Priority    string `json:"priority" validate:"required,oneof=high medium low"`
	IsCompleted bool   `json:"isCompleted"`
	AssigneeID  *int   `json:"assigneeId,omitempty"`
}

type AssignTaskPayload struct {
	AssigneeID *int `json:"assigneeId"`
}

type UserLookup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserTasksBoard struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Tasks []*Task `json:"tasks"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserProfile struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateProfilePayload struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=255"`
	LastName  string `json:"lastName" validate:"required,min=1,max=255"`
}

type UpdatePasswordPayload struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=3,max=130"`
	NewPassword     string `json:"newPassword" validate:"required,min=3,max=130"`
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
