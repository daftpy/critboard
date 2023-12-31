package queryFeedback

import (
	"context"
	"critboard-backend/database/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Create(
	ctx context.Context,
	db *pgxpool.Pool, parentID string,
	feedbackText string,
	user common.User,
) (common.Feedback, error) {
	var feedback common.Feedback
	var commentID string

	// Create the commentable entry first
	err := db.QueryRow(ctx, `
		INSERT INTO commentables DEFAULT VALUES 
		RETURNING id
	`).Scan(&commentID)

	if err != nil {
		return common.Feedback{}, err
	}

	// Then create the feedback entry
	err = db.QueryRow(ctx, `
		INSERT INTO feedback (commentable_id, parent_commentable_id, feedback_text, author)
		VALUES ($1, $2, $3, $4)
		RETURNING commentable_id, feedback_text, created_at
	`, commentID, parentID, feedbackText, user.ID).Scan(
		&feedback.CommentID, &feedback.FeedbackText, &feedback.CreatedAt,
	)

	feedback.Author = user

	if err != nil {
		return common.Feedback{}, err
	}

	return feedback, nil
}
