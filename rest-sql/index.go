package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	DB_USER     = "gouser"
	DB_PASSWORD = "gopassword"
	DB_NAME     = "gotest"
	DB_URL      = "127.0.0.1:3306"
	//DB_URL = "localhost:3306"
)

// DB set up
func setupDB() *sql.DB {
	//gouser:gopassword@tcp(127.0.0.1:3306)/gotest
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASSWORD, DB_URL, DB_NAME)
	//printMessage(dbinfo)
	db, err := sql.Open("mysql", dbinfo)

	checkErr(err)

	return db
}

func checkSqlVesion() {
	db := setupDB()
	defer db.Close()

	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	printMessage("Connected to:" + version)
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

// Main function
func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all movies
	router.HandleFunc("/movies/", GetMovies).Methods("GET")

	// Get a specific movie by the movieID
	router.HandleFunc("/movies/{movieid}", GetMovieByID).Methods("GET")

	// Create a movie
	router.HandleFunc("/movies/", CreateMovie).Methods("POST")

	// Delete a specific movie by the movieID
	router.HandleFunc("/movies/{movieid}", DeleteMovie).Methods("DELETE")

	// Delete all movies
	router.HandleFunc("/movies/", DeleteMovies).Methods("DELETE")

	// serve the app
	fmt.Println("Server listen at 8080")
	checkSqlVesion()
	log.Print(http.ListenAndServe(":8080", router))
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all movies

// response and request handlers
func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	defer db.Close()

	printMessage("Getting movies...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var movies []Movie

	// Foreach movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		// check errors
		checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

	json.NewEncoder(w).Encode(response)
}

// Get a movie

// response and request handlers
func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	printMessage("Getting movie by ID")

	params := mux.Vars(r)

	movieID := params["movieid"]
	printMessage("movieid = " + movieID)

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()
		defer db.Close()

		printMessage("Get a movie from DB")

		rows, err := db.Query("SELECT * FROM movies where movieID = ?", movieID)

		// check errors
		checkErr(err)

		// var response []JsonResponse
		var movies []Movie

		// Foreach movie
		for rows.Next() {
			var id int
			var movieID string
			var movieName string

			err = rows.Scan(&id, &movieID, &movieName)

			// check errors
			checkErr(err)

			movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
		}

		response = JsonResponse{Type: "success", Data: movies}
	}

	json.NewEncoder(w).Encode(response)
}

// Create a movie

// response and request handlers
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("movieid")
	movieName := r.FormValue("moviename")

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()
		defer db.Close()

		printMessage("Inserting movie into DB")

		fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES(?, ?) returning id;", movieID, movieName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a movie

// response and request handlers
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	movieID := params["movieid"]

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()
		defer db.Close()

		printMessage("Deleting movie from DB")

		_, err := db.Exec("DELETE FROM movies where movieID = ?", movieID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all movies

// response and request handlers
func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	defer db.Close()

	printMessage("Deleting all movies...")

	_, err := db.Exec("DELETE FROM movies")

	// check errors
	checkErr(err)

	printMessage("All movies have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
