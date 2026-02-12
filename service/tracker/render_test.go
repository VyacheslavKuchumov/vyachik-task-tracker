package tracker

import (
	"VyacheslavKuchumov/test-backend/types"
	"strings"
	"testing"
)

func TestRenderGoalsTable(t *testing.T) {
	t.Run("renders empty state", func(t *testing.T) {
		out := renderGoalsTable(nil)
		if !strings.Contains(out, "No goals found") {
			t.Fatalf("expected empty state, got %s", out)
		}
	})

	t.Run("escapes title", func(t *testing.T) {
		goals := []*types.GoalWithTasks{
			{
				Goal: types.Goal{
					ID:          1,
					Title:       "<script>alert(1)</script>",
					Description: "Desc",
				},
			},
		}
		out := renderGoalsTable(goals)
		if strings.Contains(out, "<script>alert(1)</script>") {
			t.Fatalf("expected escaped output, got %s", out)
		}
	})
}

func TestRenderTasksTable(t *testing.T) {
	t.Run("renders empty tasks state", func(t *testing.T) {
		out := renderTasksTable(nil)
		if !strings.Contains(out, "No tasks found") {
			t.Fatalf("expected empty state, got %s", out)
		}
	})
}
