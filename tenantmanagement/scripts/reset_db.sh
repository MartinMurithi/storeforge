#!/bin/bash
# 1. Drop the migration table to start over
psql 'postgres://postgres:martin321!@localhost:5432/storeforge?sslmode=disable' -c "DROP TABLE IF EXISTS schema_migrations;"
# 2. Run migrations via the app
go run ./cmd/servergit reset
