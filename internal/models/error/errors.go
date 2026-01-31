package errormodel

type ErrorMessage string

const (
	ErrUserNotFound      ErrorMessage = "user not found"
	ErrUserBanned        ErrorMessage = "user banned"
	ErrUserAlreadyExists ErrorMessage = "user already exists"
	ErrUserNotBanned     ErrorMessage = "user is not banned"
	ErrInvalidUserName   ErrorMessage = "invalid username"
	ErrInvalidUserStatus ErrorMessage = "invalid user status"
)

const (
	ErrAdminStatusAlreadySet                ErrorMessage = "admin status is already set to the specified value"
	ErrAdminCannotChangeStatusOfDeletedUser ErrorMessage = "cannot change status of a deleted user"
)
