package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Incident struct {
    ID string `json:"id"`
    Title string `json:"title"`
    Severity string `json:"severity"`
    Service string `json:"service_name"`
    StartedAt string `json:"started_at"`
    ResolvedAt string `json:"resolved_at"`
}

type IncidentsFile struct {
    Incidents []Incident `json:"incidents"`
}

type IncidentDuration struct {
    Seconds float64
    HMSFormat string
}

type IncidentPrivate struct {
    Duration map[string]IncidentDuration
}

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
}

func NewIncidentReport() *IncidentReport {
    // This is constructor function: Prepares structure for report
    // Returns IncidentReport inizialized structure
    return &IncidentReport{
        ByServices: make(map[string]map[string]IncidentReportDetails),
        BySeverity: make(map[string]map[string]IncidentReportDetails),
    }
}

func NewIncidentPrivate() *IncidentPrivate {
    // Prepares structure for private variable by initialiting its structure to be used
    return &IncidentPrivate{
        Duration: make(map[string]IncidentDuration),
    }
}
func (f *IncidentsFile) OpenInputFile(file string) {
    // Open JSON file from input - file path

    // Read file
    data, err := os.ReadFile(file)
    if err != nil {
        log.Fatalf("Error reading file: %s", err)
    }

    // Unmarshal encoded JSON data
    err = json.Unmarshal(data, &f)
    if err != nil {
        log.Fatalf("Error unmarshal JSON file: %s", err)
    }
}

func (p *IncidentPrivate) CalcIncidentDuration(incident Incident) { // , allIncidentsDuration map[string]float64 -> as output
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
    p.Duration[incident.ID] = IncidentDuration{
        Seconds: durationSec,
        HMSFormat: hms.String(),
    }
}

func (p *IncidentPrivate) IncidentMessage(incident Incident) (string, bool) {
    // Handles message generation for GroupIncidentsBy... functions
    // For resolved incidents print how long it took in HMS format
    // Returns message string and boolean whether the incident is resolved
    var isResolved bool
    message := fmt.Sprintf("Pending (Started: %s)", incident.StartedAt)

    if incident.ResolvedAt != "" {
        isResolved = true
        resolvedIn := p.Duration[incident.ID].HMSFormat
        message = fmt.Sprintf("Resolved in %s (Ended: %s)", resolvedIn, incident.ResolvedAt)
    }

    return message, isResolved
}

func (r *IncidentReport) GroupIncidentsByService(incident Incident, private *IncidentPrivate) {
    // For each service name assign ServiceDetails struct value.
    // Service name keys are not sorted.

    message, isResolved := private.IncidentMessage(incident)

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

func (r *IncidentReport) GroupIncidentsBySeverity(incident Incident, private *IncidentPrivate) {
    // For each severity assign SeverityDetails struct value.
    // Severity keys are not sorted by severity

    message, isResolved := private.IncidentMessage(incident)

    // Check whether serviceName keys exist, create them otherwise
    _, exists := r.BySeverity[incident.Severity]
    if !exists {
        r.BySeverity[incident.Severity] = make(map[string]IncidentReportDetails)
    }

    r.BySeverity[incident.Severity][incident.ID] = IncidentReportDetails{
        Title: incident.Title,
        Service: incident.Service,
        Message: message,
        IsResolved: isResolved,
    }
}

func (r *IncidentReport) CalcMTTRSec(private *IncidentPrivate) {
    // Calculate Mean time to recovery - average across all
    // First is all incident durations summed, and averaged across only resolved incidents. 
    // Unresolved are not kept in mind because their duration is effectively zero and thus would avoid the calculations.

    // Sum all incident times (unresolved gets 0 seconds)
    var sum float64
    for _, incidentDuration := range private.Duration {
        sum += incidentDuration.Seconds
    }

    // Calculate average across only resolved ones
    if len(private.Duration) == 0 {
        log.Fatalf("Error: Cannot devide by zero")
    }
    resolvedIncidentCount := r.IncidentsCount - len(r.UnresolvedIDs)
    avgSeconds := int(sum / float64(resolvedIncidentCount))

    hms := time.Duration(avgSeconds) * time.Second
    r.MTTR = hms.String()
}

func (r *IncidentReport) ServeJsonOutput() []byte {
    // Prints report as JSON
    // MarshalIndent automatically sorts alphabetically
    // Output is []byte type that needs to be further processed

    jsonData, err := json.MarshalIndent(r, "", "  ")
    if err != nil {
        log.Fatalf("Error: Can't marshal indent data: %s", err)
    }

    return jsonData
}

func (r IncidentReport) PrintReport(output *string) {
    if *output == "stdout" {
        jsonData := r.ServeJsonOutput()
        fmt.Println(string(jsonData))
    } else {
        file, err := os.Create(*output)
        if err != nil {
            log.Fatalf("Error creating output file: %s", err)
        }
        defer file.Close()

        jsonData := r.ServeJsonOutput()
        _, err = file.Write(jsonData)
        if err != nil {
            log.Fatalf("Error writing report to output file: %s", err)
        }
    }
}

func main() {
    // --- Create instance of IncidentsFile --- //
    incidentsFile := &IncidentsFile{}

    // --- Create flags --- //
    file := flag.String("file", "", "Path to incidents JSON file")
    output := flag.String("output", "stdout", "Output type. Options: stdout (default), <your-file-name>")
    flag.Parse()

    // --- Initialize maps before they can be used --- //
    report := NewIncidentReport()
    private := NewIncidentPrivate()
    
    // --- Open JSON file at given path and calculate incidents count --- //
    incidentsFile.OpenInputFile(*file)
    report.IncidentsCount = len(incidentsFile.Incidents)

    for _, incident := range incidentsFile.Incidents {
        if incident.ResolvedAt == "" {
            // --- Add unresolved to Slice --- //
            report.UnresolvedIDs = append(report.UnresolvedIDs, incident.ID)
        }

        // --- Calculate incident duration and group incidents --- //
        private.CalcIncidentDuration(incident)
        report.GroupIncidentsByService(incident, private)
        report.GroupIncidentsBySeverity(incident, private)
    }

    // --- Calculate Mean time to recovery --- //
    report.CalcMTTRSec(private)

    // --- Print Results --- //
    report.PrintReport(output)
}