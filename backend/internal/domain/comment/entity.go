package comment

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

type Comment struct {
	ID        shared.ID
	TaskID    shared.ID
	UserID    shared.ID
	Content   string
	IsEdited  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComment(taskID, userID shared.ID, content string) (*Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, shared.NewValidationError("INVALID_CONTENT", "comment content must not be empty", "content")
	}

	now := time.Now()
	return &Comment{
		ID:        shared.NewID(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		IsEdited:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Comment) Edit(content string, editorID shared.ID) error {
	if c.UserID != editorID {
		return shared.ErrForbidden
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return shared.NewValidationError("INVALID_CONTENT", "comment content must not be empty", "content")
	}

	c.Content = content
	c.IsEdited = true
	c.UpdatedAt = time.Now()
	return nil
}
