package memorystorage

import (
	"context"
	"sync"

	"github.com/krysanov-as/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; ok {
		return storage.ErrDateBusy
	}

	s.events[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; !ok {
		return storage.ErrNotFound
	}

	s.events[e.ID] = e
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0, len(s.events))
	for _, e := range s.events {
		result = append(result, e)
	}

	return result, nil
}
