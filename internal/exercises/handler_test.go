package exercises

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeClient struct{}

func (f *FakeClient) FetchExercises(r *http.Request) ([]ExerciseRecord, error) {
	return []ExerciseRecord{
		{UUID: "1", Name: "Push-up", MuscleGroups: []string{"chest"}, Difficulty: "beginner", Category: "strength"},
		{UUID: "2", Name: "Pull-up", MuscleGroups: []string{"back"}, Difficulty: "beginner", Category: "strength"},
	}, nil
}

func TestExercisesHandler(t *testing.T) {

	w := httptest.NewRecorder()
	fakeClient := &FakeClient{}
	handler := NewExercisesHandler(fakeClient)

	req := newGetExerciseRequest([]string{"chest", "back"}, "beginner", "strength")

	handler.GetExercises(w, req)
	log.Printf("Request URL: %s", req.URL.String())
	log.Printf("Received response: %s", w.Body.String())
	exercises := getExercisesFromResponse(t, w)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := []ExerciseRecord{
		{UUID: "1", Name: "Push-up", MuscleGroups: []string{"chest"}, Difficulty: "beginner", Category: "strength"},
		{UUID: "2", Name: "Pull-up", MuscleGroups: []string{"back"}, Difficulty: "beginner", Category: "strength"},
	}

	if len(exercises) != len(expected) {
		t.Fatalf("expected %d exercises, got %d", len(expected), len(exercises))
	}

	for i, exercise := range exercises {
		log.Printf("Asserting exercise %d: got %+v, want %+v", i, exercise, expected[i])
		assertResponseBody(t, exercise, expected[i])
	}

}

func newGetExerciseRequest(muscle_groups []string, difficulty string, category string) *http.Request {
	muscleGroupParams := strings.Join(muscle_groups, "&muscle_groups[]=")
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/exercises/?muscle_groups[]=%s&difficulty=%s&category=%s", muscleGroupParams, difficulty, category), nil)
	return req
}

func getExercisesFromResponse(t testing.TB, response *httptest.ResponseRecorder) []ExerciseRecord {
	t.Helper()
	var exercises []ExerciseRecord
	err := json.Unmarshal(response.Body.Bytes(), &exercises)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Exercise, '%v'", response.Body.String(), err)
	}
	log.Printf("Parsed exercises from response: %+v", exercises)
	return exercises
}

func assertResponseBody(t testing.TB, got, want ExerciseRecord) {
	t.Helper()
	if got.UUID != want.UUID || got.Name != want.Name || got.Difficulty != want.Difficulty || got.Category != want.Category || !slicesEqual(got.MuscleGroups, want.MuscleGroups) {
		t.Errorf("response body is wrong, got %+v want %+v", got, want)
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
