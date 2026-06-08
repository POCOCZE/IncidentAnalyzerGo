package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/pococze/incidentanalyzergo/backend/core"
)

func main() {
    // Create flags
    // Todo: remove all CLI features
    // file := flag.String("file", "", "Path to incidents JSON file")
    // output := flag.String("output", "stdout", "Output type. Options: stdout (default), <your-file-name>")
    // serve := flag.Bool("serve", false, "Start an HTTP server")
    port := flag.String("port", "8080", "Port to be used with the HTTP server")
    pgConn := flag.String("db", "", "PostgreSQL conection string. Format: postgres://user:pass@address:5432/db_name.")
    runDev := flag.Bool("dev", false, "Run in memory store. Development and testing only!")
    flag.Parse()
    
    // Todo: Add env vars along with flags.
    // Note: first variant: export INCPGCONN="postgres://your_user:your_pass@localhost:5432/your_db_name"
    // Note: second variant: export INC_PG_USER="user_user", export INC_PG_PASS="your_pass", export INC_PG_NAME="your_database_name", export INC_PG_PORT="5432"
    // Note: export INC_PORT="8080"
    // os.LookupEnv()

	log.Println("+-----------------------+")
	log.Println("| Built by PradkaDotDev |")
	log.Println("+-----------------------+")
	log.Println("|  More on: pradka.dev  |")
	log.Println("+-----------------------+")
    log.Println("")

    var store core.IncidentStorage

    // Create context - mainly for database timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

    if *pgConn != "" {
        var err error
        log.Println("Using database...")
        store, err = NewPostgresStore(ctx, *pgConn)
        if err != nil {
            log.Fatalf("ERR: %s", err)
        }
        log.Printf("Successfully connected to database: %s", *pgConn)
        cancel()
    // Run development mode if postgres connection string is not specified and '-dev' parameter is specified.
    } else if *runDev && *pgConn == "" {
        // * In memory store - after restart, everything is gone. Used only for testing and development purposes!
        store = NewMemoryStore()
        log.Println("+---------------------------------------+")
        log.Println("|   !!! RUNNING DEVELOPMENT MODE !!!    |")
        log.Println("+---------------------------------------+")
        log.Println("| !!! RESTART REMOVES ALL INCIDENTS !!! |")
        log.Println("+---------------------------------------+")
        log.Println("")
    } else {
        log.Println("WARN: Specified wrong combination of parameters. See help for more info ('-h')")
        log.Println("Quitting immidiately...")
        os.Exit(1)
    }

    // Inizialize instance and open JSON
    // if *file != "" {
    //     incidentsFile := &core.IncidentsFile{}
    //     incidentsFile.OpenInputFile(*file)
    //     for _, incident := range incidentsFile.Incidents {
    //         _, err := store.Add(ctx, incident)
    //         if err != nil {
    //             fmt.Printf("%s", err)
    //             continue
    //         }
    //     }
    // }

    // Start HTTP server and pass port
    core.StartServer(ctx, *port, store)

    // if *serve {
    //     core.StartServer(ctx, *port, store)
    // } else {
    //     incidents, _ := store.GetAll(ctx)
    //     report, err := core.BuildReport(incidents)
    //     if err != nil {
    //         fmt.Printf("%s", err)
    //     }
    //     err = core.PrintReport(*output, report)
    //     if err != nil {
    //         fmt.Printf("%s", err)
    //     }
    // }
}