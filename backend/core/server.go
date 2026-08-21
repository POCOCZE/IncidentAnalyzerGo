package core

import (
	"encoding/json"
	"log"
	"net/http"
)

// Health status endpoint
func HealthHandler(store IncidentStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		EncodeJSON(w, map[string]string{"status": "ok"}, "")
	}
}

// Set handlers, middlewares and start the HTTP server on particular port.
func StartServer(port string, store IncidentStorage) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", HealthHandler(store))
	mux.HandleFunc("GET /api/report", GetReportHandler(store))
	mux.HandleFunc("GET /api/incidents", GetAllHandler(store))
	mux.HandleFunc("POST /api/incidents", AddListHandler(store))
	mux.HandleFunc("POST /api/incident", AddHandler(store))
	mux.HandleFunc("PATCH /api/incident/{id}", EditHandler(store))
	mux.HandleFunc("GET /api/incident/{id}", GetByIDHandler(store))
	mux.HandleFunc("DELETE /api/incident/{id}", DeleteByIDHandler(store))

	// !This removes all incidents forever! For testing.
	mux.HandleFunc("DELETE /api/delete-all-incidents-forever", DeleteAllHandler(store))

	// empty function during development. During compilation the parameter -production is used to 
	if err := FrontendHandler(mux); err != nil {
		log.Fatalf("[StartServer] failed to serve frontend: %s", err)
	}

	// handler := CorsMiddleware(mux)
	log.Printf("INFO: Server listening on port :%s", port)
	err := http.ListenAndServe(":"+port, mux)
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
		log.Printf("%s: %s", errMessage, err)
	}
}