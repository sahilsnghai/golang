package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
)

var App_Json string = "application/json"
var Content_Type string = "Content-Type"

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    ` json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {
	var MoviesEndpoint string = "/movies"
	var MovieEndpoint string = "/movies/{id}"

	movies = append(movies, Movie{ID: "1", Isbn: "233248", Title: "Sab ke papa", Director: &Director{Firstname: "Sahil", Lastname: "Singhai"}})
	movies = append(movies, Movie{ID: "2", Isbn: "233249", Title: "Sab ke papa part 2", Director: &Director{Firstname: "Sahil", Lastname: "Singhai"}})

	r := mux.NewRouter()
	r.HandleFunc(MoviesEndpoint, getMovies).Methods("GET")
	r.HandleFunc(MovieEndpoint, getMovie).Methods("GET")
	r.HandleFunc(MoviesEndpoint, createMovies).Methods("POST")
	r.HandleFunc(MovieEndpoint, updateMovies).Methods("PUT")
	r.HandleFunc(MovieEndpoint, deleteMovies).Methods("DELETE")

	fmt.Println("Starting Server at port at server 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func getMovies(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Request at Get Movies with 's'")
	w.Header().Set(Content_Type, App_Json)
	json.NewEncoder(w).Encode(movies)

}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(Content_Type, App_Json)
	fmt.Println("Request at delete Movie")

	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(Content_Type, App_Json)
	params := mux.Vars(r)

	fmt.Println("Request at Get Movie")

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(Content_Type, App_Json)
	fmt.Println("Request at Create Movie")

	var movie Movie

	_ = json.NewDecoder(r.Body).Decode(&movie)

	Id, _ := rand.Int(rand.Reader, big.NewInt(1000000))

	movie.ID = Id.String()
	movies = append(movies, movie)

}

func updateMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", App_Json)
	fmt.Println("Request at Update Movie")

	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)

			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

}
