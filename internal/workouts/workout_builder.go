package workouts

import (
	"fitness_bot/internal/exercises"
	exercisespkg "fitness_bot/internal/exercises"
	"fmt"
	"math/rand"
	"net/http"
)

type WorkoutBuilder struct {
	client exercises.ExercisesFetcher
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

	var exercisesSize = 3
	if len(exercises) < exercisesSize {
		exercisesSize = len(exercises)
	}

	workout.Exercises = b.PickUniqueExercises(exercises, exercisesSize)

	workout.Type = b.DefineType(workout.Exercises)
	workout.Difficulty = b.DefineDifficulty(workout.Exercises)
	return workout, nil
}

func NewWorkoutBuilder(client exercises.ExercisesFetcher) *WorkoutBuilder {
	return &WorkoutBuilder{client: client}
}

func (b *WorkoutBuilder) PickUniqueExercises(exercises []exercisespkg.ExerciseRecord, count int) []exercises.ExerciseRecord {
	uniqueExercises := make(map[string]exercisespkg.ExerciseRecord)
	for len(uniqueExercises) < count && len(uniqueExercises) < len(exercises) {
		randomIndex := rand.Intn(len(exercises))
		pick := exercises[randomIndex]

		uniqueExercises[pick.UUID] = pick
	}

	result := make([]exercisespkg.ExerciseRecord, 0, len(uniqueExercises))
	for _, exercise := range uniqueExercises {
		result = append(result, exercise)
	}
	return result
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
