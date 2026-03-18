package comment

import (
	"context"
	"time"

	"asana-clone-backend/internal/domain/shared"
	commentDomain "asana-clone-backend/internal/domain/comment"
	userDomain "asana-clone-backend/internal/domain/user"
)

// AuthorInfo holds minimal author information for comment responses.
type AuthorInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// CommentResponse is the public representation of a comment.
type CommentResponse struct {
	ID        string      `json:"id"`
	TaskID    string      `json:"task_id"`
	Author    *AuthorInfo `json:"author"`
	Content   string      `json:"content"`
	IsEdited  bool        `json:"is_edited"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// CommentService handles comment operations.
type CommentService struct {
	commentRepo commentDomain.CommentRepository
	userRepo    userDomain.UserRepository
}

// NewCommentService creates a new CommentService.
func NewCommentService(commentRepo commentDomain.CommentRepository, userRepo userDomain.UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

// AddComment creates a new comment on a task.
func (s *CommentService) AddComment(ctx context.Context, taskID, userID shared.ID, content string) (*CommentResponse, error) {
	c, err := commentDomain.NewComment(taskID, userID, content)
	if err != nil {
		return nil, err
	}

	if err := s.commentRepo.Save(ctx, c); err != nil {
		return nil, err
	}

	return s.enrichCommentResponse(ctx, c)
}

// ListByTask returns comments for a task, supporting simple offset-based pagination.
func (s *CommentService) ListByTask(ctx context.Context, taskID shared.ID, offset, limit int) ([]*CommentResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	comments, err := s.commentRepo.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Apply pagination.
	total := len(comments)
	if offset >= total {
		return []*CommentResponse{}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page := comments[offset:end]

	results := make([]*CommentResponse, 0, len(page))
	for _, c := range page {
		resp, err := s.enrichCommentResponse(ctx, c)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

// EditComment edits a comment's content. Only the author can edit.
func (s *CommentService) EditComment(ctx context.Context, commentID, editorID shared.ID, content string) (*CommentResponse, error) {
	c, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, shared.ErrNotFound
	}

	if err := c.Edit(content, editorID); err != nil {
		return nil, err
	}

	if err := s.commentRepo.Save(ctx, c); err != nil {
		return nil, err
	}

	return s.enrichCommentResponse(ctx, c)
}

// DeleteComment deletes a comment. Only the author can delete.
func (s *CommentService) DeleteComment(ctx context.Context, commentID, callerID shared.ID) error {
	c, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if c == nil {
		return shared.ErrNotFound
	}

	if c.UserID != callerID {
		return shared.ErrForbidden
	}

	return s.commentRepo.Delete(ctx, commentID)
}

// enrichCommentResponse converts a domain comment to a response, resolving author details.
func (s *CommentService) enrichCommentResponse(ctx context.Context, c *commentDomain.Comment) (*CommentResponse, error) {
	resp := &CommentResponse{
		ID:        c.ID.String(),
		TaskID:    c.TaskID.String(),
		Content:   c.Content,
		IsEdited:  c.IsEdited,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}

	u, err := s.userRepo.FindByID(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	if u != nil {
		resp.Author = &AuthorInfo{
			ID:        u.ID.String(),
			Name:      u.Name,
			AvatarURL: u.AvatarURL,
		}
	}

	return resp, nil
}
