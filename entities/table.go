package entities

// Table represents a Table object
type Table struct {
	TableID           int64 `db:"id"`
	Capacity          int64 `db:"capacity"`
	AvailableCapacity int64 `db:"acapacity"`
	PlannedCapacity   int64 `db:"pcapacity"`
	Version           int64 `db:"version"`
}
