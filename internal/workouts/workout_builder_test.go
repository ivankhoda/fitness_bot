package workouts

import (
	"fitness_bot/internal/exercises"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeClient struct {
}

func (f *FakeClient) FetchExercises(r *http.Request) ([]exercises.ExerciseRecord, error) {
	return []exercises.ExerciseRecord{
		{UUID: "1", Name: "Push-up", MuscleGroups: []string{"chest"}, Difficulty: "beginner", Category: "strength"},
		{UUID: "2", Name: "Pull-up", MuscleGroups: []string{"back"}, Difficulty: "beginner", Category: "strength"},
	}, nil
}

func TestWorkoutBuilder(t *testing.T) {
	client := &FakeClient{}
	builder := NewWorkoutBuilder(client)
	request := httptest.NewRequest(http.MethodGet, "/workouts", nil)

	workout, err := builder.BuildWorkout(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if workout.Type != "Cardio" {
		t.Errorf("Expected workout type 'Cardio', got %s", workout.Type)
	}

	if workout.Difficulty != "Medium" {
		t.Errorf("Expected difficulty 'Medium', got %s", workout.Difficulty)
	}

	if len(workout.Exercises) != 2 {
		t.Errorf("Expected 2 exercises, got %d", len(workout.Exercises))
	}
}
