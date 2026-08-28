package service

import (
	"sync"

	"bridge-trajectory/domain"
)

type TrajectoryCache struct {
	mu      sync.Mutex
	latest  []domain.TrajectoryPoint
	records map[string]domain.TrajectoryRecord
}

func NewTrajectoryCache() *TrajectoryCache {
	return &TrajectoryCache{records: make(map[string]domain.TrajectoryRecord)}
}

func (c *TrajectoryCache) Acquire(size int) []domain.TrajectoryPoint {
	c.mu.Lock()
	defer c.mu.Unlock()
	if size > cap(c.latest) {
		c.latest = make([]domain.TrajectoryPoint, 0, size)
	}
	return c.latest[:0]
}

func (c *TrajectoryCache) Remember(record domain.TrajectoryRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latest = record.Points
	c.records[record.ID] = record
}

func (c *TrajectoryCache) Get(id string) (domain.TrajectoryRecord, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record, ok := c.records[id]
	return record, ok
}

func (c *TrajectoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latest = nil
	c.records = make(map[string]domain.TrajectoryRecord)
}

func (c *TrajectoryCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.records)
}
