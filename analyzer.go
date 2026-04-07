package main

import (
	"fmt"
	"log"
	"time"
)

func (s *IncidentStore) incidentMessage(incident Incident) (string, bool) {
    // Handles message generation for GroupIncidentsBy... functions
    // For resolved incidents print how long it took in HMS format
    // Returns message string and boolean whether the incident is resolved
    var isResolved bool
    message := fmt.Sprintf("Pending (Started: %s)", incident.StartedAt)

    if incident.ResolvedAt != "" {
        isResolved = true
        resolvedIn := s.Duration[incident.ID].HMSFormat
        message = fmt.Sprintf("Resolved in %s (Ended: %s)", resolvedIn, incident.ResolvedAt)
    }

    return message, isResolved
}

func (r *IncidentReport) groupIncidentsByService(incident Incident, store *IncidentStore) {
    // For each service name assign ServiceDetails struct value.
    // Service name keys are not sorted.

    message, isResolved := store.incidentMessage(incident)

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.ByServices[incident.Service]
    if !exists {
        r.ByServices[incident.Service] = make(map[string]IncidentReportDetails)
    }

    r.ByServices[incident.Service][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Severity: incident.Severity,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) groupIncidentsBySeverity(incident Incident, store *IncidentStore) {
    // For each severity assign SeverityDetails struct value.
    // Severity keys are not sorted by severity

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.BySeverity[incident.Severity]
    if !exists {
        r.BySeverity[incident.Severity] = make(map[string]IncidentReportDetails)
    }

    message, isResolved := store.incidentMessage(incident)
    r.BySeverity[incident.Severity][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Service: incident.Service,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) groupIncidentsByID(incident Incident, store *IncidentStore) {
    message, isResolved := store.incidentMessage(incident)
    r.ByID[incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Severity: incident.Severity,
        Service: incident.Service,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) calcMTTRSec(store *IncidentStore) {
    // Calculate Mean time to recovery - average across all
    // First is all incident durations summed, and averaged across only resolved incidents. 
    // Unresolved are not kept in mind because their duration is effectively zero and thus would avoid the calculations.

    // Sum all incident times (unresolved gets 0 seconds)
    var sum float64
    for _, incidentDuration := range store.Duration {
        sum += incidentDuration.Seconds
    }

    // Calculate average across only resolved ones
    if len(store.Duration) == 0 {
        log.Fatalf("Error: Cannot devide by zero")
    }
    resolvedIncidentCount := r.IncidentsCount - len(r.UnresolvedIDs)
    avgSeconds := int(sum / float64(resolvedIncidentCount))

    hms := time.Duration(avgSeconds) * time.Second
    r.MTTR = hms.String()
}

func (s *IncidentStore) calcIncidentDuration(incident Incident) { // , allIncidentsDuration map[string]float64 -> as output
    // Calculate incident duration for all incidents
    // Unresolved incidents will have 0 seconds duration, resolved one gets calculated

    var durationSec float64
    if incident.ResolvedAt != "" {
        startedAt, err := time.Parse(time.RFC3339, incident.StartedAt)
        if err != nil {
            log.Fatalf("Error occured while converting time %s", err)
        }
    
        resolvedAt, err := time.Parse(time.RFC3339, incident.ResolvedAt)
        if err != nil {
            log.Fatalf("Error occured while converting time %s", err)
        }

        durationSec = resolvedAt.Sub(startedAt).Seconds()
    }

    // hms := SecToHMS(int(durationSec))
    hms := time.Duration(durationSec) * time.Second
    s.Duration[incident.ID] = IncidentDuration{
        Seconds: durationSec,
        HMSFormat: hms.String(),
    }
}

func (s *IncidentStore) recompute() {
    // --- Initialize maps before they can be used --- //
    report := NewIncidentReport()

    // Calculate incidents length
    report.IncidentsCount = len(s.Incidents)

    for _, incident := range s.Incidents {
        if incident.ResolvedAt == "" {
            // --- Add unresolved to Slice --- //
            report.UnresolvedIDs = append(report.UnresolvedIDs, incident.ID)
        }

        // --- Calculate incident duration and group incidents --- //
        s.calcIncidentDuration(incident)
        report.groupIncidentsByService(incident, s)
        report.groupIncidentsBySeverity(incident, s)
        report.groupIncidentsByID(incident, s)
    }

    // --- Calculate Mean time to recovery --- //
    report.calcMTTRSec(s)

    s.Report = report
}