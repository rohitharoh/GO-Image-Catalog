package system

import (
	"errors"
	"strings"
)

var (
	UnauthorisedErr       = errors.New("Unauthorized request")
	InternalServerError    = errors.New("Sorry, Something went wrong.")
	NotAnImageFile    = errors.New("Invalid File Format")

	DateFormatMismatchErr = errors.New("Please enter the date as YYYY-MM-DD (e.g. 2016-09-25)")
)

func GetErrorMessagesMap() map[error]bool {

	errorMessageMap := map[error]bool{
		InternalServerError:true,
		UnauthorisedErr:true,
		NotAnImageFile:true,
		DateFormatMismatchErr:true,
	}
	return errorMessageMap
}

func IsFunctionalError(err error) bool {
	errorMessageMap := GetErrorMessagesMap()

	if errorMessageMap[err] {
		return true
	} else if strings.Contains(err.Error(), "Error:") {
		return true

	}

	return false
}

func DefaultErr(errorMsg string) error {

	errorMsg = "Error: " + errorMsg
	return errors.New(errorMsg)
}
