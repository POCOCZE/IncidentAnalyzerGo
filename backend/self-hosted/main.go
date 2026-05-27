package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pococze/incidentanalyzergo/backend/core"
)

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

	log.Printf("+-----------------------+")
	log.Printf("| Built by PradkaDotDev |")
	log.Printf("+-----------------------+")
	log.Printf("|  More on: pradka.dev  |")
	log.Printf("+-----------------------+")

    var store core.IncidentStorage
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
    incidentsFile := &core.IncidentsFile{}
    if *file != "" {
        incidentsFile.OpenInputFile(*file)
        for _, incident := range incidentsFile.Incidents {
            // Todo: handle error. Add HTTP handler `AddList` that will loop through each incident (like here) and add list of incidents to storage. Implement this to MemoryStore for now.
            // ? This feature will allow batch processing of incidents - importing. Exporting will be added too, but after import is added.
            err := store.Add(incident)
            if err != nil {
                fmt.Printf("%s", err)
            }
        }
    }

    if *serve {
        // Start HTTP server and pass port
        core.StartServer(*port, store)
    } else {
        // CLI mode
        incidents, _ := store.GetAll()
        report, err := core.BuildReport(incidents)
        if err != nil {
            fmt.Printf("%s", err)
        }
        err = core.PrintReport(*output, report)
        if err != nil {
            fmt.Printf("%s", err)
        }
    }
}