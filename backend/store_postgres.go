package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	// Constructor function
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("Oppening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Connecting to database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) Add(incident Incident) error {
	_, err := p.db.Exec(
		`INSERT INTO incidents (id, title, severity, service_name, started_at, resolved_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		incident.ID,
		incident.Title,
		incident.Severity,
		incident.Service,
		incident.StartedAt,
		incident.ResolvedAt,
	)

	return err
}

func(p *PostgresStore) AddList(incidents []Incident) error {
	for _, incident := range incidents {
		err := p.Add(incident)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresStore) GetAll() ([]Incident, error) {
	rows, err := p.db.Query(`SELECT id, title, severity, service_name, started_at, resolved_at FROM incidents`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var inc Incident
		err := rows.Scan(&inc.ID, &inc.Title, &inc.Severity, &inc.Service, &inc.StartedAt, &inc.ResolvedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}

	return incidents, rows.Err()
}

func (p *PostgresStore) GetByID(id string) (Incident, error) {
	var inc Incident
	row := p.db.QueryRow(`SELECT id, title, severity, service_name, started_at, resolved_at FROM incidents WHERE id = $1`, id)
	err := row.Scan(&inc.ID, &inc.Title, &inc.Severity, &inc.Service, &inc.StartedAt, &inc.ResolvedAt)
	if err == sql.ErrNoRows {
		return inc, fmt.Errorf("Incident ID %q not found", id)
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