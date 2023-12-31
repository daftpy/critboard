package queryFeedback

import (
	"context"
	"critboard-backend/database/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetByParentID(ctx context.Context, db *pgxpool.Pool, parentID string) ([]*common.Feedback, error) {
	var feedbacks []*common.Feedback

	rows, err := db.Query(ctx, `
        SELECT 
            f1.commentable_id, 
            CASE 
                WHEN f1.deleted THEN 'removed'
                ELSE f1.feedback_text
            END AS feedback_text, 
            f1.created_at, 
            COUNT(f2.commentable_id) as reply_count,
			f1.deleted,
			u.id AS user_id,
            u.twitch_id,
            u.username,
            u.email
        FROM feedback f1
        LEFT JOIN feedback f2 ON f1.commentable_id = f2.parent_commentable_id
        LEFT JOIN users u ON f1.author = u.id
        WHERE f1.parent_commentable_id = $1
        GROUP BY f1.commentable_id, f1.feedback_text, f1.deleted, u.id
    `, parentID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var feedback common.Feedback
		var reply_count int

		if err := rows.Scan(
			&feedback.CommentID,
			&feedback.FeedbackText,
			&feedback.CreatedAt,
			&reply_count,
			&feedback.Removed,
			&feedback.Author.ID,
			&feedback.Author.TwitchID,
			&feedback.Author.Username,
			&feedback.Author.Email,
		); err != nil {
			return nil, err
		}

		feedback.Replies = reply_count
		feedbacks = append(feedbacks, &feedback)
	}

	return feedbacks, nil
}
