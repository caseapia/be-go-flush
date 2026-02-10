package logger

type UserPunishment string
type CommonAction string

// ! User punishments
const (
	Ban   UserPunishment = "has banned"
	Unban UserPunishment = "has unbanned"
)

// ! Common actions
const (
	CreateRank CommonAction = "has created rank"

	SearchByUsername CommonAction = "searched by username"
	SearchByUserID   CommonAction = "searched by user ID"
	SearchByAllUsers CommonAction = "searched all users"
	SearchLogs       CommonAction = "searched logs"

	SetStaffRank     CommonAction = "has set admin perm"
	SetDeveloperRank CommonAction = "has set developer perm"
	RestoreUser      CommonAction = "has restored"
	Create           CommonAction = "has created"
	ChangeFlags      CommonAction = "has changed flags"
	DeleteRank       CommonAction = "has delete rank"

	SoftDelete           UserPunishment = "has soft-deleted"
	HardDelete           UserPunishment = "has hard-deleted"
	TriedToDeleteManager UserPunishment = "has tried to delete manager's account and action has stopped"
)
