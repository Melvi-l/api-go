package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestAPI(t *testing.T) {
	err := godotenv.Load(".env.dev")
	if err != nil {
		fmt.Println("Error loading .env file: " + err.Error())
	}

	ctx := context.Background()

	dbName := os.Getenv("MYSQL_DATABASE")

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.3.0"),
		mysql.WithDatabase(dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("test1234"),
		mysql.WithScripts("db/test_schema.sql"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	defer mysqlContainer.Terminate(ctx)

	connectionString, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}
	fmt.Println(connectionString)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		t.Fatalf("Failed to open database: %s", err)
	}
	defer db.Close()

	t.Run("Test GET /comments", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/comments", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getCommentsHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Response should be OK")

		var comments []Comment
		err := json.Unmarshal(rr.Body.Bytes(), &comments)
		assert.NoError(t, err, "Response should be a valid JSON of comments")
		assert.NotEmpty(t, comments, "Response should not be empty")
	})

	t.Run("Test GET /comments/{id}", func(t *testing.T) {
		router := http.NewServeMux()
		router.HandleFunc("GET /comments/{id}", getCommentsByIdHandler)

		server := httptest.NewServer(router)
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL+"/comments/1", nil)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Failed to send GET request:", err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode, "Response should be OK")

		body, err := io.ReadAll(response.Body)
		assert.NoError(t, err, "Response body should be readable")

		var c Comment
		err = json.Unmarshal(body, &c)
		assert.NoError(t, err, "Response should be a valid JSON of comments")

		assert.Equal(t, c.Username, "test_container_user", "Comment 1 user should be test_container_user")
	})

	t.Run("Test POST /comments", func(t *testing.T) {
		newComment := Comment{Username: "testuser", Content: "A new test comment"}
		jsonComment, _ := json.Marshal(newComment)

		req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonComment))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(postCommentsHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code, "Response should be Created")
	})

	t.Run("Test PUT /comments/{id}", func(t *testing.T) {
		updatedComment := Comment{Content: "Updated content"}
		jsonComment, _ := json.Marshal(updatedComment)

		router := http.NewServeMux()
		router.HandleFunc("PUT /comments/{id}", putCommentsByIdHandler)

		server := httptest.NewServer(router)
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/comments/1", bytes.NewBuffer(jsonComment))
		req.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Failed to send PUT request:", err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode, "Response should be OK after updating")

		currentComment := Comment{}
		err = db.QueryRow("SELECT username, content FROM comments WHERE id=1").Scan(&currentComment.Username, &currentComment.Content)
		if err != nil {
			log.Printf("Error when getting comment: %v", err)
			return
		}

		assert.Equal(t, updatedComment.Content, currentComment.Content, "Comment content should be updated")
	})

	t.Run("Test DELETE /comments/{id}", func(t *testing.T) {
		router := http.NewServeMux()
		router.HandleFunc("DELETE /comments/{id}", deleteCommentsByIdHandler)

		server := httptest.NewServer(router)
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/comments/1", nil)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Failed to send DELETE request:", err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode, "Response should be OK after deleting")

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE id=1").Scan(&count)
		if err != nil {
			log.Printf("Error when checking if comment exists: %v", err)
			return
		}

		assert.Equal(t, 0, count, "Comment should no longer exist in the database")
	})
}
