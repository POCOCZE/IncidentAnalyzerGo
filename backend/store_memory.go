package main

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

type MemoryStore struct {
    Mu sync.RWMutex
    Incidents []Incident
}

func NewMemoryStore() *MemoryStore {
    // This is constructor function
    return &MemoryStore{
        Incidents: make([]Incident, 0),
    }
}

func (m *MemoryStore) Add(incident Incident) error {
	// Add write lock with mutex
	m.Mu.Lock()
	defer m.Mu.Unlock()

	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(incident)
	if err != nil {
		return err
	}

	// Append incident and rebuild report
	m.Incidents = append(m.Incidents, incident)
	BuildReport(m.Incidents)

	return nil
}

func (m *MemoryStore) GetAll() ([]Incident, error) {
	incidents := m.Incidents

	return incidents, nil
}

func (m *MemoryStore) GetByID(id string) (Incident, error) {
	for _, incident := range m.Incidents {
		if incident.ID == id {
			return incident, nil
		}
	}

	return Incident{}, fmt.Errorf("Incident ID %q not found", id)
}

func (m *MemoryStore) DeleteByID(id string) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	var newIncidents []Incident
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
		BuildReport(m.Incidents)
	} else {
        err := fmt.Errorf("Incident ID %q was not found", id)
        return err
    }

	return nil
}

func (m *MemoryStore) DeleteAll() error {
	m.Incidents = []Incident{}

	return nil
}