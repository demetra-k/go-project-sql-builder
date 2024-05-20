package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Book struct {
	id     int
	title  string
	ISBN   string
	author string
}

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createBookTable(db)

	data := []Book{}
	rows, err := db.Query("SELECT id, title, ISBN, author FROM book")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)

	var id int
	var title string
	var ISBN string
	var author string

	for rows.Next() {
		err := rows.Scan(&id, &title, &ISBN, &author)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Book{id, title, ISBN, author})
	}

	fmt.Println(data)

	// Call the updateBook function
	bookToUpdate := Book{
		id:     1, // ID of the book to update
		title:  "Updated Title",
		ISBN:   "1234567890",
		author: "Updated Author",
	}

	updatedID, err := updateBook(db, bookToUpdate)
	if err != nil {
		log.Fatalf("Failed to update book: %v", err)
	}

	fmt.Printf("Book with ID %d has been updated.\n", updatedID.id)

	// Call the insertBook function
	newBook := Book{
		title:  "New Book Title",
		ISBN:   "0987654321",
		author: "New Author",
	}

	insertedBook, err := insertBook(db, newBook)
	if err != nil {
		log.Fatalf("Failed to insert book: %v", err)
	}
	fmt.Printf("New book inserted: %+v\n", insertedBook)
}

func createBookTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS book (
		id serial primary key,
		title varchar(100) not null,
		ISBN  varchar(100) not null,
		author varchar(100) not null,
		yearOfIssue timestamp default now()
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func insertBook(db *sql.DB, book Book) (Book, error) {
	query := `INSERT INTO book (title, ISBN, author)
	VALUES ($1, $2, $3) RETURNING id, title, ISBN, author`

	var newBook Book
	err := db.QueryRow(query, book.title, book.ISBN, book.author).Scan(&newBook.id, &newBook.title, &newBook.ISBN, &newBook.author)

	if err != nil {
		log.Fatal(err)
	}
	return newBook, nil
}

func updateBook(db *sql.DB, updateBook Book) (Book, error) {
	query := `UPDATE book SET title = $1, ISBN = $2, author = $3 WHERE id = $4 RETURNING id, title, ISBN, author`

	var updatedBook Book
	err := db.QueryRow(query, updateBook.title, updateBook.ISBN, updateBook.author, updateBook.id).Scan(&updatedBook.id, &updatedBook.title, &updatedBook.ISBN, &updatedBook.author)
	if err != nil {
		return Book{}, err
	}
	return updatedBook, nil
}
