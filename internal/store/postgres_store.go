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
	INSERT INTO links (title, link)
	VALUES ($1, $2)
	RETURNING id`

	err := p.DB.QueryRow(query, data.Title, data.Link).Scan(&data.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return models.Link{}, models.ErrDuplicateLink
			}
		}
		return models.Link{}, fmt.Errorf("Couldn't save data: %w", err)
	}

	return data, nil
}

func (p *PostgresStore) GetAll() ([]models.Link, error) {
	query := `SELECT id, title, link FROM links`

	rows, err := p.DB.Query(query)
	if err != nil {
		return []models.Link{}, fmt.Errorf("something went wrong: %w", err)
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		err := rows.Scan(&l.ID, &l.Title, &l.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		links = append(links, l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return links, nil
}
