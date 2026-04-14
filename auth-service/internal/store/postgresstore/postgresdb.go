package postgresstore

import "database/sql"

type PostgresDB struct {
	DB *sql.DB
}

func New(databaseURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{
		DB: db,
	}, nil
}
