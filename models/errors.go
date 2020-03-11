package models

// IsNotFoundError returns whether an error represents a "not found" error.
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case AccountNotFoundError:
		return true
	}
	return false
}

// AccountNotFoundError represents when a user is not found.
type AccountNotFoundError struct{}

func (e AccountNotFoundError) Error() string {
	return "Account not found"
}
