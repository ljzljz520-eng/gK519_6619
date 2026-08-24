package store

import (
	"fmt"
	"sort"

	"bridge-trajectory/domain"
	"go.etcd.io/bbolt"
)

type StoreCounts struct {
	Bridges      int `json:"bridges"`
	Scenarios    int `json:"scenarios"`
	Trajectories int `json:"trajectories"`
	Views        int `json:"views"`
	Events       int `json:"events"`
}

func (s *Store) Counts() (StoreCounts, error) {
	counts := StoreCounts{}
	for _, entry := range []struct {
		bucket []byte
		set    func(int)
	}{
		{bucketNames.Bridges, func(value int) { counts.Bridges = value }},
		{bucketNames.Scenarios, func(value int) { counts.Scenarios = value }},
		{bucketNames.Trajectories, func(value int) { counts.Trajectories = value }},
		{bucketNames.Views, func(value int) { counts.Views = value }},
		{bucketNames.Events, func(value int) { counts.Events = value }},
	} {
		value, err := s.countBucket(entry.bucket)
		if err != nil {
			return StoreCounts{}, err
		}
		entry.set(value)
	}
	return counts, nil
}

func (s *Store) countBucket(bucket []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, fmt.Errorf("database is closed")
	}
	count := 0
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(key, value []byte) error {
			if value != nil {
				count++
			}
			return nil
		})
	})
	return count, err
}

func (s *Store) DeleteTrajectory(id string) error {
	if id == "" {
		return fmt.Errorf("trajectory id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketNames.Trajectories)
		if bucket.Get([]byte(id)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(id))
	})
}

func (s *Store) DeleteScenario(id string) error {
	if id == "" {
		return fmt.Errorf("scenario id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketNames.Scenarios).Delete([]byte(id)) })
}

func (s *Store) TrajectoryIDs(filter domain.TrajectoryFilter) ([]string, error) {
	items, err := s.ListTrajectories(filter)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	return ids, nil
}

func (s *Store) VerifyBuckets() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("database is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		for _, bucket := range allBuckets() {
			if tx.Bucket(bucket) == nil {
				return fmt.Errorf("missing bucket %s", bucket)
			}
		}
		return nil
	})
}
