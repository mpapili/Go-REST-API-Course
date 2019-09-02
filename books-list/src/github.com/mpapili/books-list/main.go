package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"

	"database/sql"
)

type Book struct {
	ID     int    "json:id"
	Title  string "json:tite"
	Author string "json:author"
	Year   string "json:year"
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load() // load env variables in ".env" hidden file
}

func checkErr(err error) {
	// if error is not nill, log fatal error
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	/*
		books = append(books, Book{ID: 1, Title: "Golang Pointers", Author: "Mr. Golang", Year: "2020"},
			Book{ID: 2, Title: "the adventuers of mike", Author: "Mike", Year: "1992"})
	*/

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	checkErr(err)

	log.Println(pgUrl)

	db, err = sql.Open("postgres", pgUrl)
	checkErr(err)
	db.Ping() // returns nothing if successful

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

	var book Book // single book is of type book
	books = []Book{}

	rows, err := db.Query("SELECT * FROM books")
	checkErr(err)
	defer rows.Close() // close rows once this function is done

	for rows.Next() {
		/*
			map fields to book object
			each value is given an address
			(the args we pass in) where it is
			assigned
		*/
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		checkErr(err)
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {

	var book Book // single book is of type book

	params := mux.Vars(r)

	// parameterized query to get book with matching id
	row := db.QueryRow("SELECT * FROM books WHERE id=$1", params["id"])
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	checkErr(err)
	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	var bookID int

	// decode json body into "book" type object's address
	err := json.NewDecoder(r.Body).Decode(&book)
	checkErr(err)
	log.Println(book)

	err = db.QueryRow("INSERT INTO books (title, author, year) values ($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&bookID)
	// we create a book in the DB and return "ID" which we ".Scan()" to the
	// address of our bookID integer
	checkErr(err)
	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	json.NewDecoder(r.Body).Decode(&book) // "book" becomes PUT book

	result, err := db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)
	checkErr(err)
	rowsUpdated, err := result.RowsAffected() // see which rows were effected
	checkErr(err)
	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("remove book is called")
	// get our query params
	params := mux.Vars(r)
	result, err := db.Exec("DELETE from books WHERE id = $1", params["id"])
	checkErr(err)
	rowsDeleted, err := result.RowsAffected()

	json.NewEncoder(w).Encode(rowsDeleted)

}
