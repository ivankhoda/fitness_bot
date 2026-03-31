package workouts

import (
	"fitness_bot/internal/config"
	"fitness_bot/internal/exercises"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeClient struct {
}

func (f *FakeClient) FetchExercises(r *http.Request) ([]exercises.ExerciseRecord, error) {
	return []exercises.ExerciseRecord{
		{UUID: "1luid", Name: "Push-up", MuscleGroups: []string{"chest"}, Difficulty: "beginner", Category: "strength"},
		{UUID: "2xuid", Name: "Pull-up", MuscleGroups: []string{"back"}, Difficulty: "beginner", Category: "strength"},
	}, nil
}

func TestWorkoutBuilder(t *testing.T) {
	client := &FakeClient{}
	app := config.Application{}
	builder := NewWorkoutBuilder(client, app)
	request := httptest.NewRequest(http.MethodGet, "/generate_workout", nil)

	workout, err := builder.BuildWorkout(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if workout.Type != "strength" {
		t.Errorf("Expected workout type 'strength', got %s", workout.Type)
	}

	if workout.Difficulty != "beginner" {
		t.Errorf("Expected difficulty 'beginner', got %s", workout.Difficulty)
	}

	if len(workout.Exercises) != 2 {
		t.Errorf("Expected 2 exercises, got %d", len(workout.Exercises))
	}
}
