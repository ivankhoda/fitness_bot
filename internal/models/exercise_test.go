package models

import (
	"context"
	"fitness_bot/internal/assert"
	"fitness_bot/internal/domain"
	"fitness_bot/testutils"
	"testing"
)

func TestExerciseModelExists(t *testing.T) {
	db := testutils.NewTestDB(t)
	model := ExerciseModel{DB: db}

	exercise := domain.ExerciseRecord{
		UUID:         "test-uuid",
		Name:         "Test Exercise",
		MuscleGroups: []string{"chest"},
		Difficulty:   "beginner",
		Category:     "strength",
	}

	exists, err := model.Exists(exercise.UUID)
	assert.NilError(t, err)
	assert.Equal(t, exists, false)

	id, err := model.Insert(exercise)
	assert.NilError(t, err)

	exists, err = model.Exists(exercise.UUID)
	assert.NilError(t, err)
	assert.Equal(t, exists, true)

	_, err = db.Exec(context.Background(), "DELETE FROM exercises WHERE id = $1", id)
	assert.NilError(t, err)
}
