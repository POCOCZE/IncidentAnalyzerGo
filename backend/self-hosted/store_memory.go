package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pococze/incident-analyzer-go/backend/core"
)

// ! This file is meant only for development purposes and testing.
// ! Incidents that would be written into memory store will be lost after app restart.
// * Recommended is to use a postgresql database to make incidents persistent across restarts.

type MemoryStore struct {
    Mu sync.RWMutex
    Incidents []core.Incident
}

func NewMemoryStore() *MemoryStore {
    // This is constructor function
    return &MemoryStore{
        Incidents: make([]core.Incident, 0),
    }
}

func (m *MemoryStore) Add(incident core.Incident) error {
	// Add write lock with mutex
	m.Mu.Lock()
	defer m.Mu.Unlock()

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(incident)
	if err != nil {
		return err
	}

	// Check if the key already exist, this will prevent some bugs and errors
	for _, storedInc := range m.Incidents {
		// If incident already exist in slice - return error
		// ? Not sure whether there is more effective solution that could immidiately find the incident without relying on loops - I would need to switch to `map`, but this would result in incompatibilities...
		// ? I will keep it as it as, since Go is very fast.
		if storedInc.ID == incident.ID {
			return fmt.Errorf("error: incident already exist")
		}
	}
	
	// * This will be removed in the future
	// Check if time is defined (resolvedAt can be null) and convert to UTC if needed
	if core.IsValidTime(incident.StartedAt) {
		incident.StartedAt = core.ConvertToUTC(incident.StartedAt)
	} else {
		return fmt.Errorf("startedAt time is not valid")
	}
	if core.IsValidTime(incident.ResolvedAt) {
		incident.ResolvedAt = core.ConvertToUTC(incident.ResolvedAt)
	} else {
		incident.ResolvedAt = nil
	}

	// Append incident and rebuild report
	m.Incidents = append(m.Incidents, incident)
	_, err = core.BuildReport(m.Incidents)
	if err != nil {
		return fmt.Errorf("error building report: %s", err)
	}

	return nil
}

func (m *MemoryStore) AddList(incidents []core.Incident) error {
	for _, incident := range incidents {
		err := m.Add(incident)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryStore) Edit(id string, incident core.Incident) error {
	// * Info: id - original incident ID; incident - changed incident struct; Users are not allowed to edit incident ID - returns error (frontend blocks the input field too to edit it)
	m.Mu.Lock()
	defer m.Mu.Unlock()

	log.Printf("Received incident: %v", incident)

	var isFound bool
	for i, inc := range m.Incidents {
		if inc.ID == id {
			// return error if user tries to change incident ID or startedAt (e.g. using a bug)
			// Todo: changing startedAt is not yet implemented.
			if inc.ID != incident.ID {
				return fmt.Errorf("chanding Incident ID is not allowed")
			}
			isFound = true

			// * This will be removed in the future
			// Check if time is defined (resolvedAt can be null) and convert to UTC if needed
			if core.IsValidTime(incident.StartedAt) {
				incident.StartedAt = core.ConvertToUTC(incident.StartedAt)
			} else {
				return fmt.Errorf("startedAt time is not valid")
			}
			if core.IsValidTime(incident.ResolvedAt) {
				incident.ResolvedAt = core.ConvertToUTC(incident.ResolvedAt)
			} else {
				incident.ResolvedAt = nil
			}

			m.Incidents[i] = incident
			break
		}
	}
	if isFound {
		_, err := core.BuildReport(m.Incidents)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("incident %q not found", id)
	}

	return nil
}

func (m *MemoryStore) GetAll() ([]core.Incident, error) {
	incidents := m.Incidents
	incidentsWide, err := core.IncidentsWide(incidents)
	if err != nil {
		return nil, err
	}

	return incidentsWide, nil
}

func (m *MemoryStore) GetByID(id string) (core.Incident, error) {
	for _, incident := range m.Incidents {
		if incident.ID == id {
			inc, err := core.IncidentWide(incident)
			if err != nil {
				return core.Incident{}, err
			}
			return inc, nil
		}
	}

	return core.Incident{}, fmt.Errorf("incident ID %q not found", id)
}

func (m *MemoryStore) DeleteByID(id string) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	var newIncidents []core.Incident
	var found bool
	for _, incident := range m.Incidents {
		if incident.ID == id {
			found = true
		} else {	
			newIncidents = append(newIncidents, incident)
		}
	}
	if found {
		// Set new incidents list and rebuild report
		m.Incidents = newIncidents
		_, err := core.BuildReport(m.Incidents)
		if err != nil {
			return err
		}
	} else {
        err := fmt.Errorf("incident ID %q was not found", id)
        return err
    }

	return nil
}

func (m *MemoryStore) DeleteAll() error {
	m.Incidents = []core.Incident{}

	return nil
}