package exercises

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"fitness_bot/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseModel struct {
	DB *pgxpool.Pool
}

func (e *ExerciseModel) Insert(exercise domain.ExerciseRecord) (int, error) {

	stmt := `INSERT INTO exercises (external_uuid, name, description, muscle_groups, difficulty, category) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err := e.DB.QueryRow(context.Background(), stmt, exercise.UUID, exercise.Name, exercise.Description, exercise.MuscleGroups, exercise.Difficulty, exercise.Category).Scan(&id)
	if err != nil {
		log.Printf("Error inserting exercise: %v", err)
		return 0, err
	}
	return id, nil
}

func (e *ExerciseModel) Upsert(exercise domain.ExerciseRecord) error {
	stmt := `
		INSERT INTO exercises (external_uuid, name, description, muscle_groups, difficulty, category)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (external_uuid)
		DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			muscle_groups = EXCLUDED.muscle_groups,
			difficulty = EXCLUDED.difficulty,
			category = EXCLUDED.category`

	_, err := e.DB.Exec(context.Background(), stmt, exercise.UUID, exercise.Name, exercise.Description, exercise.MuscleGroups, exercise.Difficulty, exercise.Category)
	if err != nil {
		log.Printf("Error upserting exercise: %v", err)
		return err
	}

	return nil
}

func (e *ExerciseModel) GetAll(f domain.ExercsiesFilter) ([]domain.ExerciseRecord, error) {
	var exercises []domain.ExerciseRecord
	query := "SELECT external_uuid, name, description, muscle_groups, difficulty, category FROM exercises WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if len(f.MuscleGroups) > 0 {
		placeholders := make([]string, len(f.MuscleGroups))
		for i, mg := range f.MuscleGroups {
			placeholders[i] = "$" + fmt.Sprintf("%d", argIndex)
			args = append(args, mg)
			argIndex++
		}
		query += " AND muscle_groups && ARRAY[" + strings.Join(placeholders, ", ") + "]::text[]"
	}
	if f.Category != "" {
		query += " AND category = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, f.Category)
		argIndex++
	}
	if f.Difficulty != "" {
		query += " AND difficulty = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, f.Difficulty)
		argIndex++
	}
	query += " ORDER BY RANDOM()"
	limit, err := strconv.Atoi(f.Limit)
	if err == nil && limit > 0 {
		query += " LIMIT " + fmt.Sprintf("%d", limit)
	}

	rows, err := e.DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error fetching exercises: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var exercise domain.ExerciseRecord
		err := rows.Scan(&exercise.UUID, &exercise.Name, &exercise.Description, &exercise.MuscleGroups, &exercise.Difficulty, &exercise.Category)
		if err != nil {
			log.Printf("Error scanning exercise: %v", err)
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating over exercises: %v", rows.Err())
		return nil, rows.Err()
	}

	return exercises, nil
}

func (e *ExerciseModel) GetByID(id int) (*domain.ExerciseRecord, error) {
	var exercise domain.ExerciseRecord
	err := e.DB.QueryRow(context.Background(), "SELECT external_uuid, name, description, muscle_groups, difficulty, category FROM exercises WHERE id = $1", id).Scan(&exercise.UUID, &exercise.Name, &exercise.Description, &exercise.MuscleGroups, &exercise.Difficulty, &exercise.Category)
	if err != nil {
		log.Printf("Error fetching exercise by ID: %v", err)
		return nil, err
	}
	return &exercise, nil
}

func (e *ExerciseModel) GetByUUID(uuid string) (*domain.ExerciseRecord, error) {
	var exercise domain.ExerciseRecord
	err := e.DB.QueryRow(context.Background(), "SELECT external_uuid, name, description, muscle_groups, difficulty, category FROM exercises WHERE external_uuid = $1", uuid).Scan(&exercise.UUID, &exercise.Name, &exercise.Description, &exercise.MuscleGroups, &exercise.Difficulty, &exercise.Category)
	if err != nil {
		log.Printf("Error fetching exercise by UUID: %v", err)
		return nil, err
	}
	return &exercise, nil
}

func (e *ExerciseModel) Update(id int, exercise domain.ExerciseRecord) error {
	return nil
}

func (e *ExerciseModel) Delete(id int) error {
	return nil
}
