package task

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

// Status represents the current state of a task.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusInReview   Status = "in_review"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusInReview, StatusDone, StatusCancelled:
		return true
	}
	return false
}

// allowedTransitions defines the valid state machine transitions.
var allowedTransitions = map[Status][]Status{
	StatusTodo:       {StatusInProgress, StatusCancelled},
	StatusInProgress: {StatusInReview, StatusDone, StatusTodo, StatusCancelled},
	StatusInReview:   {StatusInProgress, StatusDone, StatusCancelled},
	StatusDone:       {StatusTodo},
	StatusCancelled:  {StatusTodo},
}

func canTransition(from, to Status) bool {
	targets, ok := allowedTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}

// Priority represents the urgency level of a task.
type Priority string

const (
	PriorityNone   Priority = "none"
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityNone, PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	}
	return false
}

type Task struct {
	ID          shared.ID
	ProjectID   shared.ID
	SectionID   *shared.ID
	AssigneeID  *shared.ID
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     *time.Time
	Position    int
	LabelIDs    []shared.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(projectID shared.ID, title string, position int) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, shared.NewValidationError("INVALID_TITLE", "task title must not be empty", "title")
	}

	now := time.Now()
	return &Task{
		ID:        shared.NewID(),
		ProjectID: projectID,
		Title:     title,
		Status:    StatusTodo,
		Priority:  PriorityNone,
		Position:  position,
		LabelIDs:  make([]shared.ID, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Task) UpdateDetails(title, description string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return shared.NewValidationError("INVALID_TITLE", "task title must not be empty", "title")
	}

	t.Title = title
	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) ChangeStatus(newStatus Status) error {
	if !newStatus.IsValid() {
		return shared.NewValidationError("INVALID_STATUS", "invalid task status", "status")
	}
	if t.Status == newStatus {
		return nil
	}
	if !canTransition(t.Status, newStatus) {
		return shared.NewDomainError("INVALID_STATUS_TRANSITION",
			"cannot transition from "+string(t.Status)+" to "+string(newStatus))
	}

	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) ChangePriority(priority Priority) error {
	if !priority.IsValid() {
		return shared.NewValidationError("INVALID_PRIORITY", "invalid task priority", "priority")
	}
	t.Priority = priority
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) Assign(userID *shared.ID) {
	t.AssigneeID = userID
	t.UpdatedAt = time.Now()
}

func (t *Task) MoveToSection(sectionID *shared.ID, position int) {
	t.SectionID = sectionID
	t.Position = position
	t.UpdatedAt = time.Now()
}

func (t *Task) SetDueDate(dueDate *time.Time) {
	t.DueDate = dueDate
	t.UpdatedAt = time.Now()
}

func (t *Task) AddLabel(labelID shared.ID) error {
	for _, id := range t.LabelIDs {
		if id == labelID {
			return shared.NewDomainError("LABEL_ALREADY_ADDED", "label is already attached to this task")
		}
	}
	t.LabelIDs = append(t.LabelIDs, labelID)
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) RemoveLabel(labelID shared.ID) error {
	for i, id := range t.LabelIDs {
		if id == labelID {
			t.LabelIDs = append(t.LabelIDs[:i], t.LabelIDs[i+1:]...)
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return shared.NewDomainError("LABEL_NOT_FOUND", "label is not attached to this task")
}

// TaskFilters holds criteria for querying tasks.
type TaskFilters struct {
	ProjectID  *shared.ID
	SectionID  *shared.ID
	AssigneeID *shared.ID
	Status     *Status
	Priority   *Priority
	LabelID    *shared.ID
	DueBefore  *time.Time
	DueAfter   *time.Time
	Query      string
}
