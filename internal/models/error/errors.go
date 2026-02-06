package ErrorModel

type ErrorMessage string

const (
	ErrUserNotFound      ErrorMessage = "user not found"
	ErrUserBanned        ErrorMessage = "user banned"
	ErrUserAlreadyExists ErrorMessage = "user already exists"
	ErrUserNotBanned     ErrorMessage = "user is not banned"
	ErrInvalidUserName   ErrorMessage = "invalid username"
	ErrInvalidUserStatus ErrorMessage = "invalid user status"
	ErrReasonRequired    ErrorMessage = "reason field required"
)

const (
	ErrAdminStatusAlreadySet                ErrorMessage = "admin status is already set to the specified value"
	ErrAdminCannotChangeStatusOfDeletedUser ErrorMessage = "cannot change status of a deleted user"
	ErrAdminMaxValueExceeded                ErrorMessage = "admin max value exceeded"
	ErrAdminManagerRankCannotBeChanged      ErrorMessage = "manager rank cannot be changed"
	ErrDeveloperRankCannotBeIssued          ErrorMessage = "developer rank cannot be issued by SetStaff function"
	ErrStaffRankCannotBeIssued              ErrorMessage = "staff rank cannot be issued by SetDeveloper function"
	ErrDeletionOfManagerIsNotAllowed        ErrorMessage = "deletion of manager account is not allowed"
)
