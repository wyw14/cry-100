package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wyw14/cry-100/internal/model"
)

type Store struct {
	mu     sync.Mutex
	dir    string
	events []model.Event
}

func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	store := &Store{dir: dir}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}
func (s *Store) load() error {
	file, err := os.Open(filepath.Join(s.dir, "events.jsonl"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event model.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return err
		}
		s.events = append(s.events, event)
	}
	return scanner.Err()
}
func (s *Store) Append(event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.OpenFile(filepath.Join(s.dir, "events.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err = file.Write(append(data, '\n')); err != nil {
		return err
	}
	s.events = append(s.events, event)
	return file.Sync()
}
func (s *Store) SnapshotPath() string { return filepath.Join(s.dir, "snapshot.json") }
func (s *Store) SaveSnapshot(snapshot model.Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	temp := s.SnapshotPath() + ".tmp"
	if err := os.WriteFile(temp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(temp, s.SnapshotPath())
}
func (s *Store) LoadSnapshot() (model.Snapshot, error) {
	data, err := os.ReadFile(s.SnapshotPath())
	if err != nil {
		return model.Snapshot{}, err
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return snapshot, fmt.Errorf("decode snapshot: %w", err)
	}
	return snapshot, nil
}

func (s *Store) RecordOperation(record model.OperationRecord) error {
	return s.AppendOperation(record)
}
