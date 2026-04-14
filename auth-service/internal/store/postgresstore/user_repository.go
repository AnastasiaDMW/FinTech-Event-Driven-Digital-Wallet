package postgresstore

import (
	"github.com/AnastasiaDMW/auth-service/internal/model"
	"github.com/AnastasiaDMW/auth-service/internal/store"
	"github.com/lib/pq"
)

type UserRepository struct {
	store *PostgresDB
}

func NewUserRepository(s *PostgresDB) *UserRepository {
	return &UserRepository{store: s}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	row := r.store.DB.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email=$1",
		email,
	)

	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	return u, err
}

func (r *UserRepository) CreateUser(u *model.User) (int, error) {
	id := 0
	err := r.store.DB.QueryRow(
		"INSERT INTO users(email, encrypted_password) VALUES($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, store.ErrEmailExists
			}
		}
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) UpdateEmail(userID int, email string) error {
	_, err := r.store.DB.Exec(
		"UPDATE users SET email = $1, email_verified = false, updated_at=NOW() WHERE id = $2",
		email,
		userID,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return store.ErrEmailExists
			}
		}
		return err
	}

	return nil
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	row := r.store.DB.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id=$1",
		id,
	)

	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) VerifyEmail(userID int) error {
	res, err := r.store.DB.Exec(
		"UPDATE users SET email_verified=true, updated_at=NOW() WHERE id=$1",
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return store.ErrUserNotFound
	}

	return nil
}
