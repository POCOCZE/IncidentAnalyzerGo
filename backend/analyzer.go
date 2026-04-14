package main

import (
	"fmt"
	"time"
)

func (r *IncidentReport) incidentMessage(incident Incident, durations map[string]IncidentDuration) (string, bool) {
    // Handles message generation for GroupIncidentsBy... functions
    // For resolved incidents print how long it took in HMS format
    // Returns message string and boolean whether the incident is resolved
    var isResolved bool
    message := fmt.Sprintf("Pending (Started: %s)", incident.StartedAt)

    if incident.ResolvedAt != nil {
        isResolved = true
        resolvedIn := durations[incident.ID].HMSFormat
        message = fmt.Sprintf("Resolved in %s (Ended: %s)", resolvedIn, incident.ResolvedAt)
    }

    return message, isResolved
}

func (r *IncidentReport) groupIncidentsByService(incident Incident, durations map[string]IncidentDuration) {
    // For each service name assign ServiceDetails struct value.
    // Service name keys are not sorted.

    message, isResolved := r.incidentMessage(incident, durations)

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

func (r *IncidentReport) groupIncidentsBySeverity(incident Incident, durations map[string]IncidentDuration) {
    // For each severity assign SeverityDetails struct value.
    // Severity keys are not sorted by severity

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.BySeverity[incident.Severity]
    if !exists {
        r.BySeverity[incident.Severity] = make(map[string]IncidentReportDetails)
    }

    message, isResolved := r.incidentMessage(incident, durations)
    r.BySeverity[incident.Severity][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Service: incident.Service,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) groupIncidentsByID(incident Incident, durations map[string]IncidentDuration) {
    message, isResolved := r.incidentMessage(incident, durations)
    r.ByID[incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Severity: incident.Severity,
        Service: incident.Service,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) calcMTTRSec(durations map[string]IncidentDuration) error {
    // Calculate Mean time to recovery - average across all
    // First is all incident durations summed, and averaged across only resolved incidents. 
    // Unresolved are not kept in mind because their duration is effectively zero and thus would avoid the calculations.

    // Sum all incident times (unresolved gets 0 seconds)
    var sum float64
    for _, incidentDuration := range durations {
        sum += incidentDuration.Seconds
    }

    // Calculate average across only resolved ones
    if len(durations) == 0 {
        return fmt.Errorf("cannot devide by zero due to missing incidents")
    }
    resolvedIncidentCount := r.IncidentsCount - len(r.UnresolvedIDs)
    avgSeconds := int(sum / float64(resolvedIncidentCount))

    hms := time.Duration(avgSeconds) * time.Second
    r.MTTR = hms.String()

    return nil
}

func calcIncidentDuration(incident Incident, durations map[string]IncidentDuration) {
    // , allIncidentsDuration map[string]float64 -> as output
    // Calculate incident duration for all incidents
    // Unresolved incidents will have 0 seconds duration, resolved one gets calculated

    var durationSec float64
    if incident.ResolvedAt != nil {
        startedAt := incident.StartedAt
        resolvedAt := incident.ResolvedAt
        durationSec = resolvedAt.Sub(*startedAt).Seconds()
    }

    hms := time.Duration(durationSec) * time.Second
    durations[incident.ID] = IncidentDuration{
        Seconds: durationSec,
        HMSFormat: hms.String(),
    }
}

func BuildReport(incidents []Incident) (*IncidentReport, error) {
    // --- Initialize maps before they can be used --- //
    report := NewIncidentReport()
    report.IncidentsCount = len(incidents)

    // Create temporary duration storage
    durations := make(map[string]IncidentDuration)
    for _, incident := range incidents {
        // --- Add unresolved to Slice --- //
        if incident.ResolvedAt == nil {
            report.UnresolvedIDs = append(report.UnresolvedIDs, incident.ID)
        }

        // --- Calculate incident duration and group incidents --- //
        calcIncidentDuration(incident, durations)
        report.groupIncidentsByService(incident, durations)
        report.groupIncidentsBySeverity(incident, durations)
        report.groupIncidentsByID(incident, durations)
    }

    // --- Calculate Mean time to recovery --- //
    err := report.calcMTTRSec(durations)
    if err != nil {
        return &IncidentReport{}, err
    }

    return report, nil
}