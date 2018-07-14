package main

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	"gopkg.in/inconshreveable/log15.v2"
)

var statements = map[string]string{
	"check_if_user_exists": "SELECT exists (SELECT 1 FROM users where username = $1)",
	"get_users": `
		SELECT array_to_json(COALESCE(array_agg(t), '{}')) FROM (
			SELECT *, (SELECT up FROM (SELECT firstname, lastname FROM user_profile WHERE user_id=u.id) up) AS profile
			FROM users AS u
			ORDER BY u.id
		) AS t`,
	"create_user":         "INSERT INTO users (username, score) VALUES ($1, $2) RETURNING id, created",
	"create_user_profile": "INSERT INTO user_profile (user_id, firstname, lastname) VALUES ($1, $2, $3)",
	"update_user_score":   "UPDATE users SET score = $2, updated = $3 WHERE id = $1",
	"log_action":          "INSERT INTO actions (description) VALUES ($1)",
	"get_actions": `
		SELECT array_to_json(COALESCE(array_agg(a), '{}')) FROM (
			SELECT id, description, created
			FROM actions
			ORDER BY id
		) AS a
	`,
}

var maxDBConnections = 3

func dbConnect() (*pgx.ConnPool, error) {
	config := pgx.ConnConfig{
		Host:     "localhost",
		Database: "db",
		User:     "user",
		Password: "pass",
	}

	if Debug {
		config.Logger = log15adapter.NewLogger(log15.New("module", "pgx"))
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: maxDBConnections,
	}

	return pgx.NewConnPool(poolConfig)
}

func prepareStatements() error {
	for name, sql := range statements {
		_, err := db.Prepare(name, sql)
		if err != nil {
			return err
		}
	}

	return nil
}
