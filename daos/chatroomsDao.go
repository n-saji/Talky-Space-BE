package daos

import (
	"context"
	"errors"
	"talky-space-be/models"

	"github.com/jackc/pgx/v5"
)

func (d *PgxDao) CreateChatroom(ctx context.Context, chatroom *models.Chatroom) error {

	query := `INSERT INTO chatrooms (id, name, description, is_group, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err := d.pool.Exec(ctx, query, chatroom.ID, chatroom.Name, chatroom.Description, chatroom.IsGroup, chatroom.CreatedBy, chatroom.CreatedAt, chatroom.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) GetChatroomByID(ctx context.Context, id string) (*models.Chatroom, error) {
	query := `SELECT id, name, description, is_group, COALESCE(created_by::text, ''), created_at, updated_at
		FROM chatrooms WHERE id = $1 LIMIT 1`

	var chatroom models.Chatroom
	if err := d.pool.QueryRow(ctx, query, id).Scan(&chatroom.ID, &chatroom.Name, &chatroom.Description, &chatroom.IsGroup, &chatroom.CreatedBy, &chatroom.CreatedAt, &chatroom.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}
	return &chatroom, nil
}

func (d *PgxDao) UpdateChatroom(ctx context.Context, chatroom *models.Chatroom) error {
	query := `UPDATE chatrooms SET name = $1, description = $2, is_group = $3, created_by = $4, updated_at = $5 WHERE id = $6`

	if _, err := d.pool.Exec(ctx, query, chatroom.Name, chatroom.Description, chatroom.IsGroup, chatroom.CreatedBy, chatroom.UpdatedAt, chatroom.ID); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) DeleteChatroom(ctx context.Context, id string) error {
	query := `DELETE FROM chatrooms WHERE id = $1`

	if _, err := d.pool.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) CheckChatroomExistForSenderReceiver(ctx context.Context, senderID string, receiverID string) (*models.Chatroom, error) {
	query := `SELECT c.id, c.name, c.description, c.is_group, COALESCE(c.created_by::text, ''), c.created_at, c.updated_at
		FROM chatrooms c
		JOIN chatroom_members cm1 ON c.id = cm1.chatroom_id
		JOIN chatroom_members cm2 ON c.id = cm2.chatroom_id
		WHERE c.is_group = false AND cm1.user_id = $1 AND cm2.user_id = $2
		LIMIT 1`

	var chatroom models.Chatroom
	err := d.pool.QueryRow(ctx, query, senderID, receiverID).Scan(&chatroom.ID, &chatroom.Name, &chatroom.Description, &chatroom.IsGroup, &chatroom.CreatedBy, &chatroom.CreatedAt, &chatroom.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &chatroom, nil
}

func (d *PgxDao) FetchChatroomsForUserId(ctx context.Context, userId string) ([]*models.Chatroom, error) {
	query := `SELECT c.* from chatrooms c
				JOIN chatroom_members cm ON cm.chatroom_id = c.id
				JOIN users u ON u.id = cm.user_id
				WHERE cm.user_id = $1`
	rows, err := d.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatrooms []*models.Chatroom
	for rows.Next() {
		var chatroom models.Chatroom
		if err := rows.Scan(&chatroom.ID, &chatroom.Name, &chatroom.Description, &chatroom.IsGroup, &chatroom.CreatedBy, &chatroom.CreatedAt, &chatroom.UpdatedAt); err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, &chatroom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatrooms, nil
}
