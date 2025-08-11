// Package dto contains all data transferable objects
package dto

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())
