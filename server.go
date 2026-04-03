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

func getIncidentsBySeverity(store *IncidentStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		severity := r.URL.Query().Get("severity")
		report := store.getReport()

		if severity == "" {
			encodeJSON(w, report.ByID, "")
			log.Printf("Requested all incidents")
		} else {
			encodeJSON(w, report.BySeverity[severity], "")
			log.Printf("Requested %s severity", severity)
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
			errMessage := map[string]string{"error": fmt.Sprintf("incident ID (%s) does not exist", incidentID)}
			encodeJSON(w, errMessage, "")
			log.Printf("ERR: Requested incident ID: %s, does not exist!", incidentID)
		} else {
			encodeJSON(w, report.ByID[incidentID], "")
			log.Printf("INFO: Requested incident ID: %s", incidentID)
		}
	}
}

func startServer(port string, store *IncidentStore) {
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("GET /incidents", getIncidentsBySeverity(store))
	http.HandleFunc("POST /incidents", postIncidents(store))
	http.HandleFunc("GET /incidents/{id}", getIncidentsByID(store))

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