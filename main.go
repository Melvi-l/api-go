package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Comment struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Content   string  `json:"content"`
	CreatedAt []uint8 `json:"created_at"`
}

var db *sql.DB

func buildDataSourceName() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func connectToDatabase() {
	var err error
	db, err = sql.Open("mysql", buildDataSourceName())
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connexion succeed !")
	return
}

func main() {
	fmt.Println("Hello, Go server API !")

	connectToDatabase()
	defer db.Close()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go client API !")
	})

	commentRouter := http.NewServeMux()
	router.Handle("/comments/", http.StripPrefix("/comments", commentRouter))
	commentRouter.HandleFunc("GET /", getCommentsHandler)
	commentRouter.HandleFunc("GET /{id}", getCommentsByIdHandler)
	commentRouter.HandleFunc("POST /", postCommentsHandler)
	commentRouter.HandleFunc("PUT /{id}", putCommentsByIdHandler)
	commentRouter.HandleFunc("DELETE /{id}", deleteComments)

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Query db
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse rows
	var comments []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.Username, &c.Content, &c.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Serialized to JSON
	jsonData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		fmt.Println(err.Error())
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getCommentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse path value
	id := r.PathValue("id")

	// Query db
	c := Comment{}
	err := db.QueryRow("SELECT id, username, content, created_at FROM comments WHERE id=?", id).Scan(&c.ID, &c.Username, &c.Content, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error when getting comment: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Serialized to JSON
	jsonResponse, err := json.Marshal(c)
	if err != nil {
		log.Printf("Error when marshaling JSON: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func postCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var c Comment
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert in db
	_, err = db.Exec("INSERT INTO comments (username, content) VALUES (?, ?)", c.Username, c.Content)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func putCommentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse path value
	id := r.PathValue("id")

	// Parse request body
	var c Comment
	err := json.NewDecoder(r.Body).Decode(&c)

	// Query db
	_, err = db.Exec("UPDATE comments SET content = ? WHERE id = ?", c.Content, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteComments(w http.ResponseWriter, r *http.Request) {
	// Parse path value
	id := r.PathValue("id")

	// Query db
	_, err := db.Exec("DELETE FROM comments WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
