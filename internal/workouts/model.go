package workouts

import "fitness_bot/internal/exercises"

type Workout struct {
	Type       string                     `json:"type"`
	Difficulty string                     `json:"difficulty"`
	Exercises  []exercises.ExerciseRecord `json:"exercises"`
}
