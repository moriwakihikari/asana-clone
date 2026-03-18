package handler

import (
	"encoding/json"
	"net/http"

	taskApp "asana-clone-backend/internal/application/task"
	"asana-clone-backend/internal/domain/shared"
	taskDomain "asana-clone-backend/internal/domain/task"
	httpErrors "asana-clone-backend/internal/interfaces/errors"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
)

// TaskHandler handles task-related endpoints.
type TaskHandler struct {
	taskService *taskApp.TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(taskService *taskApp.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

type createTaskRequest struct {
	Title     string  `json:"title"`
	SectionID *string `json:"section_id,omitempty"`
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type changeStatusRequest struct {
	Status string `json:"status"`
}

type moveTaskRequest struct {
	SectionID *string `json:"section_id,omitempty"`
	Position  int     `json:"position"`
}

type assignTaskRequest struct {
	AssigneeID *string `json:"assignee_id,omitempty"`
}

type labelRequest struct {
	LabelID string `json:"label_id"`
}

// Create handles POST /projects/{projectID}/tasks.
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	var sectionID *shared.ID
	if req.SectionID != nil {
		id, err := shared.ParseID(*req.SectionID)
		if err != nil {
			httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
				Code:    "INVALID_ID",
				Message: "invalid section ID format",
			})
			return
		}
		sectionID = &id
	}

	task, err := h.taskService.Create(r.Context(), projectID, req.Title, sectionID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, task)
}

// List handles GET /projects/{projectID}/tasks.
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	filters := taskDomain.TaskFilters{
		ProjectID: &projectID,
	}

	// Parse optional query params.
	if sID := r.URL.Query().Get("section_id"); sID != "" {
		id, err := shared.ParseID(sID)
		if err == nil {
			filters.SectionID = &id
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := taskDomain.Status(status)
		filters.Status = &s
	}
	if assignee := r.URL.Query().Get("assignee_id"); assignee != "" {
		id, err := shared.ParseID(assignee)
		if err == nil {
			filters.AssigneeID = &id
		}
	}
	if q := r.URL.Query().Get("q"); q != "" {
		filters.Query = q
	}

	tasks, err := h.taskService.List(r.Context(), filters)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, tasks)
}

// GetTask handles GET /tasks/{taskID}.
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	task, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// Update handles PATCH /projects/{projectID}/tasks/{taskID}.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	task, err := h.taskService.Update(r.Context(), taskID, req.Title, req.Description)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// Delete handles DELETE /projects/{projectID}/tasks/{taskID}.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	if err := h.taskService.Delete(r.Context(), taskID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ChangeStatus handles POST /projects/{projectID}/tasks/{taskID}/status.
func (h *TaskHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req changeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	task, err := h.taskService.ChangeStatus(r.Context(), taskID, req.Status)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// Move handles POST /projects/{projectID}/tasks/{taskID}/move.
func (h *TaskHandler) Move(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req moveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	var sectionID *shared.ID
	if req.SectionID != nil {
		id, err := shared.ParseID(*req.SectionID)
		if err != nil {
			httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
				Code:    "INVALID_ID",
				Message: "invalid section ID format",
			})
			return
		}
		sectionID = &id
	}

	task, err := h.taskService.Move(r.Context(), taskID, sectionID, req.Position)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// Assign handles POST /projects/{projectID}/tasks/{taskID}/assign.
func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req assignTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	var assigneeID *shared.ID
	if req.AssigneeID != nil {
		id, err := shared.ParseID(*req.AssigneeID)
		if err != nil {
			httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
				Code:    "INVALID_ID",
				Message: "invalid assignee ID format",
			})
			return
		}
		assigneeID = &id
	}

	task, err := h.taskService.Assign(r.Context(), taskID, assigneeID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// AddLabel handles POST /projects/{projectID}/tasks/{taskID}/labels.
func (h *TaskHandler) AddLabel(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req labelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	labelID, err := shared.ParseID(req.LabelID)
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid label ID format",
		})
		return
	}

	task, err := h.taskService.AddLabel(r.Context(), taskID, labelID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// RemoveLabel handles DELETE /projects/{projectID}/tasks/{taskID}/labels/{labelID}.
func (h *TaskHandler) RemoveLabel(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	labelID, err := parseUUIDParam(r, "labelID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid label ID format",
		})
		return
	}

	task, err := h.taskService.RemoveLabel(r.Context(), taskID, labelID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, task)
}

// GetMyTasks handles GET /workspaces/{workspaceID}/my-tasks.
func (h *TaskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	tasks, err := h.taskService.GetMyTasks(r.Context(), wsID, userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, tasks)
}

// Mount registers task routes on the given router.
// Project-scoped routes are nested under /projects/{projectID}/tasks.
func (h *TaskHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Patch("/{taskID}", h.Update)
	r.Delete("/{taskID}", h.Delete)
	r.Post("/{taskID}/status", h.ChangeStatus)
	r.Post("/{taskID}/move", h.Move)
	r.Post("/{taskID}/assign", h.Assign)
	r.Post("/{taskID}/labels", h.AddLabel)
	r.Delete("/{taskID}/labels/{labelID}", h.RemoveLabel)
}

// MountDirect registers direct task routes (not nested under projects).
func (h *TaskHandler) MountDirect(r chi.Router) {
	r.Get("/{taskID}", h.GetTask)
}

// MountMyTasks registers the my-tasks route under workspaces.
func (h *TaskHandler) MountMyTasks(r chi.Router) {
	r.Get("/my-tasks", h.GetMyTasks)
}
