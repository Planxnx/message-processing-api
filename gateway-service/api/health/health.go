package health

import "net/http"

func HealthHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health  OK"))
	return
}
