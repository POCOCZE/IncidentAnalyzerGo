package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Health status endpoint
	// Set content type to application/json	
	// Write {"status": "ok"} as JSON reposnse
	status := struct {
		Status map[string]string
	}{
		Status: map[string]string{
			"status": "ok",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeJSON(w, status.Status, "")
}

func getIncidentReport(store *IncidentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		report := store.getReport()

		encodeJSON(w, report, "")
		log.Printf("INFO: Requested report")
	}
}

func getGroupedIncidents(store *IncidentStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		severity := r.URL.Query().Get("severity")
		service := r.URL.Query().Get("service")
		report := store.getReport()

		if severity != "" {
			_, exist := report.BySeverity[severity]
			if !exist {
				encodeJSON(w, map[string]string{"error": fmt.Sprintf("severity %s does not exist", severity)}, "")
				log.Printf("ERR: Severity %s does not exist", service)
			} else {
				encodeJSON(w, report.BySeverity[severity], "")
				log.Printf("INFO: Requested %s severity", severity)
			}
		} else if service != "" {
			_, exist := report.ByServices[service]
			if !exist {
				encodeJSON(w, map[string]string{"error": fmt.Sprintf("service %s does not exist", severity)}, "")
				log.Printf("ERR: Service %s does not exist", service)
			} else {
				encodeJSON(w, report.ByServices[service], "")
				log.Printf("INFO: Requested %s severity", service)
			}
		} else {
			encodeJSON(w, report.ByID, "")
			log.Printf("INFO: Requested all incidents")
		}
	}
}

func postIncidents(store *IncidentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incident Incident
		err := json.NewDecoder(r.Body).Decode(&incident)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": "invalid JSON"}, "Error occured while encoding JSON")
			log.Printf("ERR: %s", err)
			return
		}
		store.addIncident(incident)

		response, _ := json.Marshal(incident.ID)
		log.Printf("INFO: Added new incident ID: %s", string(response))
	}
}

func getIncidentsByID(store *IncidentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		incidentID := r.PathValue("id")

		report := store.getReport()
		_, exist := report.ByID[incidentID]
		if !exist {
			errMessage := map[string]string{"error": fmt.Sprintf("non-existent incident ID (%s)", incidentID)}
			encodeJSON(w, errMessage, "")
			log.Printf("ERR: Requested non-existent incident ID: %s", incidentID)
		} else {
			encodeJSON(w, report.ByID[incidentID], "")
			log.Printf("INFO: Requested incident ID: %s", incidentID)
		}
	}
}

func deleteIncidentByID(store *IncidentStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-Type", "application/json")
		incidentID := r.PathValue("id")

		var exist bool
		for _, incident := range store.Incidents {
			if incident.ID == incidentID {
				exist = true
				store.deleteIncident(incidentID)
				encodeJSON(w, map[string]string{"success": fmt.Sprintf("deleted incident ID %s", incidentID)}, "")
				log.Printf("INFO: Successfully deleted incident with ID %s", incidentID)
			}
		}
		if !exist {
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("cannot delete incident ID %s - not found", incidentID)}, "")
			log.Printf("ERR: cannot delete incident with ID %s - not found", incidentID)
		}
	}
}

func startServer(port string, store *IncidentStore) {
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("GET /report", getIncidentReport(store))
	http.HandleFunc("GET /incidents", getGroupedIncidents(store))
	http.HandleFunc("POST /incidents", postIncidents(store))
	http.HandleFunc("GET /incidents/{id}", getIncidentsByID(store))
	http.HandleFunc("DELETE /incidents/{id}", deleteIncidentByID(store))

	log.Printf("INFO: Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func (s *IncidentStore) addIncident(incident Incident) {
	// Add write lock with mutex
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Append incident and rebuild report
	s.Incidents = append(s.Incidents, incident)
	s.recompute()

}

func (s *IncidentStore) deleteIncident(removeKey string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var newIncidents []Incident
	for _, incident := range s.Incidents {
		if incident.ID != removeKey {
			newIncidents = append(newIncidents, incident)
		}
	}
	// Set new incidents list and rebuild report
	s.Incidents = newIncidents
	s.recompute()
}

func (s *IncidentStore) getReport() *IncidentReport {
	// Retrive report from store
	// Allow parallel read-only
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	report := s.Report

	return report
}

func encodeJSON(w http.ResponseWriter, content any, errMessage string) {
	if errMessage == "" {
		// If custom error message variable is empty, use default
		errMessage = "Error encoding response"
	}
	
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Fatalf("%s: %s", errMessage, err)
	}
}