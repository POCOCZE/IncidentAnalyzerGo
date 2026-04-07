package main

import "sync"

type Incident struct {
    ID string `json:"id" validate:"required"`
    Title string `json:"title" validate:"required"`
    Severity string `json:"severity" validate:"required"`
    Service string `json:"service_name" validate:"required"`
    StartedAt string `json:"started_at" validate:"required"`
    ResolvedAt string `json:"resolved_at"`
}

type IncidentsFile struct {
    Incidents []Incident `json:"incidents"`
}

// type IncidentPrivate struct {
//     Duration map[string]IncidentDuration
// }

type IncidentReportDetails struct {
    Title string `json:"title"`
    Severity string `json:"severity,omitempty"`
    Service string `json:"service,omitempty"`
    Message string `json:"message"`
    IsResolved bool `json:"is_resolved"`
}

type IncidentReport struct {
    IncidentsCount int `json:"incidents_count"`
    UnresolvedIDs []string `json:"unresolved_ids"`
    MTTR string `json:"mttr"`
    ByServices map[string]map[string]IncidentReportDetails `json:"by_services"`
    BySeverity map[string]map[string]IncidentReportDetails `json:"by_severity"`
    ByID map[string]IncidentReportDetails `json:"by_id"`
}

type IncidentDuration struct {
    Seconds float64
    HMSFormat string
}

type IncidentStore struct {
    Mu sync.RWMutex
    Incidents []Incident
    Report *IncidentReport
    Duration map[string]IncidentDuration
}