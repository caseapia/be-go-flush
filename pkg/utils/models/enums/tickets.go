package enums

// TicketStatus defines the current status of the ticket
type TicketStatus string

// TicketPriority defines the current priority of the ticket. Tickets with higher priority will be higher in tickets list
type TicketPriority string

// TicketAction defines actions performed by users involved in the ticket
type TicketAction string

const (
	// Open indicates the ticket is opened and staff member answered
	TicketStatusOpen TicketStatus = "Waiting for user"
	// Closed indicates the ticket was resolved
	TicketStatusClosed TicketStatus = "Closed"
	// Pending indicates user answered in the ticket and waiting for staff member response
	TicketStatusPending TicketStatus = "Waiting for staff"
)

const (
	TicketPriorityLow      TicketPriority = "Low"
	TicketPriorityMedium   TicketPriority = "Medium"
	TicketPriorityHigh     TicketPriority = "High"
	TicketPriorityCritical TicketPriority = "Critical"
)

const (
	TicketActionClose          TicketAction = "has closed this ticket"
	TicketActionReopen         TicketAction = "has reopened this ticket"
	TicketActionHandled        TicketAction = "has handled this ticket"
	TicketActionUnassign       TicketAction = "has unassigned himself from this ticket"
	TicketActionChangeCategory TicketAction = "has changed category of this ticket to %s"
)
