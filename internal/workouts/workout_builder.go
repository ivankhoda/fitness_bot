package workouts

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/exercises"
	"fmt"
	"net/http"
)

type WorkoutBuilder struct {
	client exercises.ExercisesFetcher
	app    *config.Application
}

func (b *WorkoutBuilder) BuildWorkout(r *http.Request) (Workout, error) {
	var workout Workout
	var err error

	exercises, err := b.client.FetchExercises(r)

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

func NewWorkoutBuilder(client exercises.ExercisesFetcher, app config.Application) *WorkoutBuilder {
	return &WorkoutBuilder{client: client, app: &app}
}

func (b *WorkoutBuilder) DefineType(exercises []exercises.ExerciseRecord) string {
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
func (b *WorkoutBuilder) DefineDifficulty(exercises []exercises.ExerciseRecord) string {
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
