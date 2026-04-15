package scheduler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"fitness_bot/internal/config"
	"fitness_bot/internal/domain"
)

type SyncExercisesService struct {
	fetcher domain.ExercisesFetcher
	repo    domain.ExerciseRepository
	state   SyncCheckpointStore
	app     *config.Application
}

func NewSyncExercisesService(fetcher domain.ExercisesFetcher, repo domain.ExerciseRepository, state SyncCheckpointStore, app *config.Application) *SyncExercisesService {
	return &SyncExercisesService{
		fetcher: fetcher,
		repo:    repo,
		state:   state,
		app:     app,
	}
}

func (service *SyncExercisesService) Run() error {
	since, err := service.state.GetExercisesCheckpoint()
	if err != nil {
		if !service.canProceedWithoutCheckpoint(err) {
			service.logError(err)
			return err
		}

		service.logWarning("checkpoint store unavailable; falling back to full sync: %v", err)
		since = nil
	}

	if service.app.InfoLog != nil {
		service.app.InfoLog.Printf("starting exercise sync, since=%v", since)
	}

	exercises, err := service.fetcher.FetchExercises(newSyncRequest(since))
	if err != nil {
		service.logError(err)
		return err
	}

	var maxUpdatedAt *time.Time
	for _, exercise := range exercises {
		if err := service.repo.Upsert(exercise); err != nil {
			service.logError(err)
			return err
		}
		service.app.InfoLog.Println(exercise.UUID, exercise.Name, exercise.UpdatedAt)
		if exercise.UpdatedAt != nil && (maxUpdatedAt == nil || exercise.UpdatedAt.After(*maxUpdatedAt)) {
			updatedAt := *exercise.UpdatedAt
			maxUpdatedAt = &updatedAt
		}
	}

	if maxUpdatedAt != nil {
		if err := service.state.SaveExercisesCheckpoint(*maxUpdatedAt); err != nil {
			if !service.canProceedWithoutCheckpoint(err) {
				service.logError(err)
				return err
			}

			service.logWarning("checkpoint was not persisted; next sync may be full: %v", err)
		}
	}

	if service.app.InfoLog != nil {
		service.app.InfoLog.Printf("exercise sync completed, processed=%d", len(exercises))
	}

	return nil
}

func (service *SyncExercisesService) logError(err error) {
	if service.app.ErrorLog != nil {
		service.app.ErrorLog.Printf("exercise sync failed: %v", err)
	}
}

func (service *SyncExercisesService) logWarning(format string, args ...any) {
	if service.app.InfoLog != nil {
		service.app.InfoLog.Printf(format, args...)
	}
}

func (service *SyncExercisesService) canProceedWithoutCheckpoint(err error) bool {
	message := err.Error()
	return strings.Contains(message, "permission denied for schema") || strings.Contains(message, "relation \"sync_state\" does not exist")
}

func newSyncRequest(since *time.Time) *http.Request {
	if since == nil {
		return nil
	}

	query := url.Values{}
	query.Set("updated_since", since.UTC().Format(time.RFC3339))

	return &http.Request{
		URL: &url.URL{RawQuery: query.Encode()},
	}
}
