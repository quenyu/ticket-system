package model

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type APIErrorDetail struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type LegacyAPIError struct {
	Error APIErrorDetail `json:"error"`
}

type TicketFilter struct {
	StatusID     *int16
	AssigneeID   *string
	DepartmentID *int16
	Q            string
	Limit        int
	Offset       int
}
