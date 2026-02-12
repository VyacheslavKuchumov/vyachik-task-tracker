package web

import (
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/service/auth"
	userService "VyacheslavKuchumov/test-backend/service/user"
	"VyacheslavKuchumov/test-backend/types"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tebeka/selenium"
)

const seleniumTestUserID = 42

func TestSelenium_PageAccessControl(t *testing.T) {
	driver := newSeleniumWebDriver(t)
	defer driver.Quit()

	server := httptest.NewServer(newSeleniumTestRouter())
	defer server.Close()

	t.Run("public pages are accessible", func(t *testing.T) {
		mustOpenURL(t, driver, server.URL+"/login")
		mustWaitForPath(t, driver, "/login")
		if _, err := driver.FindElement(selenium.ByCSSSelector, "form[action='/auth/login']"); err != nil {
			t.Fatalf("expected login form to be visible: %v", err)
		}

		mustOpenURL(t, driver, server.URL+"/register")
		mustWaitForPath(t, driver, "/register")
		if _, err := driver.FindElement(selenium.ByCSSSelector, "form[action='/auth/register']"); err != nil {
			t.Fatalf("expected register form to be visible: %v", err)
		}
	})

	t.Run("unauthorized users can not visit protected pages", func(t *testing.T) {
		if err := driver.DeleteAllCookies(); err != nil {
			t.Fatalf("failed to reset cookies: %v", err)
		}

		protectedPaths := []string{
			"/",
			"/goals",
			"/tasks",
			"/goals/edit?id=1",
			"/tasks/edit?id=1",
		}

		for _, path := range protectedPaths {
			mustOpenURL(t, driver, server.URL+path)
			mustWaitForPath(t, driver, "/login")
			mustWaitForElement(t, driver, selenium.ByCSSSelector, "form[action='/auth/login']")
		}
	})

	t.Run("authorized users can visit protected pages", func(t *testing.T) {
		attachAuthCookie(t, driver, server.URL)

		protectedPages := []struct {
			targetPath   string
			expectedPath string
		}{
			{targetPath: "/", expectedPath: "/goals"},
			{targetPath: "/goals", expectedPath: "/goals"},
			{targetPath: "/tasks", expectedPath: "/tasks"},
			{targetPath: "/goals/edit?id=1", expectedPath: "/goals/edit"},
			{targetPath: "/tasks/edit?id=1", expectedPath: "/tasks/edit"},
		}

		for _, tc := range protectedPages {
			mustOpenURL(t, driver, server.URL+tc.targetPath)
			mustWaitForPath(t, driver, tc.expectedPath)

			switch tc.expectedPath {
			case "/goals":
				mustWaitForElement(t, driver, selenium.ByCSSSelector, "#goalsTable .grid-table")
			case "/tasks":
				mustWaitForElement(t, driver, selenium.ByCSSSelector, "#tasksTable .grid-table")
			case "/goals/edit":
				mustWaitForElement(t, driver, selenium.ByCSSSelector, "form[hx-post='/htmx/goals/save']")
			case "/tasks/edit":
				mustWaitForElement(t, driver, selenium.ByCSSSelector, "form[hx-post='/htmx/tasks/save']")
			}
		}
	})
}

func newSeleniumWebDriver(t *testing.T) selenium.WebDriver {
	t.Helper()

	remoteURL := strings.TrimSpace(os.Getenv("SELENIUM_URL"))
	if remoteURL == "" {
		t.Skip("set SELENIUM_URL to run Selenium UI tests")
	}

	browserName := strings.TrimSpace(os.Getenv("SELENIUM_BROWSER"))
	if browserName == "" {
		browserName = "chrome"
	}

	caps := selenium.Capabilities{
		"browserName": browserName,
	}

	switch strings.ToLower(browserName) {
	case "chrome":
		caps["goog:chromeOptions"] = map[string]any{
			"args": []string{
				"--headless",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--window-size=1400,1000",
			},
		}
	case "firefox":
		caps["moz:firefoxOptions"] = map[string]any{
			"args": []string{
				"-headless",
				"--width=1400",
				"--height=1000",
			},
		}
	}

	driver, err := selenium.NewRemote(caps, remoteURL)
	if err != nil {
		t.Fatalf("failed to connect to selenium server (%s): %v", remoteURL, err)
	}

	return driver
}

func newSeleniumTestRouter() http.Handler {
	r := chi.NewRouter()

	userStore := &seleniumUserStore{
		usersByID: map[int]*types.User{
			seleniumTestUserID: {
				ID:        seleniumTestUserID,
				Email:     "ui-test@example.com",
				FirstName: "UI",
				LastName:  "Test",
			},
		},
	}
	webHandler := NewHandler(&seleniumGoalTaskStore{})
	userHandler := userService.NewHandler(userStore)

	userService.RegisterRoutes(r, userHandler)
	RegisterRoutes(r, webHandler, userStore)

	return r
}

func attachAuthCookie(t *testing.T, driver selenium.WebDriver, baseURL string) {
	t.Helper()

	if err := driver.DeleteAllCookies(); err != nil {
		t.Fatalf("failed to clear cookies: %v", err)
	}

	mustOpenURL(t, driver, baseURL+"/login")
	mustWaitForPath(t, driver, "/login")

	token, err := auth.CreateJWT([]byte(config.Envs.JWTSecret), seleniumTestUserID)
	if err != nil {
		t.Fatalf("failed to create jwt token: %v", err)
	}

	if err := driver.AddCookie(&selenium.Cookie{
		Name:  auth.AuthCookieName,
		Value: token,
		Path:  "/",
	}); err != nil {
		t.Fatalf("failed to add auth cookie: %v", err)
	}
}

func mustOpenURL(t *testing.T, driver selenium.WebDriver, target string) {
	t.Helper()
	if err := driver.Get(target); err != nil {
		t.Fatalf("failed to open %s: %v", target, err)
	}
}

func mustWaitForPath(t *testing.T, driver selenium.WebDriver, expectedPath string) {
	t.Helper()

	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		currentURL, err := driver.CurrentURL()
		if err == nil {
			path, parseErr := extractPath(currentURL)
			if parseErr == nil && path == expectedPath {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	currentURL, err := driver.CurrentURL()
	if err != nil {
		t.Fatalf("failed to read current URL while waiting for %s: %v", expectedPath, err)
	}

	currentPath, parseErr := extractPath(currentURL)
	if parseErr != nil {
		t.Fatalf("failed to parse current URL %q: %v", currentURL, parseErr)
	}

	t.Fatalf("expected current path %q, got %q (url=%s)", expectedPath, currentPath, currentURL)
}

func mustWaitForElement(t *testing.T, driver selenium.WebDriver, by, selector string) {
	t.Helper()

	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		element, err := driver.FindElement(by, selector)
		if err == nil && element != nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("expected to find element by %s=%q", by, selector)
}

func extractPath(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse url %q: %w", raw, err)
	}
	return parsed.Path, nil
}

type seleniumUserStore struct {
	usersByID map[int]*types.User
}

func (s *seleniumUserStore) GetUserByEmail(string) (*types.User, error) {
	return nil, fmt.Errorf("not found")
}

func (s *seleniumUserStore) GetUserByID(id int) (*types.User, error) {
	user, ok := s.usersByID[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return user, nil
}

func (s *seleniumUserStore) CreateUser(types.User) error {
	return nil
}

type seleniumGoalTaskStore struct{}

func (s *seleniumGoalTaskStore) CreateGoal(int, types.CreateGoalPayload) (*types.Goal, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *seleniumGoalTaskStore) UpdateGoal(int, int, types.CreateGoalPayload) (*types.Goal, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *seleniumGoalTaskStore) GetGoalsByOwner(int) ([]*types.GoalWithTasks, error) {
	return []*types.GoalWithTasks{
		{
			Goal: types.Goal{
				ID:          1,
				Title:       "UI Test Goal",
				Description: "Goal used by selenium tests",
				OwnerID:     seleniumTestUserID,
				OwnerName:   "UI Test",
				CreatedAt:   time.Now(),
			},
			Tasks: []*types.Task{
				{
					ID:            1,
					GoalID:        1,
					GoalTitle:     "UI Test Goal",
					Title:         "UI Test Task",
					Description:   "Task used by selenium tests",
					Status:        "todo",
					CreatedBy:     seleniumTestUserID,
					CreatedByName: "UI Test",
					CreatedAt:     time.Now(),
				},
			},
		},
	}, nil
}

func (s *seleniumGoalTaskStore) CreateTask(int, int, types.CreateTaskPayload) (*types.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *seleniumGoalTaskStore) UpdateTask(int, int, types.UpdateTaskPayload) (*types.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *seleniumGoalTaskStore) AssignTask(int, int, types.AssignTaskPayload) (*types.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *seleniumGoalTaskStore) GetAssignedTasks(int) ([]*types.Task, error) {
	return []*types.Task{}, nil
}

func (s *seleniumGoalTaskStore) ListUsers() ([]*types.UserLookup, error) {
	return []*types.UserLookup{
		{
			ID:   seleniumTestUserID,
			Name: "UI Test",
		},
	}, nil
}
