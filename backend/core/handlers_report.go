package core

import (
	"fmt"
	"log"
	"net/http"
)

// Build IncidentReport and return it.
func GetReportHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		severity := r.URL.Query().Get("severity")
		service := r.URL.Query().Get("service")
		id := r.URL.Query().Get("id")
		w.Header().Add("Content-Type", "application/json")
		incidents, err := store.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("%s", err)
		}
		report, err := BuildReport(incidents)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("%s", err)
		} else {
			// Check if requested only groupped incidents by Severity. If so, return it.
			if severity != "" {
				_, exist := report.BySeverity[severity]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					EncodeJSON(w, map[string]string{"error": fmt.Sprintf("severity %s does not exist", severity)}, "")
					log.Printf("ERR: Severity %s does not exist", service)
				} else {
					w.WriteHeader(http.StatusOK)
					EncodeJSON(w, report.BySeverity[severity], "")
					log.Printf("INFO: Requested %s severity", severity)
				}
			// Check if requested only groupped incidents by ServiceName. If so, return it.
			} else if service != "" {
				_, exist := report.ByServices[service]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					EncodeJSON(w, map[string]string{"error": fmt.Sprintf("service %s does not exist", service)}, "")
					log.Printf("ERR: Service %s does not exist", service)
				} else {
					w.WriteHeader(http.StatusOK)
					EncodeJSON(w, report.ByServices[service], "")
					log.Printf("INFO: Requested %s severity", service)
				}
			// Check if requested only groupped incidents by incident ID. If so, return it.
			} else if id != "" {
				_, exist := report.ByID[id]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					EncodeJSON(w, map[string]string{"error": fmt.Sprintf("service %s does not exist", id)}, "")
					log.Printf("ERR: Service %s does not exist", id)
				} else {
					w.WriteHeader(http.StatusOK)
					EncodeJSON(w, report.ByID[id], "")
					log.Printf("INFO: Requested %s severity", id)
				}
			// Otherwise return whole report.
			} else {
				w.WriteHeader(http.StatusOK)
				EncodeJSON(w, report, "")
				log.Printf("INFO: Requested all incidents")
			}
		}
	}
}