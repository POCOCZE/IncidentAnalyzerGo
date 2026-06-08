package core

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// Health status endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	EncodeJSON(w, map[string]string{"status": "ok"}, "")
}

// Set handlers, middlewares and start the HTTP server on particular port.
func StartServer(ctx context.Context, port string, store IncidentStorage) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("GET /report", GetReportHandler(store))
	mux.HandleFunc("GET /incidents", GetAllHandler(store))
	mux.HandleFunc("POST /incidents", AddListHandler(store))
	mux.HandleFunc("POST /incident", AddHandler(store))
	mux.HandleFunc("PATCH /incident/{id}", EditHandler(store))
	mux.HandleFunc("GET /incidents/{id}", GetByIDHandler(store))
	mux.HandleFunc("DELETE /incidents/{id}", DeleteByIDHandler(store))

	// !This removes all incidents forever! For testing.
	mux.HandleFunc("DELETE /delete-all-incidents-forever", DeleteAllHandler(store))

	handler := CorsMiddleware(mux)
	log.Printf("INFO: Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

// Used to structure JSON encoded responses.
func EncodeJSON(w http.ResponseWriter, content any, errMessage string) {
	if errMessage == "" {
		// If custom error message variable is empty, use default
		errMessage = "Error encoding response"
	}

	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Fatalf("%s: %s", errMessage, err)
	}
}