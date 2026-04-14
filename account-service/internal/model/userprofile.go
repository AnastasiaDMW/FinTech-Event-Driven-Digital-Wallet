package model

import "time"

type UserProfile struct {
	UserID    int64      `json:"user_id"`
	Phone     *string    `json:"phone,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
}
