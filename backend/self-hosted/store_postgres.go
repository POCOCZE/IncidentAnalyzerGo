package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pococze/incidentanalyzergo/backend/core"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	// Constructor function
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("oppening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) Add(incident core.Incident) error {
	var err error
	_, err = p.db.Exec(
		`INSERT INTO incidents (id, title, severity, service_name, started_at, resolved_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		incident.ID,
		incident.Title,
		incident.Severity,
		incident.Service,
		incident.StartedAt,
		incident.ResolvedAt,
	)
	if err != nil {
		return err
	}
	incidents, err := p.GetAll()
	if err != nil {
		return err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return err
	}
	return nil
}

func(p *PostgresStore) AddList(incidents []core.Incident) error {
	for _, incident := range incidents {
		err := p.Add(incident)
		if err != nil {
			return err
		}
	}
	return nil
}

func(p *PostgresStore) Edit(id string, incident core.Incident) error {
	var inc core.Incident
	inc, err := p.GetByID(id)
	if err != nil {
		return err
	}
	
	// Todo: restrict changing StartedAt too. Not yet implemented.
	if inc.ID != incident.ID {
		return fmt.Errorf("changing Incident ID is not allowed")
	}
	_, err = p.db.Exec(
		`UPDATE incidents SET id = $1, title = $2, severity = $3, service_name = $4, started_at = $5, resolved_at = $6 WHERE id = $7`,
		incident.ID,
		incident.Title,
		incident.Severity,
		incident.Service,
		incident.StartedAt,
		incident.ResolvedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("incident ID %q not found", id)
	}

	// Get All incidents to satisfy function - rebuild report
	incidents, err := p.GetAll()
	if err != nil {
		return err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) GetAll() (incidentsWide []core.Incident, err error) {
	rows, err := p.db.Query(`SELECT id, title, severity, service_name, started_at, resolved_at FROM incidents`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var incidents []core.Incident
	for rows.Next() {
		var inc core.Incident
		err := rows.Scan(&inc.ID, &inc.Title, &inc.Severity, &inc.Service, &inc.StartedAt, &inc.ResolvedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	incidentsWide, err = core.IncidentsWide(incidents)
	if err != nil {
		return nil, err
	}

	return incidentsWide, rows.Err()
}

func (p *PostgresStore) GetByID(id string) (core.Incident, error) {
	var inc core.Incident
	row := p.db.QueryRow(`SELECT id, title, severity, service_name, started_at, resolved_at FROM incidents WHERE id = $1`, id)
	err := row.Scan(&inc.ID, &inc.Title, &inc.Severity, &inc.Service, &inc.StartedAt, &inc.ResolvedAt)
	if err == sql.ErrNoRows {
		return inc, fmt.Errorf("incident ID %q not found", id)
	}
	return inc, err
}

func (p *PostgresStore) DeleteByID(id string) error {
	result, err := p.db.Exec(`DELETE FROM incidents WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("incident ID %q not found", id)
	}

	// Get All incidents to satisfy function - rebuild report
	incidents, err := p.GetAll()
	if err != nil {
		return err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStore) DeleteAll() error {
	// !This removes everything from the database forever!
	_, err := p.db.Exec(`DELETE FROM incidents`)
	if err != nil {
		return err
	}

	return nil
}