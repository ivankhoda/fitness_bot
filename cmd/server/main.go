package main

import (
	"fitness_bot/internal/exercises"
	"fitness_bot/internal/workouts"
	"log"
	"net/http"
	"os"
)

func main() {
	// errEnv := godotenv.Load("../../.env")
	// if errEnv != nil {
	// 	fmt.Println("Error loading .env file")
	// 	return
	// }
	client := exercises.NewExercisesClient(os.Getenv("BOT_TOKEN"), os.Getenv("EXTERNAL_PROVIDER_URL")+"/api/api/exercises")

	workoutBuilder := workouts.NewWorkoutBuilder(client)
	workoutHandler := workouts.NewWorkoutHandler(workoutBuilder)

	http.Handle("/workouts", workoutHandler)
	log.Println("Server is running on port 5000...")
	http.ListenAndServe(":5000", nil)

}
