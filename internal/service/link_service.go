package service

import (
	"fmt"
	"net/url"
	"strings"
	"url_saver/internal/models"
)

type LinkService struct {
	Store LinkStore
}

type LinkStore interface {
	Save(data models.Link) (models.Link, error)
	//GetAll() ([]models.Link, error)
	//Update(upId int64, data models.Link) (models.Link, error)
	//Delete(delId int64) (models.Link, error)
}

func NewLinkService(l LinkStore) *LinkService {
	return &LinkService{
		Store: l,
	}
}

func (l *LinkService) ValidateLink(data models.Link) (models.Link, error) {
	if data.Link == "" || data.Title == "" {
		return models.Link{}, fmt.Errorf("Title and link cannot be empty.")
	}

	if !strings.HasPrefix(data.Link, "http://") && !strings.HasPrefix(data.Link, "https://") {
		data.Link = "https://" + data.Link
	}

	if _, err := url.ParseRequestURI(data.Link); err != nil {
		return models.Link{}, fmt.Errorf("Bad URL format")
	}

	return data, nil
}

func (l *LinkService) CreateLink(data models.Link) (models.Link, error) {
	Data, err := l.ValidateLink(data)
	if err != nil {
		return models.Link{}, err
	}
	return l.Store.Save(Data)
}

// func (l *LinkService) GetAllLinks() ([]models.Link, error) {
// 	data, err := l.Store.GetAll()
// 	return data, err
// }

// func (l *LinkService) UpdateLink(updId int64, data models.Link) (models.Link, error) {
// 	data, err := l.ValidateLink(data)
// 	if err != nil {
// 		return models.Link{}, err
// 	}

// 	link, err := l.Store.Update(updId, data)
// 	return link, err
// }

// func (l *LinkService) DeleteLink(delId int64) (models.Link, error) {
// 	data, err := l.Store.Delete(delId)
// 	return data, err
// }
