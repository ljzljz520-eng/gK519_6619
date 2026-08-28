package store

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"bridge-trajectory/domain"
	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("record not found")

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range allBuckets() {
			if _, createErr := tx.CreateBucketIfNotExists(bucket); createErr != nil {
				return createErr
			}
		}
		return nil
	}); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return &Store{db: db, path: path}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

func (s *Store) PutBridge(record domain.BridgeRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.put(bucketNames.Bridges, record.ID, record)
}

func (s *Store) GetBridge(id string) (domain.BridgeRecord, error) {
	var result domain.BridgeRecord
	err := s.get(bucketNames.Bridges, id, &result)
	return result, err
}

func (s *Store) ListBridges() ([]domain.BridgeRecord, error) {
	var result []domain.BridgeRecord
	err := s.list(bucketNames.Bridges, func(data []byte) error {
		var item domain.BridgeRecord
		if err := decode(data, &item); err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) PutScenario(record domain.WindScenarioRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.put(bucketNames.Scenarios, record.ID, record)
}

func (s *Store) GetScenario(id string) (domain.WindScenarioRecord, error) {
	var result domain.WindScenarioRecord
	err := s.get(bucketNames.Scenarios, id, &result)
	return result, err
}

func (s *Store) ListScenarios(bridgeID string) ([]domain.WindScenarioRecord, error) {
	var result []domain.WindScenarioRecord
	err := s.list(bucketNames.Scenarios, func(data []byte) error {
		var item domain.WindScenarioRecord
		if err := decode(data, &item); err != nil {
			return err
		}
		if bridgeID == "" || item.BridgeID == bridgeID {
			result = append(result, item)
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) PutTrajectory(record domain.TrajectoryRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.put(bucketNames.Trajectories, record.ID, record.Clone())
}

func (s *Store) GetTrajectory(id string) (domain.TrajectoryRecord, error) {
	var result domain.TrajectoryRecord
	err := s.get(bucketNames.Trajectories, id, &result)
	if err == nil {
		result = result.Clone()
	}
	return result, err
}

func (s *Store) ListTrajectories(filter domain.TrajectoryFilter) ([]domain.TrajectoryRecord, error) {
	var result []domain.TrajectoryRecord
	limit := domain.NormalizeLimit(filter.Limit)
	err := s.list(bucketNames.Trajectories, func(data []byte) error {
		var item domain.TrajectoryRecord
		if err := decode(data, &item); err != nil {
			return err
		}
		if !filter.Match(item) {
			return nil
		}
		result = append(result, item.Clone())
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedUnix > result[j].CreatedUnix })
	if len(result) > limit {
		result = result[:limit]
	}
	return result, err
}

func (s *Store) PutViewPreference(record domain.ViewPreferenceRecord) error {
	if record.UserID == "" {
		return fmt.Errorf("user id is required")
	}
	record.Mode = domain.DefaultView(record.Mode)
	return s.put(bucketNames.Views, record.UserID, record)
}

func (s *Store) GetViewPreference(userID string) (domain.ViewPreferenceRecord, error) {
	var result domain.ViewPreferenceRecord
	err := s.get(bucketNames.Views, userID, &result)
	return result, err
}

func (s *Store) AppendEvent(event domain.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	return s.put(bucketNames.Events, event.ID, event)
}

func (s *Store) ListEvents(subject string) ([]domain.AuditEvent, error) {
	var result []domain.AuditEvent
	err := s.list(bucketNames.Events, func(data []byte) error {
		var item domain.AuditEvent
		if err := decode(data, &item); err != nil {
			return err
		}
		if subject == "" || item.Subject == subject {
			result = append(result, item)
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].AtUnix < result[j].AtUnix })
	return result, err
}

func (s *Store) put(bucket []byte, key string, value any) error {
	if key == "" {
		return fmt.Errorf("record key is required")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(key), data) })
}

func (s *Store) get(bucket []byte, key string, target any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return ErrNotFound
		}
		return decode(cloneBytes(value), target)
	})
}

func (s *Store) list(bucket []byte, visit func([]byte) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			return visit(cloneBytes(value))
		})
	})
}

func (s *Store) Exists() bool {
	if s == nil {
		return false
	}
	_, err := os.Stat(s.Path())
	return err == nil
}
