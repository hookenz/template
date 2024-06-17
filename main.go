package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var db *sql.DB

type Todo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	// Initialize zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to SQLite database
	var err error
	db, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// Create table if not exists
	createTable()

	// Routes
	e.POST("/todos", createTodo)
	e.GET("/todos", getTodos)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func createTable() {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        title TEXT,
        completed BOOLEAN
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create table")
	}
}

func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := "INSERT INTO todos (user_id, title, completed) VALUES (?, ?, ?)"
	result, err := db.Exec(query, todo.UserID, todo.Title, todo.Completed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	id, _ := result.LastInsertId()
	todo.ID = int(id)
	return c.JSON(http.StatusCreated, todo)
}

func getTodos(c echo.Context) error {
	rows, err := db.Query("SELECT id, user_id, title, completed FROM todos")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		todos = append(todos, todo)
	}
	return c.JSON(http.StatusOK, todos)
}

func updateTodo(c echo.Context) error {
	id := c.Param("id")
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := "UPDATE todos SET title = ?, completed = ? WHERE id = ?"
	_, err := db.Exec(query, todo.Title, todo.Completed, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func deleteTodo(c echo.Context) error {
	id := c.Param("id")

	query := "DELETE FROM todos WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
