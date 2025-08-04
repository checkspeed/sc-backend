package models

type CreateFeedback struct {
	Subject string `json:"subject,omitempty"`
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}
