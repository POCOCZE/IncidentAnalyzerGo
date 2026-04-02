package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

func NewIncidentReport() *IncidentReport {
    // This is constructor function: Prepares structure for report
    // Returns IncidentReport inizialized structure
    return &IncidentReport{
        ByServices: make(map[string]map[string]IncidentReportDetails),
        BySeverity: make(map[string]map[string]IncidentReportDetails),
    }
}

func NewIncidentStore(incidents []Incident) *IncidentStore {
    // This is constructor function
    store := &IncidentStore{
        Incidents: make([]Incident, 0),
    }
    report := buildReport(incidents)

    store.Report = report
    return store
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

func main() {
    // --- Create flags --- //
    file := flag.String("file", "", "Path to incidents JSON file")
    output := flag.String("output", "stdout", "Output type. Options: stdout (default), <your-file-name>")
    serve := flag.Bool("serve", false, "Start an HTTP server")
    port := flag.String("port", "8080", "Port to be used with the HTTP server")
    flag.Parse()
    
    // --- Inizialize instance and open JSON --- //
    incidentsFile := &IncidentsFile{}
    incidentsFile.OpenInputFile(*file)

    store := NewIncidentStore(incidentsFile.Incidents)

    if *serve {
        // Start HTTP server and pass port
        startServer(*port, store)
    } else {
        // CLI mode
        // runCLI(file, output)
        printReport(*output, store.Report)
    }
}