package service

import (
	"errors"
	"testing"
	"url_saver/internal/models"

	"github.com/stretchr/testify/assert"
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

// 2. Tests for CreateLink & ValidateLink
func TestCreateLink(t *testing.T) {
	mockStore := &MockStore{}
	service := NewLinkService(mockStore)

	type testCase struct {
		name         string
		input        models.Link
		expectedLink string
		wantErr      bool
	}

	tests := []testCase{
		{
			name:         "Success - clean link remains unchanged",
			input:        models.Link{Title: "Google", Link: "https://google.com"},
			expectedLink: "https://google.com",
			wantErr:      false,
		},
		{
			name:         "Success - prepends https prefix if missing",
			input:        models.Link{Title: "Google", Link: "google.com"},
			expectedLink: "https://google.com",
			wantErr:      false,
		},
		{
			name:         "Error - empty title",
			input:        models.Link{Title: "", Link: "google.com"},
			expectedLink: "",
			wantErr:      true,
		},
		{
			name:         "Error - empty link",
			input:        models.Link{Title: "Google", Link: ""},
			expectedLink: "",
			wantErr:      true,
		},
		{
			name:         "Error - bad URL format string",
			input:        models.Link{Title: "Google", Link: "not a valid url !!"},
			expectedLink: "",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := service.CreateLink(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLink, actual.Link)
			}
		})
	}
}

// 3. Simple Tests for remaining CRUD Methods
func TestGetAllLinks(t *testing.T) {
	mockStore := &MockStore{}
	service := NewLinkService(mockStore)

	links, err := service.GetAllLinks()
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "Google", links[0].Title)
}

func TestUpdateLink(t *testing.T) {
	mockStore := &MockStore{}
	service := NewLinkService(mockStore)

	// Valid update
	input := models.Link{Title: "New Title", Link: "newurl.com"}
	updatedLink, err := service.UpdateLink(5, input)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), updatedLink.ID)
	assert.Equal(t, "https://newurl.com", updatedLink.Link)

	// Invalid update (fails validation)
	invalidInput := models.Link{Title: "", Link: "newurl.com"}
	_, err = service.UpdateLink(5, invalidInput)
	assert.Error(t, err)
}

func TestDeleteLink(t *testing.T) {
	mockStore := &MockStore{}
	service := NewLinkService(mockStore)

	err := service.DeleteLink(10)
	assert.NoError(t, err)
}
