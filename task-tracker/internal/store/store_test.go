package store

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/itbem-corp/task-tracker/internal/models"
)

func newTestTask(id, title string, createdAt time.Time) models.Task {
	return models.Task{
		ID:          id,
		Title:       title,
		Description: "test description",
		Status:      models.StatusPending,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
}

func strPtr(s string) *string { return &s }

func TestTaskStore_Create(t *testing.T) {
	t.Run("creates task and retrieves it", func(t *testing.T) {
		s := New()
		task := newTestTask("id-1", "Task One", time.Now().UTC())

		created := s.Create(task)

		if created.ID != task.ID {
			t.Errorf("Create() returned ID = %q, want %q", created.ID, task.ID)
		}
		if created.Title != task.Title {
			t.Errorf("Create() returned Title = %q, want %q", created.Title, task.Title)
		}

		got, err := s.GetByID("id-1")
		if err != nil {
			t.Fatalf("GetByID() after Create() returned error: %v", err)
		}
		if got.Title != "Task One" {
			t.Errorf("GetByID() Title = %q, want %q", got.Title, "Task One")
		}
	})

	t.Run("multiple creates maintain separate entries", func(t *testing.T) {
		s := New()
		now := time.Now().UTC()

		s.Create(newTestTask("id-1", "Task One", now))
		s.Create(newTestTask("id-2", "Task Two", now.Add(time.Second)))
		s.Create(newTestTask("id-3", "Task Three", now.Add(2*time.Second)))

		all := s.GetAll()
		if len(all) != 3 {
			t.Errorf("GetAll() returned %d tasks, want 3", len(all))
		}
	})

	t.Run("create overwrites task with same ID", func(t *testing.T) {
		s := New()
		now := time.Now().UTC()

		s.Create(newTestTask("id-1", "Original", now))
		s.Create(newTestTask("id-1", "Overwritten", now))

		got, err := s.GetByID("id-1")
		if err != nil {
			t.Fatalf("GetByID() returned error: %v", err)
		}
		if got.Title != "Overwritten" {
			t.Errorf("GetByID() Title = %q, want %q", got.Title, "Overwritten")
		}

		all := s.GetAll()
		if len(all) != 1 {
			t.Errorf("GetAll() returned %d tasks, want 1 (duplicate ID should overwrite)", len(all))
		}
	})
}

func TestTaskStore_GetAll(t *testing.T) {
	t.Run("empty store returns empty slice", func(t *testing.T) {
		s := New()
		all := s.GetAll()
		if all == nil {
			t.Error("GetAll() returned nil, want empty slice")
		}
		if len(all) != 0 {
			t.Errorf("GetAll() returned %d tasks, want 0", len(all))
		}
	})

	t.Run("returns tasks sorted by created_at descending", func(t *testing.T) {
		s := New()
		now := time.Now().UTC()

		// Insert in ascending order
		s.Create(newTestTask("oldest", "Oldest Task", now))
		s.Create(newTestTask("middle", "Middle Task", now.Add(time.Minute)))
		s.Create(newTestTask("newest", "Newest Task", now.Add(2*time.Minute)))

		all := s.GetAll()
		if len(all) != 3 {
			t.Fatalf("GetAll() returned %d tasks, want 3", len(all))
		}

		// First element should be newest
		if all[0].ID != "newest" {
			t.Errorf("GetAll()[0].ID = %q, want %q (newest first)", all[0].ID, "newest")
		}
		if all[1].ID != "middle" {
			t.Errorf("GetAll()[1].ID = %q, want %q", all[1].ID, "middle")
		}
		if all[2].ID != "oldest" {
			t.Errorf("GetAll()[2].ID = %q, want %q (oldest last)", all[2].ID, "oldest")
		}
	})

	t.Run("single task returns slice of one", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("only", "Only Task", time.Now().UTC()))

		all := s.GetAll()
		if len(all) != 1 {
			t.Fatalf("GetAll() returned %d tasks, want 1", len(all))
		}
		if all[0].ID != "only" {
			t.Errorf("GetAll()[0].ID = %q, want %q", all[0].ID, "only")
		}
	})
}

func TestTaskStore_GetByID(t *testing.T) {
	t.Run("existing ID returns task", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Task One", time.Now().UTC()))

		got, err := s.GetByID("id-1")
		if err != nil {
			t.Fatalf("GetByID() returned error: %v", err)
		}
		if got.ID != "id-1" {
			t.Errorf("GetByID() ID = %q, want %q", got.ID, "id-1")
		}
		if got.Title != "Task One" {
			t.Errorf("GetByID() Title = %q, want %q", got.Title, "Task One")
		}
	})

	t.Run("non-existent ID returns ErrTaskNotFound", func(t *testing.T) {
		s := New()

		_, err := s.GetByID("nonexistent")
		if err == nil {
			t.Fatal("GetByID() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("GetByID() error = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("empty string ID returns ErrTaskNotFound", func(t *testing.T) {
		s := New()

		_, err := s.GetByID("")
		if err == nil {
			t.Fatal("GetByID() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("GetByID() error = %v, want ErrTaskNotFound", err)
		}
	})
}

func TestTaskStore_Update(t *testing.T) {
	t.Run("update title only leaves other fields unchanged", func(t *testing.T) {
		s := New()
		now := time.Now().UTC()
		original := newTestTask("id-1", "Original Title", now)
		original.Description = "Original Description"
		original.Status = models.StatusPending
		s.Create(original)

		updated, err := s.Update("id-1", models.UpdateTaskRequest{
			Title: strPtr("New Title"),
		})
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}
		if updated.Title != "New Title" {
			t.Errorf("Update() Title = %q, want %q", updated.Title, "New Title")
		}
		if updated.Description != "Original Description" {
			t.Errorf("Update() Description changed to %q, want %q", updated.Description, "Original Description")
		}
		if updated.Status != models.StatusPending {
			t.Errorf("Update() Status changed to %q, want %q", updated.Status, models.StatusPending)
		}
	})

	t.Run("update status only", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Task", time.Now().UTC()))

		updated, err := s.Update("id-1", models.UpdateTaskRequest{
			Status: strPtr(models.StatusDone),
		})
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}
		if updated.Status != models.StatusDone {
			t.Errorf("Update() Status = %q, want %q", updated.Status, models.StatusDone)
		}
		if updated.Title != "Task" {
			t.Errorf("Update() Title changed to %q, want %q", updated.Title, "Task")
		}
	})

	t.Run("update all fields", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Old", time.Now().UTC()))

		updated, err := s.Update("id-1", models.UpdateTaskRequest{
			Title:       strPtr("New Title"),
			Description: strPtr("New Desc"),
			Status:      strPtr(models.StatusInProgress),
		})
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}
		if updated.Title != "New Title" {
			t.Errorf("Update() Title = %q, want %q", updated.Title, "New Title")
		}
		if updated.Description != "New Desc" {
			t.Errorf("Update() Description = %q, want %q", updated.Description, "New Desc")
		}
		if updated.Status != models.StatusInProgress {
			t.Errorf("Update() Status = %q, want %q", updated.Status, models.StatusInProgress)
		}
	})

	t.Run("non-existent ID returns ErrTaskNotFound", func(t *testing.T) {
		s := New()

		_, err := s.Update("nonexistent", models.UpdateTaskRequest{
			Title: strPtr("Title"),
		})
		if err == nil {
			t.Fatal("Update() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("Update() error = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("UpdatedAt changes after update", func(t *testing.T) {
		s := New()
		past := time.Now().UTC().Add(-time.Hour)
		s.Create(newTestTask("id-1", "Task", past))

		beforeUpdate := time.Now().UTC()
		updated, err := s.Update("id-1", models.UpdateTaskRequest{
			Title: strPtr("Updated"),
		})
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}
		if !updated.UpdatedAt.After(past) {
			t.Error("Update() UpdatedAt should be after original creation time")
		}
		if updated.UpdatedAt.Before(beforeUpdate) {
			t.Error("Update() UpdatedAt should be >= the time before update call")
		}
	})

	t.Run("update persists in store", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Before", time.Now().UTC()))

		_, err := s.Update("id-1", models.UpdateTaskRequest{
			Title: strPtr("After"),
		})
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}

		got, err := s.GetByID("id-1")
		if err != nil {
			t.Fatalf("GetByID() returned error: %v", err)
		}
		if got.Title != "After" {
			t.Errorf("GetByID() after Update() Title = %q, want %q", got.Title, "After")
		}
	})
}

func TestTaskStore_Delete(t *testing.T) {
	t.Run("delete existing task succeeds", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Task", time.Now().UTC()))

		err := s.Delete("id-1")
		if err != nil {
			t.Errorf("Delete() returned error: %v", err)
		}
	})

	t.Run("delete non-existent task returns ErrTaskNotFound", func(t *testing.T) {
		s := New()

		err := s.Delete("nonexistent")
		if err == nil {
			t.Fatal("Delete() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("Delete() error = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("deleted task no longer in GetAll", func(t *testing.T) {
		s := New()
		now := time.Now().UTC()
		s.Create(newTestTask("id-1", "Task One", now))
		s.Create(newTestTask("id-2", "Task Two", now.Add(time.Second)))

		err := s.Delete("id-1")
		if err != nil {
			t.Fatalf("Delete() returned error: %v", err)
		}

		all := s.GetAll()
		if len(all) != 1 {
			t.Fatalf("GetAll() after Delete() returned %d tasks, want 1", len(all))
		}
		if all[0].ID != "id-2" {
			t.Errorf("GetAll() remaining task ID = %q, want %q", all[0].ID, "id-2")
		}
	})

	t.Run("deleted task not found by GetByID", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Task", time.Now().UTC()))

		_ = s.Delete("id-1")

		_, err := s.GetByID("id-1")
		if err == nil {
			t.Fatal("GetByID() after Delete() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("GetByID() after Delete() error = %v, want ErrTaskNotFound", err)
		}
	})

	t.Run("double delete returns error", func(t *testing.T) {
		s := New()
		s.Create(newTestTask("id-1", "Task", time.Now().UTC()))

		err := s.Delete("id-1")
		if err != nil {
			t.Fatalf("First Delete() returned error: %v", err)
		}

		err = s.Delete("id-1")
		if err == nil {
			t.Fatal("Second Delete() expected error, got nil")
		}
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("Second Delete() error = %v, want ErrTaskNotFound", err)
		}
	})
}

func TestTaskStore_Concurrency(t *testing.T) {
	t.Run("concurrent creates are safe", func(t *testing.T) {
		s := New()
		var wg sync.WaitGroup

		numGoroutines := 100
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(idx int) {
				defer wg.Done()
				id := "id-" + intToStr(idx)
				s.Create(newTestTask(id, "Task "+intToStr(idx), time.Now().UTC()))
			}(i)
		}
		wg.Wait()

		all := s.GetAll()
		if len(all) != numGoroutines {
			t.Errorf("GetAll() returned %d tasks after concurrent creates, want %d", len(all), numGoroutines)
		}
	})

	t.Run("concurrent reads and writes are safe", func(t *testing.T) {
		s := New()
		// Pre-populate
		for i := 0; i < 10; i++ {
			s.Create(newTestTask("id-"+intToStr(i), "Task", time.Now().UTC()))
		}

		var wg sync.WaitGroup
		wg.Add(30)

		// 10 concurrent reads
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				_ = s.GetAll()
			}()
		}

		// 10 concurrent updates
		for i := 0; i < 10; i++ {
			go func(idx int) {
				defer wg.Done()
				_, _ = s.Update("id-"+intToStr(idx), models.UpdateTaskRequest{
					Title: strPtr("Updated"),
				})
			}(i)
		}

		// 10 concurrent GetByID
		for i := 0; i < 10; i++ {
			go func(idx int) {
				defer wg.Done()
				_, _ = s.GetByID("id-" + intToStr(idx))
			}(i)
		}

		wg.Wait()
		// If we get here without a race condition panic, the test passes
	})
}

// intToStr is a simple int-to-string helper to avoid importing strconv.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if neg {
		result = "-" + result
	}
	return result
}
