package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

type Book struct  {
	ID	int	"json:id"
	Title	string	"json:tite"
	Author	string	"json:author"
	Year	string	"json:year"
}

var books []Book

func main() {

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
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single book is called")
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("add book is called")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("update book is called")
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("remove book is called")
}

