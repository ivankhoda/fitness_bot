package exercises

type ExerciseRecord struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MuscleGroups []string `json:"muscle_groups"`
	Difficulty   string   `json:"difficulty"`
	Category     string   `json:"category"`
}
