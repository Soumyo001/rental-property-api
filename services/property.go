package services

import (
	"encoding/json"
	"fmt"
	"os"

	"rental-property-api/models"
)

type Store struct {
	properties []models.SourceProperty
	index      map[string]int
}

var defaultStore *Store

func (s *Store) Count() int {
	return len(s.properties)
}

func newStore(records []models.SourceProperty) *Store {
	index := make(map[string]int, len(records))
	for i, record := range records {
		index[record.ID] = i
	}

	return &Store{
		properties: records,
		index:      index,
	}
}

func loadFromFile(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var decoded []models.SourceProperty
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, fmt.Errorf("failed to decode file %s: %w", path, err)
	}

	return newStore(decoded), nil
}

func Init(path string) error {
	store, err := loadFromFile(path)
	if err != nil {
		return err
	}

	defaultStore = store
	return nil
}

func GetStore() *Store {
	return defaultStore
}
