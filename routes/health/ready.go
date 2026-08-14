package health

import "net/http"

func (r *healthRouter) HandleReady(w http.ResponseWriter, req *http.Request) {
	err := r.db.PingContext(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
