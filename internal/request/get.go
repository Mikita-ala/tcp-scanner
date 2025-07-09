package request

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetScanResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		w.Header().Set("Content-Type", "application/json")

		value, ok := tasks.Load(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		result := value.(ScanResult)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(result)
	}
}
