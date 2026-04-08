package exercises

import (
	"encoding/json"
	"fitness_bot/internal/domain"
	"io"
	"net/http"
)

type ExercisesClient struct {
	token string
	url   string
}

func (p *ExercisesClient) FetchExercises(r *http.Request) ([]domain.ExerciseRecord, error) {

	var exercises []domain.ExerciseRecord
	var err error

	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.token)

	buildQuery(req, r)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &exercises)
	if err != nil {
		return nil, err
	}

	return exercises, nil
}

func buildQuery(req *http.Request, r *http.Request) {
	if r == nil {
		return
	}

	q := req.URL.Query()
	if q != nil {
		for _, mg := range r.URL.Query()["muscle_groups[]"] {
			q.Add("muscle_groups[]", mg)
		}

		if v := r.URL.Query().Get("limit"); v != "" {
			q.Set("limit", v)
		}

		if v := r.URL.Query().Get("lang"); v != "" {
			q.Add("lang", v)
		}

		if v := r.URL.Query().Get("category"); v != "" {
			q.Add("category", v)
		}

		if v := r.URL.Query().Get("difficulty"); v != "" {
			q.Add("difficulty", v)
		}

		if v := r.URL.Query().Get("updated_since"); v != "" {
			q.Set("updated_since", v)
		}

		req.URL.RawQuery = q.Encode()
	}
}

func NewExercisesClient(token, url string) *ExercisesClient {
	return &ExercisesClient{token: token, url: url}
}
