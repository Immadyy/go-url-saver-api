package store

import (
	"database/sql"
	"fmt"
	"url_saver/internal/models"

	"github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

func (p *PostgresStore) Save(data models.Link) (models.Link, error) {
	query := `
	INSERT INTO links (title, link, created_at)
	VALUES ($1, $2, DEFAULT)
	RETURNING id`

	err := p.DB.QueryRow(query, data.Title, data.Link).Scan(&data.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return models.Link{}, models.ErrDuplicateLink
			}
		}
		return models.Link{}, fmt.Errorf("Could'nt save data: %w", err)
	}

	return data, nil
}
