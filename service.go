package main

import (
	"time"
)

type User struct {
	ID       int        `json:"id"`
	Username string     `json:"username"`
	Score    float64    `json:"score"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated"`
	Profile  struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	} `json:"profile"`
}

func userExists(user *User) (exists bool, err error) {
	row := db.QueryRow("check_if_user_exists", user.Username)
	err = row.Scan(&exists)
	return
}

func getUsers() (users []User, err error) {
	rows, err := db.Query("get_users")
	if err != nil {
		return
	}

	rows.Next()
	err = rows.Scan(&users)
	if err != nil {
		return
	}

	rows.Close()
	return
}

func userScoreFormula(score, usersCount float64) float64 {
	return 1 + score + score/usersCount
}

type Action struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

func getActions() (actions []Action, err error) {
	rows, err := db.Query("get_actions")
	if err != nil {
		return
	}

	rows.Next()
	err = rows.Scan(&actions)
	if err != nil {
		return
	}

	rows.Close()
	return
}
