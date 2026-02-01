package state

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
	SeenURLs map[string]bool `json:"seen_urls"`
	mu       sync.RWMutex    `json:"-"`
	FilePath string          `json:"-"`
}

func NewStore(filename string) (*Store, error) {
	s := &Store{
		SeenURLs: make(map[string]bool),
		FilePath: filename,
	}

	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Store) HasSeen(url string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SeenURLs[url]
}

func (s *Store) MarkSeen(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SeenURLs[url] = true
	return s.save()
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0644)
}
