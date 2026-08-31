package models

// ChatMessage represents a single turn in a quiz conversation between the
// user and the AI assistant.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateQuizRequest is the payload for POST /quiz/generate.
//
// ConversationID is optional when Messages is empty (starting a new
// conversation) but required and validated against a previously-issued id
// when Messages is non-empty.
type GenerateQuizRequest struct {
	ConversationID string        `json:"conversationId"`
	Messages       []ChatMessage `json:"messages"`
}

// GenerateQuizResponse is the response for POST /quiz/generate. Messages is
// the full conversation so far: the request's Messages with the newly
// generated assistant message appended.
type GenerateQuizResponse struct {
	ConversationID string        `json:"conversationId"`
	Messages       []ChatMessage `json:"messages"`
}
