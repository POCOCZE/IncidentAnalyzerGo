package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func editHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var incident Incident
		id := r.PathValue("id")
		err := json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, err, "")
			log.Printf("ERR: %s", err)
		}
		err = store.Edit(id, incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeJSON(w, err, "")
			log.Printf("DEBUG: Incident ID: %s | Incident data: %v", id, incident)
			log.Printf("ERR: %s", err)
		} else {
			w.WriteHeader(http.StatusNoContent)
			log.Printf("Successfully edited incident %s", id)
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