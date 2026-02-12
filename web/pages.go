package web

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates/*.html static/*
var assetsFS embed.FS

type Handler struct {
	store     types.GoalTaskStore
	templates *template.Template
	staticFS  fs.FS
}

type PageData struct {
	ErrorMessage string
	OKMessage    string
}

type EditPageData struct {
	Title   string
	BackURL string
	CardURL string
}

func NewHandler(store types.GoalTaskStore) *Handler {
	templates := template.Must(template.ParseFS(assetsFS, "templates/*.html"))
	staticFS, err := fs.Sub(assetsFS, "static")
	if err != nil {
		panic(err)
	}

	return &Handler{
		store:     store,
		templates: templates,
		staticFS:  staticFS,
	}
}

func (h *Handler) StaticHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.FS(h.staticFS)))
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(auth.AuthCookieName); err == nil && strings.TrimSpace(cookie.Value) != "" {
		http.Redirect(w, r, "/goals", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "login.html", PageData{
		ErrorMessage: strings.TrimSpace(r.URL.Query().Get("error")),
		OKMessage:    strings.TrimSpace(r.URL.Query().Get("ok")),
	})
}

func (h *Handler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "register.html", PageData{
		ErrorMessage: strings.TrimSpace(r.URL.Query().Get("error")),
	})
}

func (h *Handler) HandleGoalsPage(w http.ResponseWriter, _ *http.Request) {
	h.renderPage(w, "goals.html", nil)
}

func (h *Handler) HandleTasksPage(w http.ResponseWriter, _ *http.Request) {
	h.renderPage(w, "tasks.html", nil)
}

func (h *Handler) HandleGoalEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := requiredIDQueryParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.renderPage(w, "goal-edit.html", EditPageData{
		Title:   "Edit Goal",
		BackURL: "/goals",
		CardURL: fmt.Sprintf("/htmx/goals/card?id=%d&redirectTo=/goals", id),
	})
}

func (h *Handler) HandleTaskEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := requiredIDQueryParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.renderPage(w, "task-edit.html", EditPageData{
		Title:   "Edit Task",
		BackURL: "/tasks",
		CardURL: fmt.Sprintf("/htmx/tasks/card?id=%d&redirectTo=/tasks", id),
	})
}

func (h *Handler) renderPage(w http.ResponseWriter, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, tmpl, data); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func requiredIDQueryParam(r *http.Request, key string) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, fmt.Errorf("query parameter %q is required", key)
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("query parameter %q must be a positive integer", key)
	}

	return id, nil
}
