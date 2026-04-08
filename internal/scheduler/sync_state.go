package scheduler

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const exercisesSyncKey = "exercises"

type SyncCheckpointStore interface {
	GetExercisesCheckpoint() (*time.Time, error)
	SaveExercisesCheckpoint(syncedAt time.Time) error
}

type SyncStateStore struct {
	DB *pgxpool.Pool
}

func NewSyncStateStore(db *pgxpool.Pool) *SyncStateStore {
	return &SyncStateStore{DB: db}
}

func (s *SyncStateStore) GetExercisesCheckpoint() (*time.Time, error) {
	var syncedAt time.Time
	err := s.DB.QueryRow(context.Background(), "SELECT synced_at FROM sync_state WHERE sync_key = $1", exercisesSyncKey).Scan(&syncedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &syncedAt, nil
}

func (s *SyncStateStore) SaveExercisesCheckpoint(syncedAt time.Time) error {
	_, err := s.DB.Exec(context.Background(), `
		INSERT INTO sync_state (sync_key, synced_at)
		VALUES ($1, $2)
		ON CONFLICT (sync_key)
		DO UPDATE SET synced_at = EXCLUDED.synced_at`, exercisesSyncKey, syncedAt.UTC())
	return err
}
