package main

import (
	"fmt"
	"log"
	"net/http"

	"flashcards/config"
	"flashcards/db"
	"flashcards/handlers"
	"flashcards/services"

	"github.com/gorilla/mux"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	todoRepo, err := db.NewPostgresTodoRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer todoRepo.Close()

	noteRepo, err := db.NewPostgresNoteRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer noteRepo.Close()

	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	llm, err := openai.New(
		openai.WithModel(services.LlmModel),
		openai.WithToken(cfg.OpenAIAPIKey),
	)
	if err != nil {
		log.Fatalf("Failed to initialize OpenAI client: %v", err)
	}

	quizService := services.NewQuizService(noteService, llm)
	quizHandler := handlers.NewQuizHandler(quizService)

	router := mux.NewRouter()

	router.Use(corsMiddleware)
	router.Use(jsonMiddleware)

	todoHandler.RegisterRoutes(router)
	noteHandler.RegisterRoutes(router)
	quizHandler.RegisterRoutes(router)

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// gorilla/mux returns a bare 405 (bypassing corsMiddleware) for methods
	// with no matching route. Browsers send an OPTIONS preflight before
	// cross-origin requests that carry a JSON Content-Type, so without this
	// catch-all those preflights fail and the browser blocks the real
	// request. Matching OPTIONS here routes it through corsMiddleware,
	// which sets the CORS headers and short-circuits with 200.
	router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := ":" + cfg.Port
	fmt.Printf("Server starting on port %s\n", cfg.Port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
