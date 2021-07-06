package presenter

// Table represents a Table object
type Table struct {
	TableID  int64 `json:"id,omitempty"`
	Capacity int64 `json:"capacity,omitempty"`
}
