package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Health struct {
	Status map[string]string
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Health status endpoint
	// Set content type to application/json	
	// Write {"status": "ok"} as JSON reposnse
	status := Health{
		Status: make(map[string]string),
	}
	status.Status["status"] = "ok"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(status.Status)
	if err != nil {
		log.Fatalf("Error encoding response: %s", err)
	}
}

func getIncidentsBySeverity(report *IncidentReport) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		severity := r.URL.Query().Get("severity")
		w.Header().Set("Content-Type", "application/json")
		if severity == "" {
			err := json.NewEncoder(w).Encode(report)
			if err != nil {
				log.Fatalf("Error encoding response: %s", err)
			}
		} else {
			err := json.NewEncoder(w).Encode(report.BySeverity[severity])
			if err != nil {
				log.Fatalf("Error encoding response: %s", err)
			}
		}
	}
}

func postIncidents(store *IncidentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incident Incident
		err := json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERR: %s", err)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
			if err != nil {
				log.Fatalf("Error encoding JSON error: %s", err)
			}
			return
		}
		store.Incidents = append(store.Incidents, incident)
		buildReport(store.Incidents)
		// json.NewEncoder(w).Encode(incident)
		response, _ := json.Marshal(incident.ID)
		log.Printf("INFO: Processed incident ID: %s", string(response))
	}
}

func startServer(port string, store *IncidentStore) {
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("GET /incidents", getIncidentsBySeverity(store.Report))
	http.HandleFunc("POST /incidents", postIncidents(store))

	log.Printf("Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}