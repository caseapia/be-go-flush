package enums

type LoggerType string
type LoggerAction string

const (
	StaffCommonLogger     = "staffCommon"
	StaffPunishmentLogger = "staffPunish"
	AdminTicketLogger     = "adminTickets"
	AdminAuthLogger       = "adminAuth"
)

const (
	// ! Punishments
	Ban   LoggerAction = "has banned"
	Unban LoggerAction = "has unbanned"

	// ! Common actions
	CreateRank             LoggerAction = "has created rank"
	SearchByUsername       LoggerAction = "searched by username"
	SearchByUserID         LoggerAction = "searched by user ID"
	SearchLogs             LoggerAction = "searched logs"
	SetStaffRank           LoggerAction = "has set admin rank"
	SetDeveloperRank       LoggerAction = "has set developer rank"
	RestoreUser            LoggerAction = "has restored"
	Create                 LoggerAction = "has created user"
	ChangeFlags            LoggerAction = "has changed flags"
	DeleteRank             LoggerAction = "has delete rank"
	SoftDelete             LoggerAction = "has soft-deleted"
	HardDelete             LoggerAction = "has hard-deleted"
	TriedToDeleteManager   LoggerAction = "has tried to delete manager's account and action has stopped"
	CreateInvite           LoggerAction = "has created invite code"
	DeleteInvite           LoggerAction = "has deleted invite code"
	ChangeUserNickname     LoggerAction = "has changed user's nickname"
	ChangeUserPassword     LoggerAction = "has changed user's password"
	ChangeUserEmail        LoggerAction = "has changed user's email"
	ChangeUserStatus       LoggerAction = "has changed user's status"
	ResetUserSensetiveData LoggerAction = "has reset user IPs and last seen info"
	EditRank               LoggerAction = "has edited rank"
	LookupNotifications    LoggerAction = "lookup user notifications"
	SendNotification       LoggerAction = "send notify"
	DeleteNotification     LoggerAction = "has deleted notification"
	AssignedToTicket       LoggerAction = "has assigned to the ticket"
	UnassignFromTicket     LoggerAction = "has unassign himself from ticket"
	CloseTicket            LoggerAction = "has closed ticket"
	ChangeTicketCategory   LoggerAction = "has changed ticket category"
	DeleteTicket           LoggerAction = "has deleted a ticket"
	RevealSensetiveData    LoggerAction = "revealed sensetive data"
	UserRegister           LoggerAction = "just registered"
	UserLogin              LoggerAction = "just log in"
	UserLogout             LoggerAction = "just log out"
	CreateBadge            LoggerAction = "created a badge"
	EditBadge              LoggerAction = "edited a badge"
	DeleteBadge            LoggerAction = "deleted a badge"
	AwardUser              LoggerAction = "has awarded"
	LinkDiscord            LoggerAction = "linked his discord"
	UnlinkDiscord          LoggerAction = "unlinked his discord"
	ForceUnlinkUserDiscord LoggerAction = "force unlinked discord"
	CheckedSessions        LoggerAction = "just checked sessions"
	SetDonatePoints        LoggerAction = "has set donate points"
	BoughtForPoints        LoggerAction = "has bought for points"
	EnableService          LoggerAction = "has enabled service"
	DisableService         LoggerAction = "has disable service"
)
