package models

import (
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
)

type CreateUserPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

func (c *CreateUserPayload) Validate() error {
	if !validator.IsStringEmpty(c.Name) ||
		!validator.IsStringEmpty(c.Email) ||
		!validator.IsStringEmpty(c.Country) ||
		!validator.IsStringEmpty(c.CountryCode) ||
		!validator.IsStringEmpty(c.Phone) ||
		!validator.ValidateStringIsCountryCode(c.CountryCode) ||
		!validator.ValidateStringIsEmail(c.Email) ||
		!validator.ValidateStringIsPhoneNumber(c.Phone) {
		return errorhandler.NewDomainError(errorhandler.ErrUserRequestPayloadBadRequest, errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequest), nil)
	}
	return nil
}
