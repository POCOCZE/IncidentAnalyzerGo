package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Add Slice of incidents to selected store.
func AddListHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incidents *IncidentsFile
		err := json.NewDecoder(r.Body).Decode(&incidents)
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to decode JSON: %s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		}
		// Guard that checks if user sent correct incidents body structure
		incidentsFileLength := len(*incidents.Incidents)
		if incidentsFileLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": "no incidents provided. did you specified correct endpoint?"}, "")
			log.Printf("ERR: User tried to write 0 incidents. Maybe he specificed bad endpoint.")
			return
		}
		totalDuplicateCount, err := store.AddList(r.Context(), incidents.Incidents)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to add list of incidents: %s", err)}, "")
			log.Printf("ERR: %s", err)
			return
		}
		uniqueIncCount := incidentsFileLength - totalDuplicateCount
		w.WriteHeader(http.StatusCreated)
		if uniqueIncCount == 0 {
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("skipping %v duplicate incidents", incidentsFileLength)}, "")
			log.Printf("DEBUG: skipping %v duplicate incidents", incidentsFileLength)
			return
		}
		EncodeJSON(w, map[string]string{"message": fmt.Sprintf("added %v unique out of %v incidents", uniqueIncCount, incidentsFileLength)}, "")
		log.Printf("DEBUG: added %v/%v unique incidents", uniqueIncCount, incidentsFileLength)
	}
}

// Add one incident to Store
func AddHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incident *Incident
		err := json.NewDecoder(r.Body).Decode(&incident)

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": "invalid JSON"}, "Error occured while encoding JSON")
			log.Printf("ERR: %s", err)
			return
		}

		totalDuplicateCount, _, err := store.Add(r.Context(), incident)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to add incident %s: %s", incident.Name, err)}, "")
			log.Printf("ERR: invalid structure: %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if totalDuplicateCount != 0 {
			EncodeJSON(w, map[string]string{"info": fmt.Sprintf("skipping duplicate incident %v", incident.Name)}, "")
			log.Printf("DEBUG: skipping duplicate incident %v", incident.Name)
			return
		}
		EncodeJSON(w, map[string]string{"message": fmt.Sprintf("added new unique incident: %v", incident.Name)}, "")
		log.Printf("DEBUG: added new unique incident: %v", incident.Name)
	}
}

// Edit one particular incident, updated incident is in the request body.
func EditHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var incident *Incident
		id := r.PathValue("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("failed to parse string to UUIDv7: %s", err)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("failed to decode request body: %s", err)
			return
		}
		if err := store.Edit(r.Context(), uuid, incident); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to edit incident %s", incident.Name)}, "")
			log.Printf("failed to edit incident: %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, map[string]string{"message": fmt.Sprintf("successfully edited incident %s", incident.Name)}, "")
		log.Printf("INFO: successfully edited incident %s", id)
	}
}

// Return details of requested incident.
func GetByIDHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		id := r.PathValue("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to parse string to UUIDv7: %s", err)
			return
		}
		inc, err := store.GetByID(r.Context(), uuid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to get incident: %s", err)}, "")
			log.Printf("ERROR: failed to get incident by ID: %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, inc, "")
		log.Printf("INFO: Requested incident ID: %s", inc.ID.String())
	}
}

// Return details of all incidents.
func GetAllHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		incidents, err := store.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("failed to get incidents: %s", err)}, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, incidents, "")
	}
}

// Delete one particular incident.
func DeleteByIDHandler(store IncidentStorage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-Type", "application/json")
		id := r.PathValue("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to parse string to UUIDv7: %s", err)
			return
		}

		if err := store.DeleteByID(r.Context(), uuid); err != nil {
			w.WriteHeader(http.StatusNotFound)
			EncodeJSON(w, map[string]string{"message": fmt.Sprintf("incident ID %s not found", id)}, "")
			log.Printf("WARN: Incident with ID %s not found", id)
			return
		}
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, map[string]string{"message": fmt.Sprintf("deleted incident ID %s", id)}, "")
		log.Printf("INFO: Successfully deleted incident with ID %s", id)
	}
}

// WARNING! Deletes all incidents.
func DeleteAllHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := store.DeleteAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			EncodeJSON(w, map[string]string{"error": fmt.Sprintf("%s", err)}, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, map[string]string{"message": "deleted all incidents"}, "")
		log.Printf("DEBUG: deleted all incidents")
	}
}