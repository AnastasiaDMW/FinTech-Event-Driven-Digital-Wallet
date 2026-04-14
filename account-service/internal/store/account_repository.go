package store

import (
	"log/slog"

	"github.com/AnastasiaDMW/account-service/internal/dto"
	"github.com/AnastasiaDMW/account-service/internal/model"
	"github.com/google/uuid"
)

type AccountRepository struct {
	store  *PostgresDB
	logger *slog.Logger
}

func NewAccountRepository(s *PostgresDB, logger *slog.Logger) *AccountRepository {
	return &AccountRepository{store: s, logger: logger}
}

func (r *AccountRepository) CreateUserProfile(e dto.ChangedUserEvent) error {
	_, err := r.store.DB.Exec(
		"INSERT INTO users_profile (user_id) VALUES($1)",
		e.UserID,
	)
	if err != nil {
		return err
	}
	r.logger.Debug("User added successfully")

	return nil
}

func (r *AccountRepository) GetUserProfileByUserID(userID int64) (*model.UserProfile, error) {
	var u model.UserProfile

	err := r.store.DB.QueryRow(
		`
		SELECT user_id, phone, birth_date, first_name, last_name
		FROM users_profile
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&u.UserID,
		&u.Phone,
		&u.BirthDate,
		&u.FirstName,
		&u.LastName,
	)
	if err != nil {
		return nil, err
	}
	r.logger.Debug("User profile getting", "user_id", u.UserID)

	return &u, nil
}

func (r *AccountRepository) UpdateUserProfile(userId int64, req model.UserProfile) error {

	_, err := r.store.DB.Exec(
		`
			UPDATE users_profile
			SET
				phone = COALESCE(NULLIF($1, ''), phone),
				birth_date = COALESCE($2, birth_date),
				first_name = COALESCE(NULLIF($3, ''), first_name),
				last_name = COALESCE(NULLIF($4, ''), last_name)
			WHERE user_id = $5
		`,
		req.Phone,
		req.BirthDate,
		req.FirstName,
		req.LastName,
		userId,
	)
	if err != nil {
		r.logger.Debug("Failed to update user profile", "error", err)
		return err
	}
	r.logger.Debug("User profile updated", "user_id", userId)

	return nil
}

func (r *AccountRepository) CreateAccount(userID int64) error {
	var count int

	err := r.store.DB.QueryRow(
		`SELECT COUNT(*) FROM accounts WHERE user_id = $1`,
		userID,
	).Scan(&count)
	if err != nil {
		r.logger.Debug("Failed to count accounts", "error", err)
		return err
	}

	if count >= 3 {
		r.logger.Debug("Account limit reached", "user_id", userID)
		return ErrAccountLimit
	}

	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		accountNumber := uuid.NewString()

		res, err := r.store.DB.Exec(
			`
			INSERT INTO accounts (user_id, account_number)
			VALUES ($1, $2)
			ON CONFLICT (account_number) DO NOTHING
			`,
			userID,
			accountNumber,
		)
		if err != nil {
			r.logger.Debug("Failed to create account", "error", err)
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rows > 0 {
			r.logger.Debug("Account created", "user_id", userID)
			return nil
		}

		r.logger.Debug("Account_number collision, retrying", "attempt", i+1)
	}

	return ErrGenerateUniqueAccNum
}

func (r *AccountRepository) IsValidAccountNumberTo(accountNumber string) (bool, error) {
	var exists bool
	err := r.store.DB.QueryRow(
		`
		SELECT EXISTS (
			SELECT 1
			FROM accounts
			WHERE account_number = $1
			  AND status = 'active'
		)
		`,
		accountNumber,
	).Scan(&exists)

	if err != nil {
		r.logger.Debug("Failed to check account", "error", err)
		return false, err
	}

	return exists, nil
}

func (r *AccountRepository) IsValidAccountNumberFrom(userID int64, accountNumber string) (bool, error) {
	var exists bool
	err := r.store.DB.QueryRow(
		`
		SELECT EXISTS (
			SELECT 1
			FROM accounts
			WHERE user_id=$1 AND account_number = $2
			  AND status = 'active'
		)
		`,
		userID,
		accountNumber,
	).Scan(&exists)

	if err != nil {
		r.logger.Debug("Failed to check account", "error", err)
		return false, err
	}

	return exists, nil
}

func (r *AccountRepository) GetAccountsByUserID(userID int64) ([]model.Account, error) {
	rows, err := r.store.DB.Query(
		`
		SELECT id, user_id, account_number, status
		FROM accounts
		WHERE user_id = $1
		`,
		userID,
	)
	if err != nil {
		r.logger.Debug("Failed to query accounts", "error", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account

	for rows.Next() {
		var acc model.Account

		err := rows.Scan(
			&acc.ID,
			&acc.UserID,
			&acc.AccountNumber,
			&acc.Status,
		)
		if err != nil {
			r.logger.Debug("Failed to scan account", "error", err)
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) UpdateAccountStatus(accountNumber, status string) error {
	res, err := r.store.DB.Exec(
		`
		UPDATE accounts
		SET status = $1
		WHERE account_number = $2
		`,
		status,
		accountNumber,
	)
	if err != nil {
		r.logger.Debug("Failed to update account status", "error", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrAccountNotFound
	}

	r.logger.Debug("Account status updated", "account_number", accountNumber, "status", status)

	return nil
}
