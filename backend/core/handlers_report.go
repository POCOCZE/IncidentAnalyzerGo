package core

import (
	"fmt"
	"log"
	"net/http"
)

func getReportHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		severity := r.URL.Query().Get("severity")
		service := r.URL.Query().Get("service")
		id := r.URL.Query().Get("id")
		w.Header().Add("Content-Type", "application/json")
		incidents, _ := store.GetAll()
		report, err := BuildReport(incidents)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("%s", err)
		} else {
			if severity != "" {
				_, exist := report.BySeverity[severity]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					encodeJSON(w, map[string]string{"error": fmt.Sprintf("severity %s does not exist", severity)}, "")
					log.Printf("ERR: Severity %s does not exist", service)
				} else {
					w.WriteHeader(http.StatusOK)
					encodeJSON(w, report.BySeverity[severity], "")
					log.Printf("INFO: Requested %s severity", severity)
				}
			} else if service != "" {
				_, exist := report.ByServices[service]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					encodeJSON(w, map[string]string{"error": fmt.Sprintf("service %s does not exist", service)}, "")
					log.Printf("ERR: Service %s does not exist", service)
				} else {
					w.WriteHeader(http.StatusOK)
					encodeJSON(w, report.ByServices[service], "")
					log.Printf("INFO: Requested %s severity", service)
				}
			} else if id != "" {
				_, exist := report.ByID[id]
				if !exist {
					w.WriteHeader(http.StatusBadRequest)
					encodeJSON(w, map[string]string{"error": fmt.Sprintf("service %s does not exist", id)}, "")
					log.Printf("ERR: Service %s does not exist", id)
				} else {
					w.WriteHeader(http.StatusOK)
					encodeJSON(w, report.ByID[id], "")
					log.Printf("INFO: Requested %s severity", id)
				}
			} else {
				w.WriteHeader(http.StatusOK)
				encodeJSON(w, report, "")
				log.Printf("INFO: Requested all incidents")
			}
		}
	}
}