package ports

import "time"

type OrderListQuery struct {
	InstallationID string
	Limit          int
	Cursor         OrderCursor
	Status         string
	DateFrom       *time.Time
	DateTo         *time.Time
	Q              string
}
