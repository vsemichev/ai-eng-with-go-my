package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"flashcards/models"
	"flashcards/services"

	"github.com/gorilla/mux"
)

type QuizHandler struct {
	service *services.QuizService
}

func NewQuizHandler(service *services.QuizService) *QuizHandler {
	return &QuizHandler{service: service}
}

func (h *QuizHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/quiz/generate", h.GenerateQuiz).Methods("POST")
}

func (h *QuizHandler) GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode quiz generate request body", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	slog.Info("received quiz generate request", "conversationId", req.ConversationID, "messageCount", len(req.Messages))

	resp, err := h.service.GenerateNext(r.Context(), &req)
	if err != nil {
		// Error already logged where it originated in the service layer;
		// just translate it into the appropriate HTTP response here.
		if errors.Is(err, services.ErrInvalidQuizRequest) {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	slog.Info("quiz generate request completed successfully", "conversationId", resp.ConversationID, "messageCount", len(resp.Messages))

	h.writeJSONResponse(w, http.StatusOK, resp)
}

func (h *QuizHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *QuizHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
