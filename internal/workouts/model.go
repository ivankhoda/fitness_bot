package workouts

import "fitness_bot/internal/domain"

type Workout struct {
	Type       string                  `json:"type"`
	Difficulty string                  `json:"difficulty"`
	Exercises  []domain.ExerciseRecord `json:"exercises"`
}
