package repo

import "github.com/jmoiron/sqlx"

type Repo struct {
	db *sqlx.DB
}

func NewNachinka(db *sqlx.DB) Repo {
	return Repo{db: db}
}
