package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"flashcards/models"

	"github.com/tmc/langchaingo/llms"
)

const roleAssistant = "assistant"

// ErrInvalidQuizRequest is returned for client input errors (missing/unknown
// conversationId, nil request), so handlers can map it to a 400 response as
// opposed to a 500 for downstream (e.g. LLM/DB) failures.
var ErrInvalidQuizRequest = errors.New("invalid quiz request")

// QuizService drives the AI quiz conversation. It keeps track of issued
// conversation ids in memory so it can validate that a client is continuing
// a conversation it actually started, and it uses NoteService to ground the
// generated questions in the user's own notes.
type QuizService struct {
	mu            sync.Mutex
	conversations map[string]struct{}
	noteService   *NoteService
	llm           llms.Model
}

func NewQuizService(noteService *NoteService, llm llms.Model) *QuizService {
	return &QuizService{
		conversations: make(map[string]struct{}),
		noteService:   noteService,
		llm:           llm,
	}
}

// GenerateNext handles both starting a new conversation (empty Messages) and
// continuing an existing one (non-empty Messages + valid ConversationID).
func (s *QuizService) GenerateNext(ctx context.Context, req *models.GenerateQuizRequest) (*models.GenerateQuizResponse, error) {
	if req == nil {
		err := fmt.Errorf("request cannot be nil: %w", ErrInvalidQuizRequest)
		slog.Error("failed to generate next quiz message", "error", err)
		return nil, err
	}

	slog.Info("generating next quiz message", "conversationId", req.ConversationID, "messageCount", len(req.Messages))

	if len(req.Messages) == 0 {
		conversationID, err := s.startConversation()
		if err != nil {
			return nil, err
		}

		nextMessage, err := s.generateNextMessage(ctx, req.Messages)
		if err != nil {
			return nil, err
		}

		slog.Info("generated next quiz message successfully", "conversationId", conversationID, "messageCount", 1)

		return &models.GenerateQuizResponse{
			ConversationID: conversationID,
			Messages:       []models.ChatMessage{nextMessage},
		}, nil
	}

	if req.ConversationID == "" {
		err := fmt.Errorf("conversationId is required when messages is not empty: %w", ErrInvalidQuizRequest)
		slog.Error("failed to generate next quiz message", "error", err)
		return nil, err
	}

	if !s.isKnownConversation(req.ConversationID) {
		err := fmt.Errorf("unknown conversationId: %s: %w", req.ConversationID, ErrInvalidQuizRequest)
		slog.Error("failed to generate next quiz message", "conversationId", req.ConversationID, "error", err)
		return nil, err
	}

	nextMessage, err := s.generateNextMessage(ctx, req.Messages)
	if err != nil {
		return nil, err
	}

	messages := make([]models.ChatMessage, 0, len(req.Messages)+1)
	messages = append(messages, req.Messages...)
	messages = append(messages, nextMessage)

	slog.Info("generated next quiz message successfully", "conversationId", req.ConversationID, "messageCount", len(messages))

	return &models.GenerateQuizResponse{
		ConversationID: req.ConversationID,
		Messages:       messages,
	}, nil
}

func (s *QuizService) startConversation() (string, error) {
	id, err := newConversationID()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to generate conversationId: %w", err)
		slog.Error("failed to start new quiz conversation", "error", wrappedErr)
		return "", wrappedErr
	}

	s.mu.Lock()
	s.conversations[id] = struct{}{}
	s.mu.Unlock()

	slog.Info("started new quiz conversation", "conversationId", id)

	return id, nil
}

func (s *QuizService) isKnownConversation(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.conversations[id]
	return ok
}

// generateNextMessage grounds a new quiz question in the user's notes (or a
// hardcoded fallback set if there are none) plus the conversation so far,
// and asks the LLM to produce the next assistant message.
func (s *QuizService) generateNextMessage(ctx context.Context, history []models.ChatMessage) (models.ChatMessage, error) {
	slog.Info("generating quiz question via LLM", "historyLength", len(history))

	notesText, err := s.notesText()
	if err != nil {
		// Error already logged in notesText where it originated.
		return models.ChatMessage{}, fmt.Errorf("failed to load notes: %w", err)
	}

	prompt := fmt.Sprintf(quizPromptTemplate, quizSystemPrompt, notesText, historyText(history))

	completion, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt, llms.WithTemperature(llmTemperature))
	if err != nil {
		wrappedErr := fmt.Errorf("failed to generate quiz question: %w", err)
		slog.Error("failed to generate quiz question via LLM", "error", wrappedErr)
		return models.ChatMessage{}, wrappedErr
	}

	message := models.ChatMessage{Role: roleAssistant, Content: strings.TrimSpace(completion)}

	slog.Info("generated quiz question via LLM successfully", "contentLength", len(message.Content))

	return message, nil
}

// notesText fetches all notes from the database and joins their content
// into a single block of text. If there are no notes (or the note service
// isn't available), it falls back to a small set of predefined notes so a
// quiz can still be generated.
func (s *QuizService) notesText() (string, error) {
	slog.Info("fetching notes for quiz generation")

	var contents []string

	if s.noteService != nil {
		notes, err := s.noteService.GetAllNotes()
		if err != nil {
			wrappedErr := fmt.Errorf("failed to fetch notes: %w", err)
			slog.Error("failed to fetch notes for quiz generation", "error", wrappedErr)
			return "", wrappedErr
		}
		for _, note := range notes {
			contents = append(contents, note.Content)
		}
	}

	usedFallback := len(contents) == 0
	if usedFallback {
		contents = quizFallbackNotes
	}

	slog.Info("fetched notes for quiz generation successfully", "noteCount", len(contents), "usedFallback", usedFallback)

	return strings.Join(contents, "\n- "), nil
}

// historyText serializes the conversation history into flat "role: content"
// lines for inclusion in the single LLM prompt string.
func historyText(history []models.ChatMessage) string {
	if len(history) == 0 {
		return noHistoryPlaceholder
	}

	lines := make([]string, 0, len(history))
	for _, msg := range history {
		lines = append(lines, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
	}

	return strings.Join(lines, "\n")
}

func newConversationID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
