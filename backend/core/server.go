package core

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Health status endpoint
	// Set content type to application/json	
	// Write OK as JSON reposnse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeJSON(w, map[string]string{"status": "ok"}, "")
}

func StartServer(ctx context.Context, port string, store IncidentStorage) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("GET /report", getReportHandler(ctx, store))
	mux.HandleFunc("GET /incidents", getAllHandler(ctx, store))
	mux.HandleFunc("POST /incidents", addListHandler(ctx, store))
	mux.HandleFunc("POST /incident", addHandler(ctx, store))
	mux.HandleFunc("PATCH /incident/{id}", editHandler(ctx, store))
	mux.HandleFunc("GET /incidents/{id}", getByIDHandler(ctx, store))
	mux.HandleFunc("DELETE /incidents/{id}", deleteByIDHandler(ctx, store))

	// !This removes all incidents forever! For testing.
	mux.HandleFunc("DELETE /delete-all-incidents-forever", deleteAllHandler(ctx, store))

	handler := corsMiddleware(mux)
	log.Printf("INFO: Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func encodeJSON(w http.ResponseWriter, content any, errMessage string) {
	if errMessage == "" {
		// If custom error message variable is empty, use default
		errMessage = "Error encoding response"
	}
	
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Fatalf("%s: %s", errMessage, err)
	}
}