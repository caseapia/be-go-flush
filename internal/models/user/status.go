package usermodel

import (
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

type UserStatus int

const (
	Player      UserStatus = 0
	Tester      UserStatus = 1
	TrialAdmin  UserStatus = 2
	AdminLvl1   UserStatus = 3
	AdminLvl2   UserStatus = 4
	SeniorAdmin UserStatus = 5
	LeadAdmin   UserStatus = 6
	Manager     UserStatus = 7
)

type DeveloperStatus int

const (
	NotDeveloper  DeveloperStatus = 0
	Developer     DeveloperStatus = 1
	LeadDeveloper DeveloperStatus = 2
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
