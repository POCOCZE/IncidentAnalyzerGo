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
	// Write OK as JSON reposnse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeJSON(w, map[string]string{"status": "ok"}, "")
}

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

func getAllHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		incidents, err := store.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]error{"error": err}, "")
		} else {
			w.WriteHeader(http.StatusOK)
			encodeJSON(w, incidents, "")
		}
	}
}

func addListHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incidents IncidentsFile
		err := json.NewDecoder(r.Body).Decode(&incidents)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		}
		// Guard that checks if user sent correct incidents body structure
		if len(incidents.Incidents) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": "no incidents provided. did you specified correct endpoint?"}, "")
			log.Printf("ERR: User tried to write 0 incidents. Maybe he specificed bad endpoint.")
		}
		err = store.AddList(incidents.Incidents)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			log.Printf("INFO: Added list of %v incidents", len(incidents.Incidents))
		}
	}
}

func addHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incident Incident
		err := json.NewDecoder(r.Body).Decode(&incident)

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": "invalid JSON"}, "Error occured while encoding JSON")
			log.Printf("ERR: %s", err)
			return
		}

		err = store.Add(incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("invalid structure: %s", err)}, "")
			log.Printf("ERR: Invalid structure: %s", err)
		} else {
			w.WriteHeader(http.StatusCreated)
			// This is not needed at all
			// encodeJSON(w, "", "")
			log.Printf("INFO: Added new incident ID: %s", incident.ID)
		}
	}
}

func getByIDHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		incidentID := r.PathValue("id")

		inc, err := store.GetByID(incidentID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := map[string]error{"error": err}
			encodeJSON(w, errMessage, "")
			log.Printf("ERR: %s", err)
		} else {
			w.WriteHeader(http.StatusOK)
			encodeJSON(w, inc, "")
			log.Printf("INFO: Requested incident ID: %s", incidentID)
		}
	}
}

func deleteByIDHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-Type", "application/json")
		id := r.PathValue("id")

		// If id gathered from path actually exist
		// If not, return error 404
		err := store.DeleteByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			encodeJSON(w, map[string]string{"warn": fmt.Sprintf("incident ID %s not found", id)}, "")
			log.Printf("WARN: Incident with ID %s not found", id)
		} else {
			// w.WriteHeader(http.StatusNoContent)
			w.WriteHeader(http.StatusOK)
			encodeJSON(w, map[string]string{"success": fmt.Sprintf("deleted incident ID %s", id)}, "")
			log.Printf("INFO: Successfully deleted incident with ID %s", id)
		}
	}
}

func deleteAllHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := store.DeleteAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
		} else {
			w.WriteHeader(http.StatusOK)
			encodeJSON(w, map[string]string{"success": "deleted all incidents"}, "")
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func startServer(port string, store IncidentStorage) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("GET /report", getReportHandler(store))
	mux.HandleFunc("GET /incidents", getAllHandler(store))
	mux.HandleFunc("POST /incidents", addListHandler(store))
	mux.HandleFunc("POST /incident", addHandler(store))
	mux.HandleFunc("GET /incidents/{id}", getByIDHandler(store))
	mux.HandleFunc("DELETE /incidents/{id}", deleteByIDHandler(store))

	// !This removes all incidents forever! For testing.
	mux.HandleFunc("DELETE /delete-all-incidents-forever", deleteAllHandler(store))

	handler := corsMiddleware(mux)
	log.Printf("INFO: Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
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