package services

import (
	"fmt"
	"log/slog"
	"strings"

	"flashcards/db"
	"flashcards/models"
)

type NoteService struct {
	repo db.NoteRepository
}

func NewNoteService(repo db.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(req *models.CreateNoteRequest) (*models.Note, error) {
	slog.Info("creating note")

	if err := s.validateCreateRequest(req); err != nil {
		slog.Error("failed to create note", "error", err)
		return nil, err
	}

	note := &models.Note{
		Content: strings.TrimSpace(req.Content),
	}

	if err := s.repo.CreateNote(note); err != nil {
		wrappedErr := fmt.Errorf("failed to create note: %w", err)
		slog.Error("failed to create note", "error", wrappedErr)
		return nil, wrappedErr
	}

	slog.Info("created note successfully", "noteId", note.ID)

	return note, nil
}

func (s *NoteService) GetNoteByID(id int) (*models.Note, error) {
	slog.Info("getting note by id", "noteId", id)

	if id <= 0 {
		err := fmt.Errorf("invalid note ID: %d", id)
		slog.Error("failed to get note by id", "noteId", id, "error", err)
		return nil, err
	}

	note, err := s.repo.GetNoteByID(id)
	if err != nil {
		slog.Error("failed to get note by id", "noteId", id, "error", err)
		return nil, err
	}

	slog.Info("got note by id successfully", "noteId", id)

	return note, nil
}

func (s *NoteService) GetAllNotes() ([]*models.Note, error) {
	slog.Info("getting all notes")

	notes, err := s.repo.GetAllNotes()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get notes: %w", err)
		slog.Error("failed to get all notes", "error", wrappedErr)
		return nil, wrappedErr
	}

	slog.Info("got all notes successfully", "noteCount", len(notes))

	return notes, nil
}

func (s *NoteService) UpdateNote(id int, req *models.UpdateNoteRequest) (*models.Note, error) {
	slog.Info("updating note", "noteId", id)

	if id <= 0 {
		err := fmt.Errorf("invalid note ID: %d", id)
		slog.Error("failed to update note", "noteId", id, "error", err)
		return nil, err
	}

	if err := s.validateUpdateRequest(req); err != nil {
		slog.Error("failed to update note", "noteId", id, "error", err)
		return nil, err
	}

	updates := make(map[string]any)

	if req.Content != nil {
		trimmedContent := strings.TrimSpace(*req.Content)
		if trimmedContent == "" {
			err := fmt.Errorf("content cannot be empty")
			slog.Error("failed to update note", "noteId", id, "error", err)
			return nil, err
		}
		updates["content"] = trimmedContent
	}

	if len(updates) == 0 {
		err := fmt.Errorf("no valid updates provided")
		slog.Error("failed to update note", "noteId", id, "error", err)
		return nil, err
	}

	if err := s.repo.UpdateNote(id, updates); err != nil {
		slog.Error("failed to update note", "noteId", id, "error", err)
		return nil, err
	}

	note, err := s.repo.GetNoteByID(id)
	if err != nil {
		slog.Error("failed to update note", "noteId", id, "error", err)
		return nil, err
	}

	slog.Info("updated note successfully", "noteId", id)

	return note, nil
}

func (s *NoteService) DeleteNote(id int) error {
	slog.Info("deleting note", "noteId", id)

	if id <= 0 {
		err := fmt.Errorf("invalid note ID: %d", id)
		slog.Error("failed to delete note", "noteId", id, "error", err)
		return err
	}

	if err := s.repo.DeleteNote(id); err != nil {
		slog.Error("failed to delete note", "noteId", id, "error", err)
		return err
	}

	slog.Info("deleted note successfully", "noteId", id)

	return nil
}

func (s *NoteService) validateCreateRequest(req *models.CreateNoteRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return fmt.Errorf("content is required")
	}

	return nil
}

func (s *NoteService) validateUpdateRequest(req *models.UpdateNoteRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Content == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return nil
}
