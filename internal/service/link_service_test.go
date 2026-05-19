package service

import (
	"errors"
	"url_saver/internal/models"
)

// 1. Simple Mock Implementation
type MockStore struct {
	shouldFail bool // Toggle this if you want to mock database errors
}

func (m *MockStore) Save(data models.Link) (models.Link, error) {
	if m.shouldFail {
		return models.Link{}, errors.New("db error")
	}
	return data, nil
}

func (m *MockStore) GetAll() ([]models.Link, error) {
	return []models.Link{{ID: 1, Title: "Google", Link: "https://google.com"}}, nil
}

func (m *MockStore) Update(upId int64, data models.Link) (models.Link, error) {
	data.ID = upId
	return data, nil
}

func (m *MockStore) Delete(delId int64) error {
	return nil
}

//END
