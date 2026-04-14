package store

import (
	"database/sql"
	"errors"
)

type NotificationRepository struct {
	store *PostgresDB
}

func NewNotificationRepository(s *PostgresDB) *NotificationRepository {
	return &NotificationRepository{store: s}
}

func (r *NotificationRepository) AddUser(userID int, email string) error {
	_, err := r.store.DB.Exec(`
		INSERT INTO notification_users(user_id, email)
		VALUES($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET 
			email = EXCLUDED.email,
			updated_at = NOW()
	`, userID, email)

	return err
}

func (r *NotificationRepository) GetEmailByUserID(userID int) (string, error) {
	var email string

	err := r.store.DB.QueryRow(
		"SELECT email FROM notification_users WHERE user_id=$1",
		userID,
	).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return email, nil
}
