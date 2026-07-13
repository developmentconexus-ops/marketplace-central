package domain

import "errors"

var ErrInvalidInternalProductID = errors.New("PRODUCT_LINKS_INTERNAL_PRODUCT_INVALID")

func ValidateInternalProductID(id int) error {
	if id <= 0 {
		return ErrInvalidInternalProductID
	}
	return nil
}
