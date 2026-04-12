package workouts

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/domain"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type WorkoutBuilder struct {
	client domain.ExercisesFetcher
	app    *config.Application
}

func (b *WorkoutBuilder) BuildWorkout(r *http.Request) (Workout, error) {
	var workout Workout

	buildQuery(r)
	var muscleGroups []string
	for _, mg := range r.URL.Query()["muscle_groups[]"] {
		for _, m := range strings.Split(mg, ",") {
			if m = strings.TrimSpace(m); m != "" {
				muscleGroups = append(muscleGroups, m)
			}
		}
	}
	exercises, err := b.app.Exercise.GetAll(domain.ExercsiesFilter{
		MuscleGroups: muscleGroups,
		Category:     r.URL.Query().Get("category"),
		Difficulty:   r.URL.Query().Get("difficulty"),
		Limit:        r.URL.Query().Get("limit"),
	})

	if err != nil {
		return Workout{}, err
	}
	if len(exercises) == 0 {
		return Workout{}, fmt.Errorf("no exercises available")
	}

	workout.Exercises = exercises
	workout.Type = b.DefineType(exercises)
	workout.Difficulty = b.DefineDifficulty(exercises)
	return workout, nil
}

func NewWorkoutBuilder(client domain.ExercisesFetcher, app config.Application) *WorkoutBuilder {
	return &WorkoutBuilder{client: client, app: &app}
}

func (b *WorkoutBuilder) DefineType(exercises []domain.ExerciseRecord) string {
	mapCategories := make(map[string]int)
	for _, exercise := range exercises {
		mapCategories[exercise.Category]++
	}

	var maxCategory string
	var maxCount int
	for category, count := range mapCategories {
		if count > maxCount {
			maxCategory = category
			maxCount = count
		}
	}

	return maxCategory
}

func (b *WorkoutBuilder) DefineDifficulty(exercises []domain.ExerciseRecord) string {
	mapDifficulty := make(map[string]int)
	for _, exercise := range exercises {
		if exercise.Difficulty == "advanced" {
			return "advanced"
		}
		mapDifficulty[exercise.Difficulty]++
	}

	var maxDifficulty string
	var maxCount int
	for difficulty, count := range mapDifficulty {
		if count > maxCount {
			maxDifficulty = difficulty
			maxCount = count
		}
	}

	return maxDifficulty
}

func buildQuery(req *http.Request) {
	if req == nil {
		return
	}

	q := req.URL.Query()
	if q != nil {
		for _, mg := range req.URL.Query()["muscle_groups[]"] {
			q.Add("muscle_groups[]", mg)
		}

		if v := req.URL.Query().Get("limit"); v != "" {
			q.Set("limit", v)
		} else {
			q.Set("limit", "3")
		}

		if v := req.URL.Query().Get("lang"); v != "" {
			q.Add("lang", v)
		}

		if v := req.URL.Query().Get("category"); v != "" {
			q.Add("category", v)
		}

		if v := req.URL.Query().Get("difficulty"); v != "" {
			q.Add("difficulty", v)
		}
		log.Printf("Built query: %s", q.Encode())
		req.URL.RawQuery = q.Encode()
	}
}
