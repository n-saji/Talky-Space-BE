package daos

import (
	"context"
	"errors"
	"talky-space-be/models"

	"github.com/jackc/pgx/v5"
)

func (d *PgxDao) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, email, phone_number, password_hash, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	if _, err := d.pool.Exec(ctx, query, user.Id, user.Username, user.Email, user.PhoneNumber, user.PasswordHash, user.AvatarURL, user.CreatedAt, user.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, phone_number, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE email = $1 LIMIT 1`

	var user models.User
	if err := d.pool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *PgxDao) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, email, phone_number, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE id = $1 LIMIT 1`

	var user models.User
	if err := d.pool.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *PgxDao) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `SELECT id, username, email, phone_number, password_hash, avatar_url, created_at, updated_at
		FROM users WHERE phone_number = $1 LIMIT 1`

	var user models.User
	if err := d.pool.QueryRow(ctx, query, phoneNumber).Scan(&user.Id, &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *PgxDao) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users
		SET username = $1, email = $2, phone_number = $3, password_hash = $4, avatar_url = $5, updated_at = $6
		WHERE id = $7`

	if _, err := d.pool.Exec(ctx, query, user.Username, user.Email, user.PhoneNumber, user.PasswordHash, user.AvatarURL, user.UpdatedAt, user.Id); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) DeleteUser(ctx context.Context, user *models.User) error {
	query := `DELETE FROM users WHERE id = $1`

	if _, err := d.pool.Exec(ctx, query, user.Id); err != nil {
		return err
	}
	return nil
}

func (d *PgxDao) LookUpUser(ctx context.Context, query string) ([]*models.User, error) {
	sql := `SELECT id, username, email, phone_number, password_hash, avatar_url, created_at, updated_at
		FROM users
		WHERE username ILIKE $1 OR email ILIKE $2 OR phone_number ILIKE $3`

	var users []*models.User
	pattern := "%" + query + "%"

	rows, err := d.pool.Query(ctx, sql, pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (d *PgxDao) GetPrivateChatroomOtherMember(ctx context.Context, chatroomID string, userId string) (*models.User, error) {
	var member models.User
	query := `SELECT users.id, users.username, users.email, users.phone_number, users.avatar_url 
	FROM users  
	JOIN chatroom_members ON chatroom_members.user_id = users.id
	WHERE chatroom_id = $1 AND user_id != $2`

	row := d.pool.QueryRow(ctx, query, chatroomID, userId)

	if err := row.Scan(&member.Id, &member.Username, &member.Email, &member.PhoneNumber, &member.AvatarURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}

	return &member, nil
}
