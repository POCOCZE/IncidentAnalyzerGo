-- Basic incidents table
CREATE TABLE incidents (
    id VARCHAR(50) PRIMARY KEY,
    title TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ,
    is_resolved BOOLEAN
);