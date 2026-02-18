package workouts

import (
	"fitness_bot/internal/exercises"
	"testing"
)

func TestWorkout(t *testing.T) {
	// Create a sample workout
	workout := Workout{
		Type:       "Cardio",
		Difficulty: "Medium",
		Exercises: []exercises.ExerciseRecord{
			{
				Name:         "Running",
				Description:  "Run at a steady pace for 30 minutes.",
				MuscleGroups: []string{"chest"},
			},
			{
				Name:         "Jumping Jacks",
				Description:  "Perform jumping jacks for 1 minute.",
				MuscleGroups: []string{"chest"},
			},
		},
	}

	if workout.Type != "Cardio" {
		t.Errorf("Expected workout type 'Cardio', got '%s'", workout.Type)
	}
	if workout.Difficulty != "Medium" {
		t.Errorf("Expected workout difficulty 'Medium', got '%s'", workout.Difficulty)
	}
	if len(workout.Exercises) != 2 {
		t.Errorf("Expected 2 exercises, got %d", len(workout.Exercises))
	}
	if workout.Exercises[0].Name != "Running" {
		t.Errorf("Expected first exercise name 'Running', got '%s'", workout.Exercises[0].Name)
	}
	if workout.Exercises[1].Name != "Jumping Jacks" {
		t.Errorf("Expected second exercise name 'Jumping Jacks', got '%s'", workout.Exercises[1].Name)
	}
}
