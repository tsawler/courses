package clientmodels

import "time"

// Sample is a sample model
type Sample struct {
	ID        int
	Name      string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
