package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func NewIncidentReport() *IncidentReport {
    // This is constructor function: Prepares structure for report
    // Returns IncidentReport inizialized structure
    return &IncidentReport{
        ByServices: make(map[string]map[string]IncidentReportDetails),
        BySeverity: make(map[string]map[string]IncidentReportDetails),
        ByID: make(map[string]IncidentReportDetails),
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
    // This is useless because pgConn string already defines whether postgres or memory should be used. So, if pgConn is defined - postgres is used, memory if empty.
    // storage := flag.String("storage", "memory", "Storage backend. Options: memory(default), postgres")
    pgConn := flag.String("db", "", "PostgreSQL conection string. Format: postgres://user:pass@address:5432/db_name")
    flag.Parse()
    
    var store IncidentStorage
    if *pgConn != "" {
        var err error
        fmt.Println("Using database...")
        store, err = NewPostgresStore(*pgConn)
        if err != nil {
            log.Fatalf("ERR: %s", err)
        }
        log.Printf("Successfully connected to database: %s", *pgConn)
    } else {
        // memory
        store = NewMemoryStore()
    }

    // --- Inizialize instance and open JSON --- //
    incidentsFile := &IncidentsFile{}
    if *file != "" {
        incidentsFile.OpenInputFile(*file)
        for _, incident := range incidentsFile.Incidents {
            // Todo: handle error. Add HTTP handler `AddList` that will loop through each incident (like here) and add list of incidents to storage. Implement this to MemoryStore for now.
            // ? This feature will allow batch processing of incidents - importing. Exporting will be added too, but after import is added.
            store.Add(incident)
        }
    }

    if *serve {
        // Start HTTP server and pass port
        startServer(*port, store)
    } else {
        // CLI mode
        incidents, _ := store.GetAll()
        report, err := BuildReport(incidents)
        if err != nil {
            fmt.Printf("ERR: %s", err)
        }
        printReport(*output, report)
    }
}