package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book struct
type Book struct {
	ID     string `json:"id"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Author `json:"author"`
}

var books []Book

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Encode the books slice as JSON as write to response
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get route params
	params := mux.Vars(r)

	// Loop over all books
	for _, b := range books {
		if b.ID == params["id"] {
			json.NewEncoder(w).Encode(b)
			return
		}
	}

	json.NewEncoder(w).Encode(Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book

	// Decode JSON book in request body and write to the variable book
	json.NewDecoder(r.Body).Decode(&book)

	// Set random ID
	book.ID = strconv.Itoa(rand.Intn(1000000))

	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	for i, b := range books {
		if b.ID == params["id"] {
			// Remove the current version of the book
			books = append(books[:i], books[i+1:]...)
			// Add the updated version of the book
			books = append(books, book)
		}
	}

	json.NewEncoder(w).Encode(books)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for i, book := range books {
		if book.ID == params["id"] {

			// One method
			books = append(books[:i], books[i+1:]...)

			// Alternate method
			// books[len(books)-1], books[i] = books[i], books[len(books)-1]
			// books = books[:len(books)-1]

			json.NewEncoder(w).Encode(book)
			return
		}
	}

}

func main() {

	// Initialize router
	r := mux.NewRouter()

	// Set up dummy data
	books = append(books, Book{ID: "1", Isbn: "44872", Title: "Book One", Author: Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "44871", Title: "Book Two", Author: Author{Firstname: "John", Lastname: "No"}})

	// Route handlers
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start server; Log errors
	log.Fatal(http.ListenAndServe(":8000", r))
}
