package repository

import(
	"context"

)

// RowScanner interface used by repo
type MockRow struct {
    ScanErr error
}

func (r *MockRow) Scan(dest ...interface{}) error {
    return r.ScanErr
}

// Mock DB implements only QueryRow
type MockDB struct {
    RowErr error
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) repository.RowScanner {
    return &MockRow{ScanErr: m.RowErr}
}