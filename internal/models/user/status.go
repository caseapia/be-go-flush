package usermodel

import (
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

type UserStatus int

const (
	Player UserStatus = iota
	Tester
	TrialAdmin
	AdminLvl1
	AdminLvl2
	SeniorAdmin
	LeadAdmin
	Manager
	TotalStatuses
)

type DeveloperStatus int

const (
	NotDeveloper DeveloperStatus = iota
	Developer
	LeadDeveloper
	WebDeveloper
)

func (s UserStatus) String() string {
	switch s {
	case Player:
		return "Player"
	case Tester:
		return "Tester"
	case TrialAdmin:
		return "Trial Admin"
	case AdminLvl1:
		return "Admin Lvl. 1"
	case AdminLvl2:
		return "Admin Lvl. 2"
	case SeniorAdmin:
		return "Senior Admin"
	case LeadAdmin:
		return "Lead Admin"
	case Manager:
		return "Manager"
	default:
		return "Unknown"
	}
}

func (s DeveloperStatus) String() string {
	switch s {
	case NotDeveloper:
		return "Not Developer"
	case Developer:
		return "Developer"
	case LeadDeveloper:
		return "Lead Developer"
	case WebDeveloper:
		return "Web Developer"
	default:
		return "Unknown"
	}
}

func ParseUserStatus(s string) (UserStatus, error) {
	switch s {
	case "0":
		return Player, nil
	case "1":
		return Tester, nil
	case "2":
		return TrialAdmin, nil
	case "3":
		return AdminLvl1, nil
	case "4":
		return AdminLvl2, nil
	case "5":
		return SeniorAdmin, nil
	case "6":
		return LeadAdmin, nil
	case "7":
		return Manager, nil
	default:
		return 0, UserError.UserInvalidStatus()
	}
}

func ParseDeveloperStatus(s string) (DeveloperStatus, error) {
	switch s {
	case "0":
		return NotDeveloper, nil
	case "1":
		return Developer, nil
	case "2":
		return LeadDeveloper, nil
	case "3":
		return WebDeveloper, nil
	default:
		return 0, UserError.UserInvalidStatus()
	}
}
