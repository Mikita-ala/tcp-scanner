package request

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func GetScanResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		w.Header().Set("Content-Type", "application/json")

		slog.Info("Getting scan result for " + id)

		value, ok := tasks.Load(id)
		if !ok {
			slog.Error("Could not find task " + id)
			http.NotFound(w, r)
			return
		}
		slog.Info("Got task " + id)
		result := value.(ScanResult)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(result)
	}
}
