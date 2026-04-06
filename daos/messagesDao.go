package daos

import (
	"context"
	"talky-space-be/models"
)

func (d *PgxDao) CreateMessage(ctx context.Context, req models.Messages) (*models.Messages, error) {

	query := `INSERT INTO messages (id, chatroom_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`

	if _, err := d.pool.Exec(ctx, query, req.Id, req.ChatroomID, req.UserId, req.Content, req.CreatedAt); err != nil {
		return nil, err
	}
	return &req, nil
}

func (d *PgxDao) GetMessagesByChatroomID(ctx context.Context, chatroomID string) ([]models.Messages, error) {
	query := `SELECT id, chatroom_id::text, user_id, content, created_at FROM messages WHERE chatroom_id = $1 ORDER BY created_at`

	var messages []models.Messages

	rows, err := d.pool.Query(ctx, query, chatroomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message models.Messages
		if err := rows.Scan(&message.Id, &message.ChatroomID, &message.UserId, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (d *PgxDao) DeleteMessagesByChatroomID(ctx context.Context, chatroomID string) error {
	query := `DELETE FROM messages WHERE chatroom_id = $1`

	if _, err := d.pool.Exec(ctx, query, chatroomID); err != nil {
		return err
	}
	return nil
}
