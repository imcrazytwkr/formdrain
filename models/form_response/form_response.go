package form_response

import (
	"net/netip"
	"time"
)

// FormResponse is the persisted submission envelope (SQLite form_responses row).
type FormResponse struct {
	Id            string
	FormId        int64
	CreatedAt     time.Time
	SchemaVersion int
	ClientIP      netip.Addr
	Payload       map[string]any
}
