package main

import (
	"fmt"

	"github.com/Stiffjobs/toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	m.data[distance.OBUID] += distance.Value
	return nil
}

func (m *MemoryStore) Get(obuID int) (float64, error) {
	dist, ok := m.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance in the storage: %d", obuID)
	}
	return dist, nil
}
