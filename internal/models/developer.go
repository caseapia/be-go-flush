package models

import "github.com/caseapia/goproject-flush/pkg/utils/enums"

type ServiceInteractionRequest struct {
	Name   string              `json:"name"`
	Action enums.ServiceAction `json:"action"`
}
