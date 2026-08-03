package notification

import "time"

// Bound outbound Discord/Brevo work after the request is done; matches DefaultClient timeout.
const notificationTimeout = 120 * time.Second
