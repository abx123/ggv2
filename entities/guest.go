package entities

// Guest represents a Guest object
type Guest struct {
	ID                 int64  `db:"guestid"`
	Name               string `db:"name"`
	TableID            int64  `db:"tableid"`
	TotalGuests        int64  `db:"total_guests"`
	TotalArrivedGuests int64  `db:"total_arrived_guests"`
	ArrivalTime        string `db:"arrivaltime"`
	Version            int64  `db:"version"`
}
