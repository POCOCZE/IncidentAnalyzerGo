package core

import (
	"encoding/json"
	"log"
	"os"
)

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