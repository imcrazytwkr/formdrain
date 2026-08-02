package brevo

import "golang.org/x/time/rate"

// @NOTE: just to be safe it's better to share a single limiter for all
// BrevoNotifier instances
var rateLimiter = rate.NewLimiter(rate.Limit(httpRateLimit), httpMaxBurst)
