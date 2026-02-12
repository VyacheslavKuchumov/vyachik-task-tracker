package tracker

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestTrackerHandlers(t *testing.T) {
	store := &mockGoalTaskStore{}
	handler := NewHandler(store)

	t.Run("create goal rejects invalid payload", func(t *testing.T) {
		payload := types.CreateGoalPayload{
			Title:       "x",
			Description: "short",
		}
		body, _ := json.Marshal(payload)
		req := newRequestWithUser(http.MethodPost, "/api/v1/goals", body, 1)
		rr := httptest.NewRecorder()

		handler.HandleCreateGoal(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("create goal returns created goal", func(t *testing.T) {
		payload := types.CreateGoalPayload{
			Title:       "Launch MVP",
			Description: "Ship first version",
		}
		body, _ := json.Marshal(payload)
		req := newRequestWithUser(http.MethodPost, "/api/v1/goals", body, 2)
		rr := httptest.NewRecorder()

		handler.HandleCreateGoal(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("create task validates goal path param", func(t *testing.T) {
		payload := types.CreateTaskPayload{
			Title:       "Break down tasks",
			Description: "Create subtasks",
		}
		body, _ := json.Marshal(payload)
		req := newRequestWithUser(http.MethodPost, "/api/v1/goals/wrong/tasks", body, 2)
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("goalID", "wrong")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.HandleCreateTask(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("assign task maps forbidden error", func(t *testing.T) {
		store.assignErr = ErrForbidden
		defer func() { store.assignErr = nil }()

		payload := types.AssignTaskPayload{AssigneeID: intPtr(3)}
		body, _ := json.Marshal(payload)
		req := newRequestWithUser(http.MethodPut, "/api/v1/tasks/10/assign", body, 2)
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("taskID", "10")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.HandleAssignTask(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("assigned tasks returns ok", func(t *testing.T) {
		req := newRequestWithUser(http.MethodGet, "/api/v1/tasks/assigned", nil, 4)
		rr := httptest.NewRecorder()

		handler.HandleGetAssignedTasks(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockGoalTaskStore struct {
	assignErr error
}

func (m *mockGoalTaskStore) CreateGoal(ownerID int, payload types.CreateGoalPayload) (*types.Goal, error) {
	return &types.Goal{
		ID:          1,
		Title:       payload.Title,
		Description: payload.Description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *mockGoalTaskStore) GetGoalsByOwner(ownerID int) ([]*types.GoalWithTasks, error) {
	return []*types.GoalWithTasks{
		{
			Goal: types.Goal{
				ID:          1,
				Title:       "Goal",
				Description: "Desc",
				OwnerID:     ownerID,
				CreatedAt:   time.Now(),
			},
			Tasks: []*types.Task{},
		},
	}, nil
}

func (m *mockGoalTaskStore) UpdateGoal(goalID, ownerID int, payload types.CreateGoalPayload) (*types.Goal, error) {
	return &types.Goal{
		ID:          goalID,
		Title:       payload.Title,
		Description: payload.Description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *mockGoalTaskStore) CreateTask(goalID, creatorID int, payload types.CreateTaskPayload) (*types.Task, error) {
	return &types.Task{
		ID:          1,
		GoalID:      goalID,
		Title:       payload.Title,
		Description: payload.Description,
		Status:      "todo",
		AssigneeID:  payload.AssigneeID,
		CreatedBy:   creatorID,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *mockGoalTaskStore) UpdateTask(taskID, requesterID int, payload types.UpdateTaskPayload) (*types.Task, error) {
	return &types.Task{
		ID:          taskID,
		GoalID:      payload.GoalID,
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
		AssigneeID:  payload.AssigneeID,
		CreatedBy:   requesterID,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *mockGoalTaskStore) AssignTask(taskID, requesterID int, payload types.AssignTaskPayload) (*types.Task, error) {
	if m.assignErr != nil {
		return nil, m.assignErr
	}

	return &types.Task{
		ID:         taskID,
		GoalID:     1,
		Title:      "task",
		Status:     "todo",
		AssigneeID: payload.AssigneeID,
		CreatedBy:  requesterID,
		CreatedAt:  time.Now(),
	}, nil
}

func (m *mockGoalTaskStore) GetAssignedTasks(userID int) ([]*types.Task, error) {
	return []*types.Task{
		{
			ID:         1,
			GoalID:     1,
			Title:      "Task A",
			Status:     "todo",
			AssigneeID: &userID,
			CreatedBy:  2,
			CreatedAt:  time.Now(),
		},
	}, nil
}

func (m *mockGoalTaskStore) ListUsers() ([]*types.UserLookup, error) {
	return []*types.UserLookup{
		{ID: 1, Name: "Alice Doe"},
		{ID: 2, Name: "Bob Doe"},
	}, nil
}

func newRequestWithUser(method, path string, payload []byte, userID int) *http.Request {
	var body *bytes.Buffer
	if payload != nil {
		body = bytes.NewBuffer(payload)
	} else {
		body = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, body)
	ctx := context.WithValue(req.Context(), auth.UserKey, userID)
	return req.WithContext(ctx)
}

func intPtr(v int) *int { return &v }
