// Package dto provides all data structures used by the api
package dto

import (
	"github.com/go-playground/validator/v10"
)

// Validate is a validator instance for JSON transferable objects
var Validate = validator.New(validator.WithRequiredStructEnabled())
