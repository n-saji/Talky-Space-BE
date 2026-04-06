package daos

import (
	"context"
	"talky-space-be/models"
)

func (d *PgxDao) SaveSession(ctx context.Context, session *models.Sessions) error {
	query := `INSERT INTO sessions (id, user_id, refresh_token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`

	if _, err := d.pool.Exec(ctx, query, session.Id, session.UserId, session.RefreshToken, session.ExpiresAt, session.CreatedAt); err != nil {
		return err
	}
	return nil
}
