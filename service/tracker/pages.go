package tracker

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"html"
	"net/http"
	"strings"
)

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(auth.AuthCookieName); err == nil && strings.TrimSpace(cookie.Value) != "" {
		http.Redirect(w, r, "/goals", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	errorMsg := html.EscapeString(r.URL.Query().Get("error"))
	okMsg := html.EscapeString(r.URL.Query().Get("ok"))
	page := loginPage
	if errorMsg != "" {
		page = strings.Replace(page, "{{status}}", `<p class="status error">`+errorMsg+`</p>`, 1)
	} else if okMsg != "" {
		page = strings.Replace(page, "{{status}}", `<p class="status ok">`+okMsg+`</p>`, 1)
	} else {
		page = strings.Replace(page, "{{status}}", "", 1)
	}
	writeHTML(w, page)
}

func (h *Handler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	errorMsg := html.EscapeString(r.URL.Query().Get("error"))
	page := registerPage
	if errorMsg != "" {
		page = strings.Replace(page, "{{status}}", `<p class="status error">`+errorMsg+`</p>`, 1)
	} else {
		page = strings.Replace(page, "{{status}}", "", 1)
	}
	writeHTML(w, page)
}

func (h *Handler) HandleGoalsPage(w http.ResponseWriter, _ *http.Request) {
	writeHTML(w, goalsPage)
}

func (h *Handler) HandleTasksPage(w http.ResponseWriter, _ *http.Request) {
	writeHTML(w, tasksPage)
}

func writeHTML(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(content))
}

const baseStyles = `
  <style>
    :root { --bg:#f4f9ff; --ink:#0f172a; --line:#cbd5e1; --card:#ffffff; --accent:#0f766e; --muted:#475569; --soft:#ecfeff; }
    * { box-sizing:border-box; }
    body { margin:0; font-family: "Trebuchet MS", "Segoe UI", sans-serif; color:var(--ink); background:linear-gradient(135deg,#eef6ff,#f9fff6); }
    .wrap { max-width:1080px; margin:0 auto; padding:24px 16px; }
    .panel { background:var(--card); border:1px solid var(--line); border-radius:14px; padding:16px; box-shadow:0 8px 24px rgba(15,23,42,.06); margin-bottom:16px; }
    h1,h2,h3 { margin:0 0 10px; }
    p { margin:0 0 10px; color:var(--muted); }
    .row { display:grid; grid-template-columns:2fr 1fr; gap:16px; }
    .stack { display:flex; flex-direction:column; gap:10px; }
    input, textarea, select, button { width:100%; border-radius:10px; border:1px solid var(--line); padding:10px 11px; }
    textarea { min-height:90px; resize:vertical; }
    button { background:var(--accent); color:#fff; border:none; font-weight:700; cursor:pointer; }
    .nav { display:flex; align-items:center; gap:10px; flex-wrap:wrap; margin-bottom:14px; }
    .link { display:inline-block; padding:8px 12px; border:1px solid var(--line); border-radius:8px; background:#fff; color:var(--ink); text-decoration:none; }
    .toolbar { display:grid; grid-template-columns:2fr 1fr auto auto; gap:10px; align-items:end; margin-bottom:14px; }
    .toolbar > * { margin:0; }
    .grid-table { width:100%; border-collapse:collapse; border:1px solid var(--line); border-radius:10px; overflow:hidden; }
    .grid-table th, .grid-table td { border-bottom:1px solid var(--line); padding:10px; text-align:left; vertical-align:top; }
    .grid-table thead th { background:var(--soft); font-size:13px; text-transform:uppercase; letter-spacing:.03em; }
    .grid-table tbody tr:last-child td { border-bottom:none; }
    .sub { margin-top:4px; color:var(--muted); font-size:13px; }
    .card-form { border:1px dashed var(--line); border-radius:12px; padding:12px; background:#fff; }
    .status { padding:10px; border-radius:10px; margin-bottom:12px; }
    .status.error { background:#fee2e2; color:#7f1d1d; border:1px solid #fecaca; }
    .status.ok { background:#dcfce7; color:#14532d; border:1px solid #bbf7d0; }
    .empty { color:var(--muted); font-style:italic; }
    @media (max-width:900px) { .row { grid-template-columns:1fr; } .toolbar { grid-template-columns:1fr; } }
  </style>
`

const loginPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Login</title>` + baseStyles + `
</head>
<body>
  <main class="wrap">
    <section class="panel" style="max-width:460px;margin:40px auto;">
      <h1>Sign In</h1>
      <p>Use your account to continue.</p>
      {{status}}
      <form class="stack" method="post" action="/auth/login">
        <input name="email" type="email" placeholder="Email" required>
        <input name="password" type="password" placeholder="Password" required>
        <button type="submit">Login</button>
      </form>
      <p style="margin-top:12px;">No account? <a href="/register">Create one</a></p>
    </section>
  </main>
</body>
</html>`

const registerPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Register</title>` + baseStyles + `
</head>
<body>
  <main class="wrap">
    <section class="panel" style="max-width:460px;margin:40px auto;">
      <h1>Create Account</h1>
      <p>Register to manage goals and tasks.</p>
      {{status}}
      <form class="stack" method="post" action="/auth/register">
        <input name="firstName" placeholder="First name" required>
        <input name="lastName" placeholder="Last name" required>
        <input name="email" type="email" placeholder="Email" required>
        <input name="password" type="password" placeholder="Password" required>
        <button type="submit">Register</button>
      </form>
      <p style="margin-top:12px;">Already registered? <a href="/login">Login</a></p>
    </section>
  </main>
</body>
</html>`

const goalsPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Goals</title>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>` + baseStyles + `
</head>
<body>
  <main class="wrap">
    <div class="nav">
      <a class="link" href="/goals">Goals</a>
      <a class="link" href="/tasks">Tasks</a>
      <form method="post" action="/auth/logout">
        <button type="submit">Logout</button>
      </form>
    </div>

    <section class="panel">
      <h1>Goals</h1>
      <p>List view with filters and toolbar operations. Create and edit one goal in the card panel.</p>
    </section>

    <div class="row">
      <section class="panel">
        <div class="toolbar">
          <input name="q" form="goalsFilter" placeholder="Filter by goal title or description">
          <select name="status" form="goalsFilter">
            <option value="">Any task status</option>
            <option value="todo">todo</option>
            <option value="in_progress">in_progress</option>
            <option value="done">done</option>
          </select>
          <button type="submit" form="goalsFilter">Filter</button>
          <button hx-get="/htmx/goals/card" hx-target="#goalCard" hx-swap="innerHTML">New Goal</button>
        </div>
        <form id="goalsFilter" hx-get="/htmx/goals" hx-target="#goalsTable" hx-swap="innerHTML"></form>
        <div id="goalsTable" hx-get="/htmx/goals" hx-trigger="load" hx-swap="innerHTML"></div>
      </section>

      <section class="panel">
        <h2>Goal Card</h2>
        <p>Create or edit one goal object.</p>
        <div id="goalCard" hx-get="/htmx/goals/card" hx-trigger="load" hx-swap="innerHTML"></div>
      </section>
    </div>

    <section class="panel">
      <h2>Task Operations</h2>
      <p>Use the dedicated tasks page for list/card task management.</p>
      <a class="link" href="/tasks">Open Tasks Page</a>
    </section>
  </main>
</body>
</html>`

const tasksPage = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Tasks</title>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>` + baseStyles + `
</head>
<body>
  <main class="wrap">
    <div class="nav">
      <a class="link" href="/goals">Goals</a>
      <a class="link" href="/tasks">Tasks</a>
      <form method="post" action="/auth/logout">
        <button type="submit">Logout</button>
      </form>
    </div>

    <section class="panel">
      <h1>Tasks</h1>
      <p>List view with filters and toolbar operations. Related goal/user data is shown by name.</p>
    </section>

    <div class="row">
      <section class="panel">
        <div class="toolbar">
          <input name="q" form="tasksFilter" placeholder="Filter by task, goal or user names">
          <select name="status" form="tasksFilter">
            <option value="">Any status</option>
            <option value="todo">todo</option>
            <option value="in_progress">in_progress</option>
            <option value="done">done</option>
          </select>
          <button type="submit" form="tasksFilter">Filter</button>
          <button hx-get="/htmx/tasks/card" hx-target="#taskCard" hx-swap="innerHTML">New Task</button>
        </div>
        <form id="tasksFilter" hx-get="/htmx/tasks" hx-target="#tasksTable" hx-swap="innerHTML"></form>
        <div id="tasksTable" hx-get="/htmx/tasks" hx-trigger="load" hx-swap="innerHTML"></div>
      </section>

      <section class="panel">
        <h2>Task Card</h2>
        <p>Create or edit one task object.</p>
        <div id="taskCard" hx-get="/htmx/tasks/card" hx-trigger="load" hx-swap="innerHTML"></div>
      </section>
    </div>
  </main>
</body>
</html>`
