package tracker

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.GoalTaskStore
}

func NewHandler(store types.GoalTaskStore) *Handler {
	return &Handler{store: store}
}

// HandleCreateGoal godoc
// @Summary Create goal
// @Description Create a goal owned by the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body types.CreateGoalPayload true "Goal payload"
// @Success 201 {object} types.Goal
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals [post]
func (h *Handler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if ownerID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	var payload types.CreateGoalPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	goal, err := h.store.CreateGoal(ownerID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, goal)
}

// HandleGetGoals godoc
// @Summary Get goals
// @Description Get all goals for the authenticated user with nested tasks
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Success 200 {array} types.GoalWithTasks
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals [get]
func (h *Handler) HandleGetGoals(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if ownerID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, goals)
}

// HandleCreateTask godoc
// @Summary Create task
// @Description Create a task under a goal. Only the goal owner can create tasks.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goalID path int true "Goal ID"
// @Param payload body types.CreateTaskPayload true "Task payload"
// @Success 201 {object} types.Task
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals/{goalID}/tasks [post]
func (h *Handler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	creatorID := auth.GetUserIDFromContext(r.Context())
	if creatorID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	goalID, err := strconv.Atoi(chi.URLParam(r, "goalID"))
	if err != nil || goalID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid goal id"))
		return
	}

	var payload types.CreateTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	task, err := h.store.CreateTask(goalID, creatorID, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, task)
}

// HandleAssignTask godoc
// @Summary Assign task
// @Description Assign or unassign a task. Only the goal owner can assign tasks.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path int true "Task ID"
// @Param payload body types.AssignTaskPayload true "Assignment payload"
// @Success 200 {object} types.Task
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /tasks/{taskID}/assign [put]
func (h *Handler) HandleAssignTask(w http.ResponseWriter, r *http.Request) {
	requesterID := auth.GetUserIDFromContext(r.Context())
	if requesterID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil || taskID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task id"))
		return
	}

	var payload types.AssignTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	task, err := h.store.AssignTask(taskID, requesterID, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

// HandleGetAssignedTasks godoc
// @Summary Get assigned tasks
// @Description Get tasks assigned to the authenticated user
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} types.Task
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /tasks/assigned [get]
func (h *Handler) HandleGetAssignedTasks(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	tasks, err := h.store.GetAssignedTasks(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}
