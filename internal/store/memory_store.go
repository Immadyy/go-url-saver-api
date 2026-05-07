package store

import (
	"fmt"
	"sync"
	"url_saver/internal/models"
)

type MemoryStore struct {
	mu    sync.Mutex
	Links []models.Link
	ID    int64
}

func (d *MemoryStore) Save(data models.Link) (models.Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.ID++
	data.ID = d.ID
	d.Links = append(d.Links, data)
	return data, nil
}

func (d *MemoryStore) GetAll() ([]models.Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	data := make([]models.Link, len(d.Links))
	copy(data, d.Links)
	return data, nil
}

func (d *MemoryStore) Update(updId int64, data models.Link) (models.Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i := range d.Links {
		if d.Links[i].ID == updId {
			d.Links[i].Title = data.Title
			d.Links[i].Link = data.Link
			return d.Links[i], nil
		}
	}
	return models.Link{}, fmt.Errorf("data not found")
}

func (d *MemoryStore) Delete(delId int64) (models.Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range d.Links {
		if d.Links[i].ID == delId {
			data := d.Links[i]
			d.Links = append(d.Links[:i], d.Links[i+1:]...)
			return data, nil
		}
	}
	return models.Link{}, fmt.Errorf("data not found")
}
