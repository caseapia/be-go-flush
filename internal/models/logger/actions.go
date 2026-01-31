package loggermodel

type LoggerAction string

const (
	// ! Admin actions
	Ban          LoggerAction = "has banned"
	Unban        LoggerAction = "has unbanned"
	Create       LoggerAction = "has created"
	SoftDelete   LoggerAction = "has soft-deleted"
	RestoreUser  LoggerAction = "has restored"
	SetAdmin     LoggerAction = "has set admin perm"
	SetDeveloper LoggerAction = "has set developer perm"

	// ! Searches
	SearchByUsername LoggerAction = "searched by username"
	SearchByUserID   LoggerAction = "searched by user ID"
	SearchByAllUsers LoggerAction = "searched all users"
	SearchLogs       LoggerAction = "searched logs"
)
