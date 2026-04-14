package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func serveJsonOutput(report *IncidentReport) []byte {
    // Prints report as JSON
    // MarshalIndent automatically sorts alphabetically
    // Output is []byte type that needs to be further processed

    jsonData, err := json.MarshalIndent(report, "", "  ")
    if err != nil {
        log.Fatalf("Error: Can't marshal indent data: %s", err)
    }

    return jsonData
}

func printReport(output string, report *IncidentReport) {
    if output == "stdout" {
        jsonData := serveJsonOutput(report)
        fmt.Println(string(jsonData))
    } else {
        file, err := os.Create(output)
        if err != nil {
            log.Fatalf("Error creating output file: %s", err)
        }
        defer file.Close()

        jsonData := serveJsonOutput(report)
        _, err = file.Write(jsonData)
        if err != nil {
            log.Fatalf("Error writing report to output file: %s", err)
        }
    }
}
