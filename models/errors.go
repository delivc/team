package models

// IsNotFoundError returns whether an error represents a "not found" error.
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case AccountNotFoundError:
		return true
	}
	return false
}

// PermissionNotFoundError represents when a user is not found.
type PermissionNotFoundError struct{}

func (e PermissionNotFoundError) Error() string {
	return "Permission not found"
}

// AccountNotFoundError represents when a user is not found.
type AccountNotFoundError struct{}

func (e AccountNotFoundError) Error() string {
	return "Account not found"
}

// RoleNotFoundError represents when role is not found.
type RoleNotFoundError struct{}

func (e RoleNotFoundError) Error() string {
	return "Role not found"
}
