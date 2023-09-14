package common

import "time"

type SubmissionType string

const (
	LINK SubmissionType = "LINK"
	FILE SubmissionType = "FILE"
)

type Submission struct {
	CommentID     string         `json:"commentId"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Type          SubmissionType `json:"type"`
	Link          string         `json:"link"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	FeedbackCount *int           `json:"feedbackCount,omitempty"`
}

type SubmissionPayload struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        SubmissionType `json:"type"`
	Link        string         `json:"link,omitempty"`
	File        string         `json:"file,omitempty"`
}

type Feedback struct {
	CommentID    string    `json:"commentId"`
	FeedbackText string    `json:"feedbackText"`
	CreatedAt    time.Time `json:"createdAt"`
	Replies      int       `json:"replies"`
}

type FeedbackPayload struct {
	CommentID    string `json:"commentId"`
	FeedbackText string `json:"feedbackText"`
}

type SubmissionsResponse struct {
	Submissions []Submission `json:"submissions"`
}