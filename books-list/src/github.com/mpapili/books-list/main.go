package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    "json:id"
	Title  string "json:tite"
	Author string "json:author"
	Year   string "json:year"
}

var books []Book

func main() {

	books = append(books, Book{ID: 1, Title: "Golang Pointers", Author: "Mr. Golang", Year: "2020"},
		Book{ID: 2, Title: "the adventuers of mike", Author: "Mike", Year: "1992"})

	// create a new mux router object
	router := mux.NewRouter()
	// path /books runs getBooks() and accepts method GET
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}

func getBooks(w http.ResponseWriter, r *http.Request) {

	log.Println("Get all books is called")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {

	log.Println("Get single book is called")
	params := mux.Vars(r)
	log.Println(params)

	reflect.TypeOf(params["id"])

	for _, book := range books {
		id_as_int, err := strconv.Atoi(params["id"])
		if err != nil {
			panic("could not cast ID to integer")
		}
		if book.ID == id_as_int {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {

	var book Book                         // new blank book object to be where r.Body is decoded
	json.NewDecoder(r.Body).Decode(&book) // decode request-body into book
	books = append(books, book)           // add books to book

	// response to client with encoded json of books slice
	json.NewEncoder(w).Encode(books)
	log.Println("add book is called")
}

func updateBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	json.NewDecoder(r.Body).Decode(&book) // "book" becomes PUT book

	// replace book with matching ID with our new book
	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}
	// return books slice
	json.NewEncoder(w).Encode(books)
	log.Println("update book is called")
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("remove book is called")
	// get our query params
	params := mux.Vars(r)
	id_as_int, err := strconv.Atoi(params["id"])

	if err != nil {
		panic("could not cast id to integer")
	}

	for i, item := range books {
		if item.ID == id_as_int {
			// everything up to index i + everything beyond index i
			books = append(books[:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
