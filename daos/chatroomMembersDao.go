package daos

import (
	"context"
	"talky-space-be/models"
)

func (d *PgxDao) CreateChatroomMember(ctx context.Context, chatroomMember *models.ChatroomMember) error {
	query := `INSERT INTO chatroom_members (id,chatroom_id, user_id, joined_at) VALUES ($1, $2, $3, $4)`

	if _, err := d.pool.Exec(ctx, query, chatroomMember.Id, chatroomMember.ChatroomID, chatroomMember.UserID, chatroomMember.JoinedAt); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) GetChatroomMembersByChatroomID(ctx context.Context, chatroomID string) ([]models.ChatroomMember, error) {
	var members []models.ChatroomMember
	query := `SELECT id, chatroom_id, user_id, joined_at FROM chatroom_members WHERE chatroom_id = $1`

	rows, err := d.pool.Query(ctx, query, chatroomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.ChatroomMember
		if err := rows.Scan(&member.Id, &member.ChatroomID, &member.UserID, &member.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (d *PgxDao) DeleteChatroomMember(ctx context.Context, chatroomID string, userID string) error {
	query := `DELETE FROM chatroom_members WHERE chatroom_id = $1 AND user_id = $2`

	if _, err := d.pool.Exec(ctx, query, chatroomID, userID); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) IsUserInChatroom(ctx context.Context, chatroomID string, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM chatroom_members WHERE chatroom_id = $1 AND user_id = $2`

	var count int64
	if err := d.pool.QueryRow(ctx, query, chatroomID, userID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}