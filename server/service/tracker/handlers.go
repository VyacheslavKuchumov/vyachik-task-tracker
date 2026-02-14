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

// HandleUpdateGoal godoc
// @Summary Update goal
// @Description Update a goal owned by the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param goalID path int true "Goal ID"
// @Param payload body types.CreateGoalPayload true "Goal payload"
// @Success 200 {object} types.Goal
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals/{goalID} [put]
func (h *Handler) HandleUpdateGoal(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if ownerID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	goalID, err := parsePathID(r, "goalID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid goal id"))
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

	goal, err := h.store.UpdateGoal(goalID, ownerID, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, goal)
}

// HandleDeleteGoal godoc
// @Summary Delete goal
// @Description Delete a goal owned by the authenticated user
// @Tags goals
// @Produce json
// @Security BearerAuth
// @Param goalID path int true "Goal ID"
// @Success 204 {object} nil
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals/{goalID} [delete]
func (h *Handler) HandleDeleteGoal(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if ownerID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	goalID, err := parsePathID(r, "goalID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid goal id"))
		return
	}

	if err := h.store.DeleteGoal(goalID, ownerID); err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetGoals godoc
// @Summary Get goals
// @Description Get all goals with nested tasks for authenticated users
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

// HandleGetGoalTasks godoc
// @Summary Get tasks by goal
// @Description Get a single goal with its tasks for authenticated users
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param goalID path int true "Goal ID"
// @Success 200 {object} types.GoalWithTasks
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /goals/{goalID}/tasks [get]
func (h *Handler) HandleGetGoalTasks(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if ownerID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	goalID, err := parsePathID(r, "goalID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid goal id"))
		return
	}

	goalWithTasks, err := h.store.GetGoalWithTasks(goalID, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, goalWithTasks)
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

// HandleUpdateTask godoc
// @Summary Update task
// @Description Update a task under a goal owned by the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path int true "Task ID"
// @Param payload body types.UpdateTaskPayload true "Task update payload"
// @Success 200 {object} types.Task
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /tasks/{taskID} [put]
func (h *Handler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	requesterID := auth.GetUserIDFromContext(r.Context())
	if requesterID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	taskID, err := parsePathID(r, "taskID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task id"))
		return
	}

	var payload types.UpdateTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	task, err := h.store.UpdateTask(taskID, requesterID, payload)
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

// HandleDeleteTask godoc
// @Summary Delete task
// @Description Delete a task from a goal owned by the authenticated user
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param taskID path int true "Task ID"
// @Success 204 {object} nil
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /tasks/{taskID} [delete]
func (h *Handler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	requesterID := auth.GetUserIDFromContext(r.Context())
	if requesterID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	taskID, err := parsePathID(r, "taskID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task id"))
		return
	}

	if err := h.store.DeleteTask(taskID, requesterID); err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

// HandleGetUsersWithCurrentTasks godoc
// @Summary Get users with current tasks
// @Description Get all users and their current assigned tasks (not completed)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} types.UserTasksBoard
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /users/tasks [get]
func (h *Handler) HandleGetUsersWithCurrentTasks(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	usersTasks, err := h.store.GetUsersWithCurrentTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, usersTasks)
}

func parsePathID(r *http.Request, key string) (int, error) {
	value := chi.URLParam(r, key)
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return id, nil
}
