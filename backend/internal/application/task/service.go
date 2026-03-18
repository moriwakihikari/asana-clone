package task

import (
	"context"
	"time"

	"asana-clone-backend/internal/domain/shared"
	labelDomain "asana-clone-backend/internal/domain/label"
	sectionDomain "asana-clone-backend/internal/domain/section"
	taskDomain "asana-clone-backend/internal/domain/task"
	userDomain "asana-clone-backend/internal/domain/user"
)

// AssigneeInfo holds minimal assignee information for responses.
type AssigneeInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// LabelInfo holds minimal label information for responses.
type LabelInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// TaskResponse is the public representation of a task.
type TaskResponse struct {
	ID          string        `json:"id"`
	ProjectID   string        `json:"project_id"`
	SectionID   *string       `json:"section_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Priority    string        `json:"priority"`
	DueDate     *string       `json:"due_date,omitempty"`
	Position    int           `json:"position"`
	Assignee    *AssigneeInfo `json:"assignee,omitempty"`
	Labels      []*LabelInfo  `json:"labels"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// TaskService handles task operations.
type TaskService struct {
	taskRepo    taskDomain.TaskRepository
	sectionRepo sectionDomain.SectionRepository
	labelRepo   labelDomain.LabelRepository
	userRepo    userDomain.UserRepository
}

// NewTaskService creates a new TaskService.
func NewTaskService(
	taskRepo taskDomain.TaskRepository,
	sectionRepo sectionDomain.SectionRepository,
	labelRepo labelDomain.LabelRepository,
	userRepo userDomain.UserRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		sectionRepo: sectionRepo,
		labelRepo:   labelRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new task in a project.
func (s *TaskService) Create(ctx context.Context, projectID shared.ID, title string, sectionID *shared.ID) (*TaskResponse, error) {
	// Determine position by counting existing tasks.
	existing, err := s.taskRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	t, err := taskDomain.NewTask(projectID, title, len(existing))
	if err != nil {
		return nil, err
	}

	// Validate section if provided.
	if sectionID != nil {
		sec, err := s.sectionRepo.FindByID(ctx, *sectionID)
		if err != nil {
			return nil, err
		}
		if sec == nil {
			return nil, shared.NewDomainError("SECTION_NOT_FOUND", "the specified section does not exist")
		}
		if sec.ProjectID != projectID {
			return nil, shared.NewDomainError("SECTION_MISMATCH", "the section does not belong to this project")
		}
		t.SectionID = sectionID
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// GetByID retrieves a task by ID with full details.
func (s *TaskService) GetByID(ctx context.Context, id shared.ID) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	return s.enrichTaskResponse(ctx, t)
}

// List retrieves tasks matching the given filters.
func (s *TaskService) List(ctx context.Context, filters taskDomain.TaskFilters) ([]*TaskResponse, error) {
	tasks, err := s.taskRepo.FindWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	results := make([]*TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp, err := s.enrichTaskResponse(ctx, t)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

// Update modifies a task's title and description.
func (s *TaskService) Update(ctx context.Context, id shared.ID, title, description string) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	if err := t.UpdateDetails(title, description); err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// ChangeStatus transitions a task's status.
func (s *TaskService) ChangeStatus(ctx context.Context, id shared.ID, status string) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	if err := t.ChangeStatus(taskDomain.Status(status)); err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// ChangePriority changes a task's priority.
func (s *TaskService) ChangePriority(ctx context.Context, id shared.ID, priority string) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	if err := t.ChangePriority(taskDomain.Priority(priority)); err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// Assign assigns or unassigns a user to/from a task.
func (s *TaskService) Assign(ctx context.Context, id shared.ID, assigneeID *shared.ID) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	// Validate assignee exists if provided.
	if assigneeID != nil {
		u, err := s.userRepo.FindByID(ctx, *assigneeID)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "the specified assignee does not exist")
		}
	}

	t.Assign(assigneeID)

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// Move moves a task to a different section and/or position.
func (s *TaskService) Move(ctx context.Context, id shared.ID, sectionID *shared.ID, position int) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	// Validate section if provided.
	if sectionID != nil {
		sec, err := s.sectionRepo.FindByID(ctx, *sectionID)
		if err != nil {
			return nil, err
		}
		if sec == nil {
			return nil, shared.NewDomainError("SECTION_NOT_FOUND", "the specified section does not exist")
		}
		if sec.ProjectID != t.ProjectID {
			return nil, shared.NewDomainError("SECTION_MISMATCH", "the section does not belong to this project")
		}
	}

	t.MoveToSection(sectionID, position)

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// SetDueDate sets or clears the due date for a task.
func (s *TaskService) SetDueDate(ctx context.Context, id shared.ID, dueDate *time.Time) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	t.SetDueDate(dueDate)

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// AddLabel attaches a label to a task.
func (s *TaskService) AddLabel(ctx context.Context, taskID, labelID shared.ID) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	// Validate label exists.
	lbl, err := s.labelRepo.FindByID(ctx, labelID)
	if err != nil {
		return nil, err
	}
	if lbl == nil {
		return nil, shared.NewDomainError("LABEL_NOT_FOUND", "the specified label does not exist")
	}

	if err := t.AddLabel(labelID); err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// RemoveLabel detaches a label from a task.
func (s *TaskService) RemoveLabel(ctx context.Context, taskID, labelID shared.ID) (*TaskResponse, error) {
	t, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, shared.ErrNotFound
	}

	if err := t.RemoveLabel(labelID); err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return s.enrichTaskResponse(ctx, t)
}

// Delete removes a task by ID.
func (s *TaskService) Delete(ctx context.Context, id shared.ID) error {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if t == nil {
		return shared.ErrNotFound
	}

	return s.taskRepo.Delete(ctx, id)
}

// GetMyTasks returns tasks assigned to the given user across a workspace.
func (s *TaskService) GetMyTasks(ctx context.Context, workspaceID, userID shared.ID) ([]*TaskResponse, error) {
	filters := taskDomain.TaskFilters{
		AssigneeID: &userID,
	}

	tasks, err := s.taskRepo.FindWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	results := make([]*TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp, err := s.enrichTaskResponse(ctx, t)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

// enrichTaskResponse converts a domain task to a response, resolving assignee and label details.
func (s *TaskService) enrichTaskResponse(ctx context.Context, t *taskDomain.Task) (*TaskResponse, error) {
	resp := &TaskResponse{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		Position:    t.Position,
		Labels:      make([]*LabelInfo, 0),
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}

	if t.SectionID != nil {
		sid := t.SectionID.String()
		resp.SectionID = &sid
	}

	if t.DueDate != nil {
		dd := t.DueDate.Format(time.RFC3339)
		resp.DueDate = &dd
	}

	// Resolve assignee.
	if t.AssigneeID != nil {
		u, err := s.userRepo.FindByID(ctx, *t.AssigneeID)
		if err != nil {
			return nil, err
		}
		if u != nil {
			resp.Assignee = &AssigneeInfo{
				ID:        u.ID.String(),
				Name:      u.Name,
				AvatarURL: u.AvatarURL,
			}
		}
	}

	// Resolve labels.
	for _, labelID := range t.LabelIDs {
		lbl, err := s.labelRepo.FindByID(ctx, labelID)
		if err != nil {
			return nil, err
		}
		if lbl != nil {
			resp.Labels = append(resp.Labels, &LabelInfo{
				ID:    lbl.ID.String(),
				Name:  lbl.Name,
				Color: lbl.Color,
			})
		}
	}

	return resp, nil
}
