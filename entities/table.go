package entities

// Table represents a Table object
type Table struct {
	TableID           int64 `json:"id,omitempty" db:"tableid"`
	Capacity          int64 `json:"capacity,omitempty" db:"capacity"`
	AvailableCapacity int64 `json:"available_capacity,omitempty" db:"acapacity"`
	PlannedCapacity   int64 `json:"available_planned_capacity,omitempty" db:"pcapacity"`
	Version           int64 `json:"version,omitempty" db:"version"`
}
