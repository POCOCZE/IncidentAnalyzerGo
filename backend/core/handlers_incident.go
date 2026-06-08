package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Add Slice of incidents to selected store.
func AddListHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incidents IncidentsFile
		err := json.NewDecoder(r.Body).Decode(&incidents)
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		}
		// Guard that checks if user sent correct incidents body structure
		if len(incidents.Incidents) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": "no incidents provided. did you specified correct endpoint?"}, "")
			log.Printf("ERR: User tried to write 0 incidents. Maybe he specificed bad endpoint.")
		}
		totalDuplicateCount, err := store.AddList(r.Context(), incidents.Incidents)
		uniqueIncCount := len(incidents.Incidents) - totalDuplicateCount
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			if uniqueIncCount == 0 {
				EncodeJSON(w, map[string]string{"info": fmt.Sprintf("skipping %v duplicate incidents", len(incidents.Incidents))}, "")
				log.Printf("INFO: skipping %v duplicate incidents", len(incidents.Incidents))
			} else {
				EncodeJSON(w, map[string]string{"info": fmt.Sprintf("added %v unique out of %v incidents", uniqueIncCount, len(incidents.Incidents))}, "")
				log.Printf("INFO: added %v/%v unique incidents", uniqueIncCount, len(incidents.Incidents))
			}
		}
	}
}

// Add one incident to Store
func AddHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incident Incident
		err := json.NewDecoder(r.Body).Decode(&incident)

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": "invalid JSON"}, "Error occured while encoding JSON")
			log.Printf("ERR: %s", err)
			return
		}

		totalDuplicateCount, err := store.Add(r.Context(), incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("invalid structure: %s", err)}, "")
			log.Printf("ERR: invalid structure: %s", err)
		} else {
			// w.WriteHeader(http.StatusCreated)
			w.WriteHeader(http.StatusOK)
			if totalDuplicateCount != 0 {
				EncodeJSON(w, map[string]string{"warn": fmt.Sprintf("skipping duplicate incident %v", incident.ID)}, "")
				log.Printf("WARN: skipping duplicate incident %v", incident.ID)
			} else {
				EncodeJSON(w, map[string]string{"info": fmt.Sprintf("added new unique incident ID: %v", incident.ID)}, "")
				log.Printf("INFO: added new unique incident ID: %v", incident.ID)
			}
		}
	}
}

// Edit one particular incident, updated incident is in the request body.
func EditHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var incident Incident
		id := r.PathValue("id")
		err := json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, err, "")
			log.Printf("ERR: %s", err)
		}
		err = store.Edit(r.Context(), id, incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, err, "")
			log.Printf("DEBUG: Incident ID: %s | Incident data: %v", id, incident)
			log.Printf("ERR: %s", err)
		} else {
			w.WriteHeader(http.StatusNoContent)
			log.Printf("Successfully edited incident %s", id)
		}
	}
}

// Return details of requested incident.
func GetByIDHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		incidentID := r.PathValue("id")

		inc, err := store.GetByID(r.Context(), incidentID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage := map[string]error{"error": err}
			EncodeJSON(w, errMessage, "")
			log.Printf("ERR: %s", err)
		} else {
			w.WriteHeader(http.StatusOK)
			EncodeJSON(w, inc, "")
			log.Printf("INFO: Requested incident ID: %s", incidentID)
		}
	}
}

// Return details of all incidents.
func GetAllHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		incidents, err := store.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]error{"error": err}, "")
		} else {
			w.WriteHeader(http.StatusOK)
			EncodeJSON(w, incidents, "")
		}
	}
}

// Delete one particular incident.
func DeleteByIDHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-Type", "application/json")
		id := r.PathValue("id")

		// If id gathered from path actually exist
		// If not, return error 404
		err := store.DeleteByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			EncodeJSON(w, map[string]string{"warn": fmt.Sprintf("incident ID %s not found", id)}, "")
			log.Printf("WARN: Incident with ID %s not found", id)
		} else {
			// w.WriteHeader(http.StatusNoContent)
			w.WriteHeader(http.StatusOK)
			EncodeJSON(w, map[string]string{"success": fmt.Sprintf("deleted incident ID %s", id)}, "")
			log.Printf("INFO: Successfully deleted incident with ID %s", id)
		}
	}
}

// WARNING! Delete all incidents
func DeleteAllHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := store.DeleteAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
		} else {
			w.WriteHeader(http.StatusOK)
			EncodeJSON(w, map[string]string{"success": "deleted all incidents"}, "")
		}
	}
}