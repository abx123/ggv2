package presenter

// Guest represents a Guest object
type Guest struct {
	ID                 int64  `json:"id,omitempty"`
	Name               string `json:"name"`
	TableID            int64  `json:"tableid,omitempty"`
	AccompanyingGuests int64  `json:"accompanying_guests"`
	ArrivalTime        string `json:"arrived_time,omitempty"`
}
