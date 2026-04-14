package dto

import (
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateUserProfileRequest struct {
	Phone     *string `json:"phone"`
	BirthDate *string `json:"birthDate"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

func (r *UpdateUserProfileRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Phone,
			validation.When(r.Phone != nil,
				validation.Required,
				validation.Match(regexp.MustCompile(`^(?:\+7|7|8)9\d{9}$`)).
					Error("invalid Russian phone number"),
			),
		),
		validation.Field(&r.FirstName,
			validation.When(r.FirstName != nil,
				validation.Required,
				validation.Length(1, 50),
				validation.Match(regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)).
					Error("first_name must contain only letters"),
			),
		),
		validation.Field(&r.LastName,
			validation.When(r.LastName != nil,
				validation.Required,
				validation.Length(1, 50),
				validation.Match(regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)).
					Error("last_name must contain only letters"),
			),
		),
		validation.Field(&r.BirthDate,
			validation.When(r.BirthDate != nil,
				validation.Required,
				validation.By(func(value interface{}) error {
					str := *(value.(*string))

					t, err := ParseBirthDate(str)
					if err != nil {
						return err
					}

					if t.After(time.Now()) {
						return fmt.Errorf("birth_date cannot be in the future")
					}

					if t.Before(time.Now().AddDate(-120, 0, 0)) {
						return fmt.Errorf("birth_date is too far in the past")
					}

					return nil
				}),
			),
		),
	)
}

func ParseBirthDate(s string) (time.Time, error) {
	layouts := []string{
		"02.01.2006",
		"02/01/2006",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, ErrInvalidFormatBirth
}
