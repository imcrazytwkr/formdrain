package health

import "net/http"

func (r *healthRouter) HandleLive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
