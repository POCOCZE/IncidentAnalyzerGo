package core

import (
	"fmt"
	"log"
	"net/http"
)

// Build and return IncidentReport
func GetReportHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		severity := r.URL.Query().Get("severity")
		service := r.URL.Query().Get("service")
		// id := r.URL.Query().Get("id")
		w.Header().Add("Content-Type", "application/json")
		incidents, err := store.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("%s", err)}, "")
			log.Printf("%s", err)
			return
		}
		report, err := BuildReport(incidents)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("%s", err)}, "")
			log.Printf("%s", err)
			return
		}

		// Check if requested only groupped incidents by Severity. If so, return it.
		if severity != "" {
			_, exist := report.BySeverity[severity]
			if !exist {
				w.WriteHeader(http.StatusBadRequest)
				EncodeJSON(w, map[string]string{"message": fmt.Sprintf("severity %s does not exist", severity)}, "")
				log.Printf("ERR: Severity %s does not exist", service)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				EncodeJSON(w, report.BySeverity[severity], "")
				log.Printf("INFO: Requested %s severity", severity)
				return
			}
		// Check if requested only groupped incidents by ServiceName. If so, return it.
		} else if service != "" {
			_, exist := report.ByServices[service]
			if !exist {
				w.WriteHeader(http.StatusBadRequest)
				EncodeJSON(w, map[string]string{"message": fmt.Sprintf("service %s does not exist", service)}, "")
				log.Printf("ERR: Service %s does not exist", service)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				EncodeJSON(w, report.ByServices[service], "")
				log.Printf("INFO: Requested %s severity", service)
				return
			}
		// Otherwise return whole report.
		} else {
			w.WriteHeader(http.StatusOK)
			EncodeJSON(w, report, "")
			log.Printf("INFO: Requested all incidents")
		}
	}
}