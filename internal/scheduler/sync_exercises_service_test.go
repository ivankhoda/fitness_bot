package scheduler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"fitness_bot/internal/domain"
)

type fakeFetcher struct {
	response []domain.ExerciseRecord
	request  *http.Request
}

func (f *fakeFetcher) FetchExercises(r *http.Request) ([]domain.ExerciseRecord, error) {
	f.request = r
	return f.response, nil
}

type fakeExerciseRepo struct {
	upserted []domain.ExerciseRecord
}

func (f *fakeExerciseRepo) Insert(exercise domain.ExerciseRecord) (int, error) {
	panic("unexpected Insert call")
}

func (f *fakeExerciseRepo) Upsert(exercise domain.ExerciseRecord) error {
	f.upserted = append(f.upserted, exercise)
	return nil
}

func (f *fakeExerciseRepo) GetAll(filter domain.ExercsiesFilter) ([]domain.ExerciseRecord, error) {
	panic("unexpected GetAll call")
}

func (f *fakeExerciseRepo) GetByID(id int) (*domain.ExerciseRecord, error) {
	panic("unexpected GetByID call")
}

func (f *fakeExerciseRepo) Delete(id int) error {
	panic("unexpected Delete call")
}

type fakeSyncStateStore struct {
	checkpoint *time.Time
	saved      *time.Time
}

func (f *fakeSyncStateStore) GetExercisesCheckpoint() (*time.Time, error) {
	return f.checkpoint, nil
}

func (f *fakeSyncStateStore) SaveExercisesCheckpoint(syncedAt time.Time) error {
	f.saved = &syncedAt
	return nil
}

func TestSyncExercisesServiceRunUsesCheckpointAndStoresLatestTimestamp(t *testing.T) {
	checkpoint := time.Date(2026, 4, 7, 9, 0, 0, 0, time.UTC)
	updatedAt1 := time.Date(2026, 4, 7, 10, 0, 0, 0, time.UTC)
	updatedAt2 := time.Date(2026, 4, 7, 11, 0, 0, 0, time.UTC)

	fetcher := &fakeFetcher{response: []domain.ExerciseRecord{
		{UUID: "1", Name: "Push-up", UpdatedAt: &updatedAt1},
		{UUID: "2", Name: "Pull-up", UpdatedAt: &updatedAt2},
	}}
	repo := &fakeExerciseRepo{}
	state := &fakeSyncStateStore{checkpoint: &checkpoint}
	service := &SyncExercisesService{
		fetcher: fetcher,
		repo:    repo,
		state:   state,
	}

	if err := service.Run(); err != nil {
		t.Fatal(err)
	}

	if fetcher.request == nil {
		t.Fatal("expected request to be passed to fetcher")
	}

	if got := fetcher.request.URL.Query().Get("updated_since"); got != checkpoint.Format(time.RFC3339) {
		t.Fatalf("expected updated_since query %q, got %q", checkpoint.Format(time.RFC3339), got)
	}

	if len(repo.upserted) != 2 {
		t.Fatalf("expected 2 upserts, got %d", len(repo.upserted))
	}
	if state.saved == nil || !state.saved.Equal(updatedAt2) {
		t.Fatalf("expected latest checkpoint %s, got %+v", updatedAt2.Format(time.RFC3339), state.saved)
	}
}

type failingSyncStateStore struct {
	err error
}

func (f *failingSyncStateStore) GetExercisesCheckpoint() (*time.Time, error) {
	return nil, f.err
}

func (f *failingSyncStateStore) SaveExercisesCheckpoint(syncedAt time.Time) error {
	return f.err
}

func TestSyncExercisesServiceRunFallsBackWhenCheckpointStoreUnavailable(t *testing.T) {
	updatedAt := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	fetcher := &fakeFetcher{response: []domain.ExerciseRecord{{UUID: "1", Name: "Push-up", UpdatedAt: &updatedAt}}}
	repo := &fakeExerciseRepo{}
	state := &failingSyncStateStore{err: errors.New("ERROR: permission denied for schema public (SQLSTATE 42501)")}
	service := &SyncExercisesService{
		fetcher: fetcher,
		repo:    repo,
		state:   state,
	}

	if err := service.Run(); err != nil {
		t.Fatal(err)
	}

	if fetcher.request != nil {
		t.Fatalf("expected full sync request without updated_since, got %+v", fetcher.request.URL)
	}

	if len(repo.upserted) != 1 {
		t.Fatalf("expected 1 upsert, got %d", len(repo.upserted))
	}
}
