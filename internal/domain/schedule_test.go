package domain

import "testing"

func TestNewSchedule(t *testing.T) {
	s := NewSchedule("admin")
	if s.ID == "" {
		t.Error("expected ID to be set")
	}
	if s.CreatedBy != "admin" {
		t.Errorf("expected createdBy 'admin', got %q", s.CreatedBy)
	}
}

func TestScheduleUpdateDays(t *testing.T) {
	s := NewSchedule("admin")
	originalUpdatedAt := s.UpdatedAt

	s.UpdateDays("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", "editor")

	days := map[string]string{
		"Monday": s.Monday, "Tuesday": s.Tuesday, "Wednesday": s.Wednesday,
		"Thursday": s.Thursday, "Friday": s.Friday, "Saturday": s.Saturday, "Sunday": s.Sunday,
	}
	expected := map[string]string{
		"Monday": "Mon", "Tuesday": "Tue", "Wednesday": "Wed",
		"Thursday": "Thu", "Friday": "Fri", "Saturday": "Sat", "Sunday": "Sun",
	}
	for day, got := range days {
		if got != expected[day] {
			t.Errorf("expected %s %q, got %q", day, expected[day], got)
		}
	}
	if s.UpdatedBy != "editor" {
		t.Errorf("expected updatedBy 'editor', got %q", s.UpdatedBy)
	}
	if s.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestScheduleEntityKind(t *testing.T) {
	s := &Schedule{}
	if kind := s.EntityKind(); kind != "Schedule" {
		t.Errorf("expected 'Schedule', got %q", kind)
	}
}
