package tracker

import (
	"VyacheslavKuchumov/test-backend/types"
	"strings"
	"testing"
)

func TestRenderGoals(t *testing.T) {
	t.Run("renders empty state", func(t *testing.T) {
		out := renderGoals(nil)
		if !strings.Contains(out, "No goals yet.") {
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
		out := renderGoals(goals)
		if strings.Contains(out, "<script>alert(1)</script>") {
			t.Fatalf("expected escaped output, got %s", out)
		}
	})
}

func TestRenderAssignedTasks(t *testing.T) {
	t.Run("renders empty assigned state", func(t *testing.T) {
		out := renderAssignedTasks(nil)
		if !strings.Contains(out, "No tasks assigned to you.") {
			t.Fatalf("expected empty state, got %s", out)
		}
	})
}
