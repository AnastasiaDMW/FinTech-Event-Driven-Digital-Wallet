package dto

import "errors"

var ErrInvalidFormatBirth = errors.New("invalid date format, use DD.MM.YYYY or DD/MM/YYYY")