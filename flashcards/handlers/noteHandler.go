package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"flashcards/models"
	"flashcards/services"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	service *services.NoteService
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

func (h *NoteHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notes", h.CreateNote).Methods("POST")
	router.HandleFunc("/notes", h.GetAllNotes).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}", h.GetNoteByID).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}", h.UpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id:[0-9]+}", h.DeleteNote).Methods("DELETE")
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode create note request body", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	slog.Info("received create note request")

	note, err := h.service.CreateNote(&req)
	if err != nil {
		// Error already logged where it originated in the service layer;
		// just translate it into the appropriate HTTP response here.
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("create note request completed successfully", "noteId", note.ID)

	h.writeJSONResponse(w, http.StatusCreated, note)
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	slog.Info("received get all notes request")

	notes, err := h.service.GetAllNotes()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve notes")
		return
	}

	slog.Info("get all notes request completed successfully", "noteCount", len(notes))

	h.writeJSONResponse(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("failed to parse note id", "id", vars["id"], "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	slog.Info("received get note by id request", "noteId", id)

	note, err := h.service.GetNoteByID(id)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve note")
		}
		return
	}

	slog.Info("get note by id request completed successfully", "noteId", id)

	h.writeJSONResponse(w, http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("failed to parse note id", "id", vars["id"], "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode update note request body", "noteId", id, "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	slog.Info("received update note request", "noteId", id)

	note, err := h.service.UpdateNote(id, &req)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	slog.Info("update note request completed successfully", "noteId", id)

	h.writeJSONResponse(w, http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("failed to parse note id", "id", vars["id"], "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	slog.Info("received delete note request", "noteId", id)

	err = h.service.DeleteNote(id)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete note")
		}
		return
	}

	slog.Info("delete note request completed successfully", "noteId", id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *NoteHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
