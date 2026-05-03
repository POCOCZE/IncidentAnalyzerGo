package main

import "time"

type Incident struct {
	ID         string `json:"id" validate:"required"`
	Title      string `json:"title" validate:"required"`
	Severity   string `json:"severity" validate:"required"`
	Service    string `json:"service_name" validate:"required"`
	StartedAt  *time.Time `json:"started_at" validate:"required"`
	ResolvedAt *time.Time `json:"resolved_at"`
	Message	   string `json:"message"`
	IsResolved bool `json:"is_resolved"`
}

type IncidentsFile struct {
	Incidents []Incident `json:"incidents"`
}

type IncidentReportDetails struct {
	Title      string `json:"title"`
	Severity   string `json:"severity,omitempty"`
	Service    string `json:"service,omitempty"`
	Message    string `json:"message"`
	IsResolved bool   `json:"is_resolved"`
}

type IncidentReport struct {
	IncidentsCount int                                         `json:"incidents_count"`
	UnresolvedIDs  []string                                    `json:"unresolved_ids"`
	MTTR           string                                      `json:"mttr"`
	ByServices     map[string]map[string]IncidentReportDetails `json:"by_services"`
	BySeverity     map[string]map[string]IncidentReportDetails `json:"by_severity"`
	ByID           map[string]IncidentReportDetails            `json:"by_id"`
}

type IncidentDuration struct {
	Seconds   float64
	HMSFormat string
}

type IncidentStorage interface {
	GetAll() ([]Incident, error)
	GetByID(id string) (Incident, error)
	Add(incident Incident) error
	DeleteByID(id string) error
	DeleteAll() error
}