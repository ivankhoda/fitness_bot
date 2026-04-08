package scheduler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fitness_bot/internal/domain"
)

type SyncExercisesService struct {
	fetcher  domain.ExercisesFetcher
	repo     domain.ExerciseRepository
	state    SyncCheckpointStore
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewSyncExercisesService(fetcher domain.ExercisesFetcher, repo domain.ExerciseRepository, state SyncCheckpointStore, infoLog, errorLog *log.Logger) *SyncExercisesService {
	return &SyncExercisesService{
		fetcher:  fetcher,
		repo:     repo,
		state:    state,
		infoLog:  infoLog,
		errorLog: errorLog,
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

	if service.infoLog != nil {
		service.infoLog.Printf("starting exercise sync, since=%v", since)
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
		service.infoLog.Println(exercise.UUID, exercise.Name, exercise.UpdatedAt)
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

	if service.infoLog != nil {
		service.infoLog.Printf("exercise sync completed, processed=%d", len(exercises))
	}

	return nil
}

func (service *SyncExercisesService) logError(err error) {
	if service.errorLog != nil {
		service.errorLog.Printf("exercise sync failed: %v", err)
	}
}

func (service *SyncExercisesService) logWarning(format string, args ...any) {
	if service.infoLog != nil {
		service.infoLog.Printf(format, args...)
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
