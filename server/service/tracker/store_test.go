package tracker

import (
	"database/sql"
	"errors"
	"testing"
	"time"
)

type stubScanner struct {
	values []any
	err    error
}

func (s stubScanner) Scan(dest ...any) error {
	if s.err != nil {
		return s.err
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *int:
			*d = s.values[i].(int)
		case *string:
			*d = s.values[i].(string)
		case *bool:
			*d = s.values[i].(bool)
		case *time.Time:
			*d = s.values[i].(time.Time)
		case *sql.NullInt64:
			*d = s.values[i].(sql.NullInt64)
		default:
			return errors.New("unsupported destination type")
		}
	}
	return nil
}

func TestScanRowIntoGoal(t *testing.T) {
	now := time.Now()
	goal, err := scanRowIntoGoal(stubScanner{
		values: []any{1, "Build app", "Description", "high", "in_progress", 2, now},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.ID != 1 || goal.OwnerID != 2 || goal.Priority != "high" || goal.Status != "in_progress" {
		t.Fatalf("unexpected goal data: %+v", goal)
	}
}

func TestScanRowIntoTask(t *testing.T) {
	now := time.Now()
	task, err := scanRowIntoTask(stubScanner{
		values: []any{
			1,
			2,
			"Task title",
			"Task description",
			"low",
			true,
			sql.NullInt64{Int64: 4, Valid: true},
			3,
			now,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != 1 || task.GoalID != 2 || !task.IsCompleted || task.Priority != "low" || task.AssigneeID == nil || *task.AssigneeID != 4 {
		t.Fatalf("unexpected task data: %+v", task)
	}
}
