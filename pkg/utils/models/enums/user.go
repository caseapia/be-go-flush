package enums

// Defines the current state of a user account
type UserStatus int

const (
	// Indicates the account is created but email is unconfirmed
	UserStatusNotVerified UserStatus = iota
	// Indicates the user has full access to the all site features
	UserStatusActive
	// Indicates the account is temporarily suspended by the user or the system
	UserStatusDisabled
	// Indicates the account is marked as deleted (soft delete)
	UserStatusDeleted
	// Indicates the user must update their password before proceeding
	UserStatusRequiresPasswordChange
)

// Defines the current state of ban
type BanStatus int

const (
	// Indicates that restriction is now active
	BanActive BanStatus = iota
	// Indicates that ban was removed by an admin
	BanRemoved
	// Indicates that ban was removed due to expiration
	BanExpired
)
