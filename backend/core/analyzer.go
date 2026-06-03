package core

import (
	"fmt"
	"time"
)

func (r *IncidentReport) groupIncidentsByService(incident Incident) {
    // For each service name assign ServiceDetails struct value.
    // Service name keys are not sorted.

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.ByServices[incident.ServiceName]
    if !exists {
        r.ByServices[incident.ServiceName] = make(map[string]IncidentReportDetails)
    }

    r.ByServices[incident.ServiceName][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Severity: incident.Severity,
    }
}

func (r *IncidentReport) groupIncidentsBySeverity(incident Incident) {
    // For each severity assign SeverityDetails struct value.
    // Severity keys are not sorted by severity

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.BySeverity[incident.Severity]
    if !exists {
        r.BySeverity[incident.Severity] = make(map[string]IncidentReportDetails)
    }

    r.BySeverity[incident.Severity][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Service: incident.ServiceName,
    }
}

func (r *IncidentReport) groupIncidentsByID(incident Incident) {
    r.ByID[incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Severity: incident.Severity,
        Service: incident.ServiceName,
    }
}

func (r *IncidentReport) CalcMTTRSec(durations map[string]IncidentDuration) error {
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
        return fmt.Errorf("WARN: cannot devide by zero due to missing incidents")
    }
    resolvedIncidentCount := r.IncidentsCount - len(r.UnresolvedIDs)
    avgSeconds := int(sum / float64(resolvedIncidentCount))

    hms := time.Duration(avgSeconds) * time.Second
    r.MTTR = hms.String()

    return nil
}

// func calcIncidentDuration(incident Incident) map[string]IncidentDuration {
//     // , allIncidentsDuration map[string]float64 -> as output
//     // Calculate incident duration for all incidents
//     // Unresolved incidents will have 0 seconds duration, resolved one gets calculated

//     var durationSec float64
//     if incident.ResolvedAt != nil {
//         startedAt := incident.StartedAt
//         resolvedAt := incident.ResolvedAt
//         durationSec = resolvedAt.Sub(*startedAt).Seconds()
//     }

//     hms := time.Duration(durationSec) * time.Second
//     var durations map[string]IncidentDuration
//     durations[incident.ID] = IncidentDuration{
//         Seconds: durationSec,
//         HMSFormat: hms.String(),
//     }
//     return durations
// }

func NewIncidentReport() *IncidentReport {
    // This is constructor function: Prepares structure for report
    // Returns IncidentReport inizialized structure
    return &IncidentReport {
        ByServices: make(map[string]map[string]IncidentReportDetails),
        BySeverity: make(map[string]map[string]IncidentReportDetails),
        ByID: make(map[string]IncidentReportDetails),
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
        // durations := calcIncidentDuration(incident)
        report.groupIncidentsByService(incident)
        report.groupIncidentsBySeverity(incident)
        report.groupIncidentsByID(incident)
    }

    // --- Calculate Mean time to recovery --- //
    err := report.CalcMTTRSec(durations)
    if err != nil {
        return &IncidentReport{}, err
    }

    return report, nil
}