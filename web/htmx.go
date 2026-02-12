package web

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/service/tracker"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleHTMXGoals(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	goals = filterGoals(goals, query, status)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderGoalsTable(goals)))
}

func (h *Handler) HandleHTMXTasks(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := flattenOwnedTasks(goals)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	tasks = filterTasks(tasks, query, status)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderTasksTable(tasks)))
}

func (h *Handler) HandleHTMXGoalCard(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var goal *types.Goal
	goalIDRaw := strings.TrimSpace(chi.URLParam(r, "goalID"))
	if goalIDRaw != "" {
		goalID, err := strconv.Atoi(goalIDRaw)
		if err != nil || goalID <= 0 {
			http.Error(w, "invalid goal id", http.StatusBadRequest)
			return
		}
		goal = findGoal(goals, goalID)
		if goal == nil {
			http.Error(w, "goal not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderGoalCard(goal)))
}

func (h *Handler) HandleHTMXTaskCard(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	goals, err := h.store.GetGoalsByOwner(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users, err := h.store.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task := &types.Task{Status: "todo"}
	if len(goals) > 0 {
		task.GoalID = goals[0].ID
	}

	taskIDRaw := strings.TrimSpace(chi.URLParam(r, "taskID"))
	if taskIDRaw != "" {
		taskID, err := strconv.Atoi(taskIDRaw)
		if err != nil || taskID <= 0 {
			http.Error(w, "invalid task id", http.StatusBadRequest)
			return
		}
		found := findTask(flattenOwnedTasks(goals), taskID)
		if found == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		task = found
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(renderTaskCard(task, goals, users)))
}

func (h *Handler) HandleHTMXGoalSave(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	payload := types.CreateGoalPayload{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}
	if err := utils.Validate.Struct(payload); err != nil {
		http.Error(w, "invalid goal payload", http.StatusBadRequest)
		return
	}

	goalIDRaw := strings.TrimSpace(r.FormValue("goalId"))
	if goalIDRaw == "" {
		if _, err := h.store.CreateGoal(ownerID, payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		goalID, err := strconv.Atoi(goalIDRaw)
		if err != nil || goalID <= 0 {
			http.Error(w, "invalid goal id", http.StatusBadRequest)
			return
		}
		if _, err := h.store.UpdateGoal(goalID, ownerID, payload); err != nil {
			status := http.StatusInternalServerError
			if err == tracker.ErrForbidden {
				status = http.StatusForbidden
			}
			http.Error(w, err.Error(), status)
			return
		}
	}

	h.HandleHTMXGoals(w, r)
}

func (h *Handler) HandleHTMXTaskSave(w http.ResponseWriter, r *http.Request) {
	ownerID := auth.GetUserIDFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	goalID, err := strconv.Atoi(strings.TrimSpace(r.FormValue("goalId")))
	if err != nil || goalID <= 0 {
		http.Error(w, "invalid goal", http.StatusBadRequest)
		return
	}

	var assigneeID *int
	if raw := strings.TrimSpace(r.FormValue("assigneeId")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			http.Error(w, "invalid assignee", http.StatusBadRequest)
			return
		}
		assigneeID = &value
	}

	taskIDRaw := strings.TrimSpace(r.FormValue("taskId"))
	if taskIDRaw == "" {
		payload := types.CreateTaskPayload{
			Title:       strings.TrimSpace(r.FormValue("title")),
			Description: strings.TrimSpace(r.FormValue("description")),
			AssigneeID:  assigneeID,
		}
		if err := utils.Validate.Struct(payload); err != nil {
			http.Error(w, "invalid task payload", http.StatusBadRequest)
			return
		}
		if _, err := h.store.CreateTask(goalID, ownerID, payload); err != nil {
			status := http.StatusInternalServerError
			if err == tracker.ErrForbidden {
				status = http.StatusForbidden
			}
			http.Error(w, err.Error(), status)
			return
		}
	} else {
		taskID, err := strconv.Atoi(taskIDRaw)
		if err != nil || taskID <= 0 {
			http.Error(w, "invalid task id", http.StatusBadRequest)
			return
		}
		payload := types.UpdateTaskPayload{
			GoalID:      goalID,
			Title:       strings.TrimSpace(r.FormValue("title")),
			Description: strings.TrimSpace(r.FormValue("description")),
			Status:      strings.TrimSpace(r.FormValue("status")),
			AssigneeID:  assigneeID,
		}
		if err := utils.Validate.Struct(payload); err != nil {
			http.Error(w, "invalid task payload", http.StatusBadRequest)
			return
		}
		if _, err := h.store.UpdateTask(taskID, ownerID, payload); err != nil {
			status := http.StatusInternalServerError
			if err == tracker.ErrForbidden {
				status = http.StatusForbidden
			}
			http.Error(w, err.Error(), status)
			return
		}
	}

	h.HandleHTMXTasks(w, r)
}

func renderGoalsTable(goals []*types.GoalWithTasks) string {
	if len(goals) == 0 {
		return `<div class="empty">No goals found for these filters.</div>`
	}

	var b strings.Builder
	b.WriteString(`<table class="grid-table"><thead><tr><th>Goal</th><th>Owner</th><th>Tasks</th><th>Created</th><th>Operations</th></tr></thead><tbody>`)
	for _, goal := range goals {
		b.WriteString(`<tr><td><strong>`)
		b.WriteString(html.EscapeString(goal.Title))
		b.WriteString(`</strong><div class="sub">`)
		b.WriteString(html.EscapeString(goal.Description))
		b.WriteString(`</div></td><td>`)
		ownerName := goal.OwnerName
		if ownerName == "" {
			ownerName = "You"
		}
		b.WriteString(html.EscapeString(ownerName))
		b.WriteString(`</td><td>`)
		b.WriteString(strconv.Itoa(len(goal.Tasks)))
		b.WriteString(`</td><td>`)
		b.WriteString(goal.CreatedAt.Format("2006-01-02"))
		b.WriteString(`</td><td><button hx-get="/htmx/goals/card/`)
		b.WriteString(strconv.Itoa(goal.ID))
		b.WriteString(`" hx-target="#goalCard" hx-swap="innerHTML">Edit</button></td></tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return b.String()
}

func renderTasksTable(tasks []*types.Task) string {
	if len(tasks) == 0 {
		return `<div class="empty">No tasks found for these filters.</div>`
	}

	var b strings.Builder
	b.WriteString(`<table class="grid-table"><thead><tr><th>Task</th><th>Goal</th><th>Status</th><th>Assignee</th><th>Created By</th><th>Operations</th></tr></thead><tbody>`)
	for _, task := range tasks {
		assigneeName := task.AssigneeName
		if assigneeName == "" {
			assigneeName = "Unassigned"
		}
		creatorName := task.CreatedByName
		if creatorName == "" {
			creatorName = "Unknown"
		}

		b.WriteString(`<tr><td><strong>`)
		b.WriteString(html.EscapeString(task.Title))
		b.WriteString(`</strong><div class="sub">`)
		b.WriteString(html.EscapeString(task.Description))
		b.WriteString(`</div></td><td>`)
		b.WriteString(html.EscapeString(task.GoalTitle))
		b.WriteString(`</td><td>`)
		b.WriteString(html.EscapeString(task.Status))
		b.WriteString(`</td><td>`)
		b.WriteString(html.EscapeString(assigneeName))
		b.WriteString(`</td><td>`)
		b.WriteString(html.EscapeString(creatorName))
		b.WriteString(`</td><td><button hx-get="/htmx/tasks/card/`)
		b.WriteString(strconv.Itoa(task.ID))
		b.WriteString(`" hx-target="#taskCard" hx-swap="innerHTML">Edit</button></td></tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return b.String()
}

func renderGoalCard(goal *types.Goal) string {
	isEdit := goal != nil
	title := "Create Goal"
	actionLabel := "Create"
	goalIDValue := ""
	goalTitle := ""
	goalDescription := ""
	if isEdit {
		title = "Edit Goal"
		actionLabel = "Update"
		goalIDValue = strconv.Itoa(goal.ID)
		goalTitle = goal.Title
		goalDescription = goal.Description
	}

	var b strings.Builder
	b.WriteString(`<div class="card-form"><h3>`)
	b.WriteString(html.EscapeString(title))
	b.WriteString(`</h3><form class="stack" hx-post="/htmx/goals/save" hx-target="#goalsTable" hx-swap="innerHTML">`)
	b.WriteString(`<input type="hidden" name="goalId" value="`)
	b.WriteString(html.EscapeString(goalIDValue))
	b.WriteString(`">`)
	b.WriteString(`<label>Title</label><input name="title" value="`)
	b.WriteString(html.EscapeString(goalTitle))
	b.WriteString(`" required>`)
	b.WriteString(`<label>Description</label><textarea name="description" required>`)
	b.WriteString(html.EscapeString(goalDescription))
	b.WriteString(`</textarea><button type="submit">`)
	b.WriteString(html.EscapeString(actionLabel))
	b.WriteString(`</button></form></div>`)
	return b.String()
}

func renderTaskCard(task *types.Task, goals []*types.GoalWithTasks, users []*types.UserLookup) string {
	if len(goals) == 0 {
		return `<div class="empty">Create at least one goal before creating tasks.</div>`
	}

	isEdit := task != nil && task.ID > 0
	if task == nil {
		task = &types.Task{Status: "todo", GoalID: goals[0].ID}
	}
	if task.Status == "" {
		task.Status = "todo"
	}

	title := "Create Task"
	actionLabel := "Create"
	taskIDValue := ""
	if isEdit {
		title = "Edit Task"
		actionLabel = "Update"
		taskIDValue = strconv.Itoa(task.ID)
	}

	var b strings.Builder
	b.WriteString(`<div class="card-form"><h3>`)
	b.WriteString(html.EscapeString(title))
	b.WriteString(`</h3><form class="stack" hx-post="/htmx/tasks/save" hx-target="#tasksTable" hx-swap="innerHTML">`)
	b.WriteString(`<input type="hidden" name="taskId" value="`)
	b.WriteString(html.EscapeString(taskIDValue))
	b.WriteString(`">`)

	b.WriteString(`<label>Goal</label><select name="goalId" required>`)
	for _, goal := range goals {
		selected := ""
		if goal.ID == task.GoalID {
			selected = ` selected`
		}
		b.WriteString(`<option value="`)
		b.WriteString(strconv.Itoa(goal.ID))
		b.WriteString(`"`)
		b.WriteString(selected)
		b.WriteString(`>`)
		b.WriteString(html.EscapeString(goal.Title))
		b.WriteString(`</option>`)
	}
	b.WriteString(`</select>`)

	b.WriteString(`<label>Title</label><input name="title" value="`)
	b.WriteString(html.EscapeString(task.Title))
	b.WriteString(`" required>`)
	b.WriteString(`<label>Description</label><textarea name="description" required>`)
	b.WriteString(html.EscapeString(task.Description))
	b.WriteString(`</textarea>`)

	b.WriteString(`<label>Status</label><select name="status" required>`)
	for _, status := range []string{"todo", "in_progress", "done"} {
		selected := ""
		if status == task.Status {
			selected = ` selected`
		}
		b.WriteString(`<option value="`)
		b.WriteString(status)
		b.WriteString(`"`)
		b.WriteString(selected)
		b.WriteString(`>`)
		b.WriteString(status)
		b.WriteString(`</option>`)
	}
	b.WriteString(`</select>`)

	b.WriteString(`<label>Assignee</label><select name="assigneeId"><option value="">Unassigned</option>`)
	for _, user := range users {
		selected := ""
		if task.AssigneeID != nil && *task.AssigneeID == user.ID {
			selected = ` selected`
		}
		b.WriteString(`<option value="`)
		b.WriteString(strconv.Itoa(user.ID))
		b.WriteString(`"`)
		b.WriteString(selected)
		b.WriteString(`>`)
		b.WriteString(html.EscapeString(user.Name))
		b.WriteString(`</option>`)
	}
	b.WriteString(`</select><button type="submit">`)
	b.WriteString(html.EscapeString(actionLabel))
	b.WriteString(`</button></form></div>`)
	return b.String()
}

func filterGoals(goals []*types.GoalWithTasks, query, status string) []*types.GoalWithTasks {
	if query == "" && status == "" {
		return goals
	}
	q := strings.ToLower(query)
	result := make([]*types.GoalWithTasks, 0)
	for _, goal := range goals {
		if q != "" {
			combined := strings.ToLower(goal.Title + " " + goal.Description)
			if !strings.Contains(combined, q) {
				continue
			}
		}
		if status != "" {
			matched := false
			for _, task := range goal.Tasks {
				if task.Status == status {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		result = append(result, goal)
	}
	return result
}

func filterTasks(tasks []*types.Task, query, status string) []*types.Task {
	if query == "" && status == "" {
		return tasks
	}
	q := strings.ToLower(query)
	result := make([]*types.Task, 0)
	for _, task := range tasks {
		if q != "" {
			combined := strings.ToLower(task.Title + " " + task.Description + " " + task.GoalTitle + " " + task.AssigneeName + " " + task.CreatedByName)
			if !strings.Contains(combined, q) {
				continue
			}
		}
		if status != "" && task.Status != status {
			continue
		}
		result = append(result, task)
	}
	return result
}

func flattenOwnedTasks(goals []*types.GoalWithTasks) []*types.Task {
	tasks := make([]*types.Task, 0)
	for _, goal := range goals {
		for _, task := range goal.Tasks {
			copyTask := *task
			if copyTask.GoalTitle == "" {
				copyTask.GoalTitle = goal.Title
			}
			if copyTask.AssigneeName == "" && copyTask.AssigneeID != nil {
				copyTask.AssigneeName = fmt.Sprintf("User %d", *copyTask.AssigneeID)
			}
			if copyTask.CreatedByName == "" && copyTask.CreatedBy > 0 {
				copyTask.CreatedByName = fmt.Sprintf("User %d", copyTask.CreatedBy)
			}
			tasks = append(tasks, &copyTask)
		}
	}
	return tasks
}

func findGoal(goals []*types.GoalWithTasks, goalID int) *types.Goal {
	for _, goal := range goals {
		if goal.ID == goalID {
			copyGoal := goal.Goal
			return &copyGoal
		}
	}
	return nil
}

func findTask(tasks []*types.Task, taskID int) *types.Task {
	for _, task := range tasks {
		if task.ID == taskID {
			copyTask := *task
			return &copyTask
		}
	}
	return nil
}
