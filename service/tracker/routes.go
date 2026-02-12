package tracker

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.GoalTaskStore
}

func NewHandler(store types.GoalTaskStore) *Handler {
	return &Handler{store: store}
}

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

func (h *Handler) HandleDashboard(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardPage))
}

func (h *Handler) HandleHTMXGoals(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderGoals(goals)))
}

func (h *Handler) HandleHTMXAssignedTasks(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	tasks, err := h.store.GetAssignedTasks(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderAssignedTasks(tasks)))
}

func (h *Handler) HandleHTMXCreateGoal(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	payload := types.CreateGoalPayload{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		http.Error(w, "invalid goal payload", http.StatusBadRequest)
		return
	}

	if _, err := h.store.CreateGoal(ownerID, payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.HandleHTMXGoals(w, r)
}

func (h *Handler) HandleHTMXCreateTask(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	goalID, err := strconv.Atoi(r.FormValue("goalId"))
	if err != nil || goalID <= 0 {
		http.Error(w, "invalid goal id", http.StatusBadRequest)
		return
	}

	var assigneeID *int
	if raw := strings.TrimSpace(r.FormValue("assigneeId")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			http.Error(w, "invalid assignee id", http.StatusBadRequest)
			return
		}
		assigneeID = &value
	}

	payload := types.CreateTaskPayload{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		AssigneeID:  assigneeID,
	}
	if err := utils.Validate.Struct(payload); err != nil {
		http.Error(w, "invalid task payload", http.StatusBadRequest)
		return
	}

	if _, err := h.store.CreateTask(goalID, ownerID, payload); err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	h.HandleHTMXGoals(w, r)
}

func (h *Handler) HandleHTMXAssignTask(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(r.FormValue("taskId"))
	if err != nil || taskID <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var assigneeID *int
	if raw := strings.TrimSpace(r.FormValue("assigneeId")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			http.Error(w, "invalid assignee id", http.StatusBadRequest)
			return
		}
		assigneeID = &value
	}

	if _, err := h.store.AssignTask(taskID, userID, types.AssignTaskPayload{AssigneeID: assigneeID}); err != nil {
		status := http.StatusInternalServerError
		if err == ErrForbidden {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	h.HandleHTMXGoals(w, r)
}

func renderGoals(goals []*types.GoalWithTasks) string {
	if len(goals) == 0 {
		return `<div class="empty">No goals yet.</div>`
	}

	var b strings.Builder
	for _, goal := range goals {
		b.WriteString(`<article class="card">`)
		b.WriteString(`<h3>`)
		b.WriteString(html.EscapeString(goal.Title))
		b.WriteString(`</h3><p>`)
		b.WriteString(html.EscapeString(goal.Description))
		b.WriteString(`</p><small>Goal ID: `)
		b.WriteString(strconv.Itoa(goal.ID))
		b.WriteString(`</small><div class="tasks">`)
		if len(goal.Tasks) == 0 {
			b.WriteString(`<div class="empty">No tasks yet.</div>`)
		} else {
			for _, task := range goal.Tasks {
				b.WriteString(`<div class="task"><strong>`)
				b.WriteString(html.EscapeString(task.Title))
				b.WriteString(`</strong><p>`)
				b.WriteString(html.EscapeString(task.Description))
				b.WriteString(`</p><small>Task ID: `)
				b.WriteString(strconv.Itoa(task.ID))
				b.WriteString(` | Status: `)
				b.WriteString(html.EscapeString(task.Status))
				b.WriteString(` | Assignee: `)
				if task.AssigneeID != nil {
					b.WriteString(strconv.Itoa(*task.AssigneeID))
				} else {
					b.WriteString(`none`)
				}
				b.WriteString(`</small></div>`)
			}
		}
		b.WriteString(`</div></article>`)
	}

	return b.String()
}

func renderAssignedTasks(tasks []*types.Task) string {
	if len(tasks) == 0 {
		return `<div class="empty">No tasks assigned to you.</div>`
	}

	var b strings.Builder
	for _, task := range tasks {
		b.WriteString(`<div class="task"><strong>`)
		b.WriteString(html.EscapeString(task.Title))
		b.WriteString(`</strong><p>`)
		b.WriteString(html.EscapeString(task.Description))
		b.WriteString(`</p><small>Task ID: `)
		b.WriteString(strconv.Itoa(task.ID))
		b.WriteString(` | Goal ID: `)
		b.WriteString(strconv.Itoa(task.GoalID))
		b.WriteString(`</small></div>`)
	}

	return b.String()
}

const dashboardPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Task Tracker</title>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
  <style>
    :root { --bg:#f5f8ff; --ink:#0f172a; --accent:#0ea5a8; --card:#ffffff; --line:#cbd5e1; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: "Trebuchet MS", "Segoe UI", sans-serif; color:var(--ink); background:linear-gradient(145deg,#eef4ff,#f8fff8); }
    .wrap { max-width:1100px; margin:0 auto; padding:20px; }
    .grid { display:grid; grid-template-columns:1fr 1fr; gap:16px; }
    .panel { background:var(--card); border:1px solid var(--line); border-radius:14px; padding:14px; box-shadow:0 8px 24px rgba(15,23,42,.06); }
    .panel h2, .panel h3 { margin:.2rem 0 .8rem; }
    input, textarea, button { width:100%; margin:.3rem 0; padding:.6rem .7rem; border-radius:10px; border:1px solid var(--line); }
    button { background:var(--accent); color:#fff; border:none; cursor:pointer; font-weight:700; }
    .card, .task { border:1px solid var(--line); border-radius:12px; background:#fff; padding:10px; margin:10px 0; }
    .empty { color:#475569; font-style:italic; padding:8px 0; }
    @media (max-width:900px){ .grid { grid-template-columns:1fr; } }
  </style>
</head>
<body>
  <div class="wrap">
    <h1>Goal-Based Task Tracker</h1>
    <div class="grid">
      <section class="panel">
        <h2>Auth</h2>
        <input id="token" placeholder="JWT token appears here after login">
        <form id="registerForm">
          <h3>Register</h3>
          <input name="firstName" placeholder="First name" required>
          <input name="lastName" placeholder="Last name" required>
          <input name="email" type="email" placeholder="Email" required>
          <input name="password" type="password" placeholder="Password" required>
          <button type="submit">Register</button>
        </form>
        <form id="loginForm">
          <h3>Login</h3>
          <input name="email" type="email" placeholder="Email" required>
          <input name="password" type="password" placeholder="Password" required>
          <button type="submit">Login</button>
        </form>
      </section>

      <section class="panel">
        <h2>Create Goal</h2>
        <form hx-post="/htmx/goals/create" hx-target="#goalsList" hx-swap="innerHTML">
          <input name="title" placeholder="Big goal title" required>
          <textarea name="description" placeholder="Goal description" required></textarea>
          <button type="submit">Create Goal</button>
        </form>

        <h2>Create Task</h2>
        <form hx-post="/htmx/tasks/create" hx-target="#goalsList" hx-swap="innerHTML">
          <input name="goalId" type="number" placeholder="Goal ID" required>
          <input name="title" placeholder="Task title" required>
          <textarea name="description" placeholder="Task description" required></textarea>
          <input name="assigneeId" type="number" placeholder="Assignee user ID (optional)">
          <button type="submit">Create Task</button>
        </form>

        <h2>Assign Task</h2>
        <form hx-post="/htmx/tasks/assign" hx-target="#goalsList" hx-swap="innerHTML">
          <input name="taskId" type="number" placeholder="Task ID" required>
          <input name="assigneeId" type="number" placeholder="Assignee user ID (empty to unassign)">
          <button type="submit">Assign</button>
        </form>
      </section>
    </div>

    <section class="panel">
      <h2>Your Goals</h2>
      <button hx-get="/htmx/goals" hx-target="#goalsList" hx-swap="innerHTML">Refresh Goals</button>
      <div id="goalsList"></div>
    </section>

    <section class="panel">
      <h2>Tasks Assigned To You</h2>
      <button hx-get="/htmx/tasks/assigned" hx-target="#assignedList" hx-swap="innerHTML">Refresh Assigned</button>
      <div id="assignedList"></div>
    </section>
  </div>

  <script>
    const tokenInput = document.getElementById('token');

    document.body.addEventListener('htmx:configRequest', function(evt) {
      const token = tokenInput.value.trim();
      if (token) {
        evt.detail.headers['Authorization'] = token;
      }
    });

    document.getElementById('registerForm').addEventListener('submit', async function(evt) {
      evt.preventDefault();
      const formData = new FormData(evt.target);
      const payload = Object.fromEntries(formData.entries());
      await fetch('/api/v1/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
    });

    document.getElementById('loginForm').addEventListener('submit', async function(evt) {
      evt.preventDefault();
      const formData = new FormData(evt.target);
      const payload = Object.fromEntries(formData.entries());
      const res = await fetch('/api/v1/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      const data = await res.json().catch(() => ({}));
      if (data.token) tokenInput.value = data.token;
    });
  </script>
</body>
</html>`
