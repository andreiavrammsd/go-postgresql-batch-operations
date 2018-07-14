package main

import (
	"context"
	"net/http"
	"time"

	"fmt"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var (
	db    *pgx.ConnPool
	Debug = true
)

func main() {
	var err error
	db, err = dbConnect()
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}

	err = prepareStatements()
	if err != nil {
		log.Fatalf("error preparing statements: %s", err)
	}

	// HTTP server
	e := echo.New()
	e.Debug = Debug
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// HTTP handlers
	ug := e.Group("/users")
	ug.POST("", createUserHandler)
	ug.GET("", listUsersHandler)

	ag := e.Group("/actions")
	ag.GET("", listActionsHandler)

	// Start server and log errors
	e.Logger.Fatal(e.Start(":8608"))
}

func createUserHandler(c echo.Context) (err error) {
	// Get POST data
	user := &User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// Check if user exists
	exists, err := userExists(user)
	if err != nil {
		return
	}

	if exists {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{
			"User exists",
		})
	}

	// Get a context with timeout
	// For a large number of queries in batch it's safer to use a timeout context
	// https://github.com/jackc/pgx/blob/39bbc98d99d7b666759f84514859becf8067128f/batch.go#L63
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	// Start a transaction
	tx, err := db.BeginEx(ctx, nil)
	if err != nil {
		return
	}

	// Create the main user record
	args := []interface{}{
		user.Username,
		user.Score,
	}
	row := tx.QueryRow("create_user", args...)

	err = row.Scan(&user.ID, &user.Created)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			c.Logger().Error(e)
		}
		return
	}

	// Create the user profile record
	args = []interface{}{
		user.ID,
		user.Profile.FirstName,
		user.Profile.LastName,
	}
	_, err = tx.Exec("create_user_profile", args...)
	if err != nil {
		return
	}

	// Get all users to update scores
	users, err := getUsers()
	if err != nil {
		if e := tx.Rollback(); e != nil {
			c.Logger().Error(e)
		}
		return
	}

	// New batch
	b := tx.BeginBatch()

	usersCount := float64(len(users))
	for _, u := range users {
		u.Score = userScoreFormula(u.Score, usersCount)
		args = []interface{}{
			u.ID,
			u.Score,
			time.Now(),
		}

		// Add user update score query to batch
		b.Queue("update_user_score", args, nil, nil)
	}

	// Send the queries to db
	err = b.Send(ctx, nil)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			c.Logger().Error(e)
		}

		// It's very important to close the batch operation on error
		if e := b.Close(); e != nil {
			c.Logger().Error(e)
		}
		return
	}

	// Close batch operation
	err = b.Close()
	if err != nil {
		if e := tx.Rollback(); e != nil {
			c.Logger().Error(e)
		}
		return
	}

	// Create action record
	args = []interface{}{
		fmt.Sprintf("new user created with id %d and username %s", user.ID, user.Username),
	}
	_, err = tx.Exec("log_action", args...)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			c.Logger().Error(e)
		}
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func listUsersHandler(c echo.Context) error {
	users, err := getUsers()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func listActionsHandler(c echo.Context) error {
	actions, err := getActions()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, actions)
}
