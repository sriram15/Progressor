package connection

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sriram15/progressor-todo-app/internal/database"
)

// DBManager manages database transactions and connections.
type DBManager struct {
	db *sql.DB
}

// NewDBManager creates a new DBManager.
func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{db: db}
}

// Execute runs a function within a database transaction.
// It automatically handles commit and rollback.
func (m *DBManager) Execute(ctx context.Context, fn func(q *database.Queries) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	q := database.New(tx)

	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// Queries returns a new Queries object for non-transactional reads.
func (m *DBManager) Queries(ctx context.Context) *database.Queries {
	return database.New(m.db)
}
