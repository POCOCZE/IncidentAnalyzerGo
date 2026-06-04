package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pococze/incidentanalyzergo/backend/core"
	"github.com/pococze/incidentanalyzergo/backend/self-hosted/database"
)

type PostgresStore struct {
	Pool *pgxpool.Pool
	Queries *database.Queries
}

func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("oppening database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	queries := database.New(pool)
	return &PostgresStore{
		Pool: pool,
		Queries: queries,
	}, nil
}

func timeToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil || t.IsZero() {
		return pgtype.Timestamptz{
			Time: time.Time{},
			InfinityModifier: 0,
			Valid: false, // This represents SQL NULL
		}
	}
	return pgtype.Timestamptz{
		Time: *t,
		InfinityModifier: 0,
		Valid: true,
	}
}

// func boolToPgBool(b bool) pgtype.Bool {
// 	return pgtype.Bool{
// 		Bool: b,
// 		Valid: true,
// 	}
// }

func pgTimestamptzToTime(pgTime pgtype.Timestamptz) *time.Time {
	if !pgTime.Valid || pgTime.Time.IsZero() {
		return nil
	}
	return &pgTime.Time
}

func (p *PostgresStore) incidentToAddParams(incident core.Incident) database.AddParams {
	params := database.AddParams{
		ID: incident.ID,
		Title: incident.Title,
		Severity: incident.Severity,
		ServiceName: incident.ServiceName,
		StartedAt: timeToPgTimestamptz(incident.StartedAt),
		ResolvedAt: timeToPgTimestamptz(incident.ResolvedAt),
	}
	return params
}

func (p *PostgresStore) incidentToEditParams(id string, incident core.Incident) database.EditParams {
	params := database.EditParams{
		ID: incident.ID,
		Title: incident.Title,
		Severity: incident.Severity,
		ServiceName: incident.ServiceName,
		StartedAt: timeToPgTimestamptz(incident.StartedAt),
		ResolvedAt: timeToPgTimestamptz(incident.ResolvedAt),
		ID_2: id, 
	}
	return params
}

func (p *PostgresStore) pgIncidentToIncident(pgIncident database.Incident) core.Incident {
	params := core.Incident{
		ID: pgIncident.ID,
		Title: pgIncident.Title,
		Severity: pgIncident.Severity,
		ServiceName: pgIncident.ServiceName,
		StartedAt: pgTimestamptzToTime(pgIncident.StartedAt),
		ResolvedAt: pgTimestamptzToTime(pgIncident.ResolvedAt),
	}
	return params
}

func (p *PostgresStore) Add(ctx context.Context, incident core.Incident) (int, error) {
	// check if incident already exist before adding.
	inc, err := p.GetByID(ctx, incident.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		// Incident does not exist - add it
		params := p.incidentToAddParams(incident)
		err = p.Queries.Add(ctx, params)
		if err != nil {
			return 0, err
		}
	}
	var duplicateCount int
	if inc.ID == incident.ID {
		duplicateCount = 1
		// log.Printf("WARN: incident %q already exist, skipping\n", incident.ID)
	}
	if err != nil {
		return 0, err
	}

	// Build report
	incidents, err := p.GetAll(ctx)
	if err != nil {
		return 0, err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return 0, err
	}
	return duplicateCount, nil
}

func(p *PostgresStore) AddList(ctx context.Context, incidents []core.Incident) (int, error) {
	var totalDuplicateCount int
	for _, incident := range incidents {
		duplicateCount, err := p.Add(ctx, incident)
		totalDuplicateCount += duplicateCount
		if err != nil {
			return 0, err
		}
	}
	return totalDuplicateCount, nil
}

func(p *PostgresStore) Edit(ctx context.Context, id string, incident core.Incident) error {
	var inc core.Incident
	inc, err := p.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Todo: restrict changing StartedAt too. Not yet implemented.
	if inc.ID != incident.ID {
		return fmt.Errorf("changing Incident ID is not allowed")
	}
	params := p.incidentToEditParams(id, incident)
	err = p.Queries.Edit(ctx, params)
	if err != nil {
		return err
	}

	// GetAll incidents to satisfy build report function
	incidents, err := p.GetAll(ctx)
	if err != nil {
		return err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) GetAll(ctx context.Context) ([]core.Incident, error) {
	pgIncidents, err := p.Queries.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var incidents []core.Incident
	for _, pgIncident := range pgIncidents {
		inc := p.pgIncidentToIncident(pgIncident)
		incidents = append(incidents, inc)
	}
	return incidents, nil
}

func (p *PostgresStore) GetByID(ctx context.Context, id string) (core.Incident, error) {
	pgIncident, err := p.Queries.GetByID(ctx, id)
	if err != nil {
		// the sqlc always fails there when adding new incidents because it did not found
		// I dont think failing to get an incident because it already exists is trully error because it says something valuable - that i can continue adding new incidents if this check fails.
		return core.Incident{}, err
	}
	inc := p.pgIncidentToIncident(pgIncident)
	return inc, nil
}

func (p *PostgresStore) DeleteByID(ctx context.Context, id string) error {
	err := p.Queries.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	// Get All incidents to satisfy function - rebuild report
	incidents, err := p.GetAll(ctx)
	if err != nil {
		return err
	}
	_, err = core.BuildReport(incidents)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) DeleteAll(ctx context.Context) error {
	// !This removes everything from the database forever!
	err := p.Queries.DeleteAll(ctx)
	if err != nil {
		return err
	}
	return nil
}