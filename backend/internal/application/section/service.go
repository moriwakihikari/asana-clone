package section

import (
	"context"
	"sort"
	"time"

	"asana-clone-backend/internal/domain/shared"
	sectionDomain "asana-clone-backend/internal/domain/section"
)

// SectionResponse is the public representation of a section.
type SectionResponse struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
}

// SectionService handles section operations.
type SectionService struct {
	sectionRepo sectionDomain.SectionRepository
}

// NewSectionService creates a new SectionService.
func NewSectionService(sectionRepo sectionDomain.SectionRepository) *SectionService {
	return &SectionService{sectionRepo: sectionRepo}
}

// Create creates a new section in a project. Position is set to the end.
func (s *SectionService) Create(ctx context.Context, projectID shared.ID, name string) (*SectionResponse, error) {
	// Determine the next position.
	sections, err := s.sectionRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	nextPosition := len(sections)

	sec, err := sectionDomain.NewSection(projectID, name, nextPosition)
	if err != nil {
		return nil, err
	}

	if err := s.sectionRepo.Save(ctx, sec); err != nil {
		return nil, err
	}

	return toSectionResponse(sec), nil
}

// ListByProject returns all sections in a project, ordered by position.
func (s *SectionService) ListByProject(ctx context.Context, projectID shared.ID) ([]*SectionResponse, error) {
	sections, err := s.sectionRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Sort by position.
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Position < sections[j].Position
	})

	results := make([]*SectionResponse, 0, len(sections))
	for _, sec := range sections {
		results = append(results, toSectionResponse(sec))
	}
	return results, nil
}

// Rename changes the name of a section.
func (s *SectionService) Rename(ctx context.Context, sectionID shared.ID, name string) (*SectionResponse, error) {
	sec, err := s.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if sec == nil {
		return nil, shared.ErrNotFound
	}

	if err := sec.Rename(name); err != nil {
		return nil, err
	}

	if err := s.sectionRepo.Save(ctx, sec); err != nil {
		return nil, err
	}

	return toSectionResponse(sec), nil
}

// Reorder moves a section to a new position, shifting other sections accordingly.
func (s *SectionService) Reorder(ctx context.Context, sectionID shared.ID, newPosition int) error {
	sec, err := s.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return err
	}
	if sec == nil {
		return shared.ErrNotFound
	}

	// Get all sections in the project.
	sections, err := s.sectionRepo.FindByProjectID(ctx, sec.ProjectID)
	if err != nil {
		return err
	}

	// Sort by current position.
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Position < sections[j].Position
	})

	// Clamp newPosition.
	maxPos := len(sections) - 1
	if newPosition < 0 {
		newPosition = 0
	}
	if newPosition > maxPos {
		newPosition = maxPos
	}

	oldPosition := sec.Position
	if oldPosition == newPosition {
		return nil
	}

	// Shift positions of affected sections.
	for _, other := range sections {
		if other.ID == sectionID {
			continue
		}
		if oldPosition < newPosition {
			// Moving down: shift items between old+1..new up by 1.
			if other.Position > oldPosition && other.Position <= newPosition {
				_ = other.MoveTo(other.Position - 1)
				if err := s.sectionRepo.Save(ctx, other); err != nil {
					return err
				}
			}
		} else {
			// Moving up: shift items between new..old-1 down by 1.
			if other.Position >= newPosition && other.Position < oldPosition {
				_ = other.MoveTo(other.Position + 1)
				if err := s.sectionRepo.Save(ctx, other); err != nil {
					return err
				}
			}
		}
	}

	// Set the target section to the new position.
	if err := sec.MoveTo(newPosition); err != nil {
		return err
	}
	return s.sectionRepo.Save(ctx, sec)
}

// Delete removes a section by ID.
func (s *SectionService) Delete(ctx context.Context, sectionID shared.ID) error {
	sec, err := s.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return err
	}
	if sec == nil {
		return shared.ErrNotFound
	}

	return s.sectionRepo.Delete(ctx, sectionID)
}

func toSectionResponse(sec *sectionDomain.Section) *SectionResponse {
	return &SectionResponse{
		ID:        sec.ID.String(),
		ProjectID: sec.ProjectID.String(),
		Name:      sec.Name,
		Position:  sec.Position,
		CreatedAt: sec.CreatedAt.Format(time.RFC3339),
	}
}
