package enums

type ServiceAction string

const (
	Enable  ServiceAction = "enable"
	Disable ServiceAction = "disable"
	Status  ServiceAction = "status"
)
