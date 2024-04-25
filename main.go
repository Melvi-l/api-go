package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Comment struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Content   string  `json:"content"`
	CreatedAt []uint8 `json:"created_at"`
}

func connectToDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:test1234@tcp(localhost:3306)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connexion succeed !")
	return db
}
func main() {
	fmt.Println("Hello, Go server API !")

	db := connectToDatabase()
	defer db.Close()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go client API !")
	})

	commentRouter := http.NewServeMux()
	router.Handle("/comments/", http.StripPrefix("/comments", commentRouter))
	commentRouter.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		comments, err := findAllComment(db)
		if err != nil {
			fmt.Println(err.Error())
		}
		rawJsonData, err := commentsToJson(comments)
          w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(rawJsonData)
	})
	commentRouter.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "post a new comment")
	})
	commentRouter.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "return the comment with id=%s", id)
	})

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func findAllComment(db *sql.DB) ([]Comment, error) {
	rows, err := db.Query("SELECT * FROM comment")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.Username, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func commentsToJson(comments []Comment) ([]byte, error) {
	jsonData, err := json.Marshal(comments)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
