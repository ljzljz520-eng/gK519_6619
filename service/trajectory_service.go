package service

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"bridge-trajectory/calc"
	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

type TrajectoryService struct {
	store    *store.Store
	sample   calc.Sampler
	cache    *TrajectoryCache
	now      func() time.Time
	sequence uint64
}

func NewTrajectoryService(database *store.Store) *TrajectoryService {
	return &TrajectoryService{store: database, sample: calc.NewSampler(), cache: NewTrajectoryCache(), now: func() time.Time { return time.Unix(0, 0) }}
}

func (s *TrajectoryService) Calculate(req domain.CalculationRequest) (domain.TrajectoryRecord, error) {
	if err := domain.ValidateCalculation(req); err != nil {
		return domain.TrajectoryRecord{}, err
	}
	bridge, err := s.store.GetBridge(req.BridgeID)
	if err != nil {
		return domain.TrajectoryRecord{}, fmt.Errorf("load bridge: %w", err)
	}
	scenario, err := s.store.GetScenario(req.ScenarioID)
	if err != nil {
		return domain.TrajectoryRecord{}, fmt.Errorf("load scenario: %w", err)
	}
	if scenario.BridgeID != bridge.ID {
		return domain.TrajectoryRecord{}, fmt.Errorf("scenario belongs to another bridge")
	}
	model := calc.NewWindModel(req.WindSpeed, req.Amplitude)
	count := int(math.Floor(req.Duration/req.Step)) + 2
	points := s.cache.Acquire(count)
	points, err = s.sample.SampleInto(points, bridge.Length, bridge.DeckHeight, model, req.Step, req.Duration)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	sequence := atomic.AddUint64(&s.sequence, 1)
	record := domain.TrajectoryRecord{ID: fmt.Sprintf("%s-%d", req.ScenarioID, sequence), BridgeID: bridge.ID, ScenarioID: scenario.ID, Points: points, View: req.View, Status: "ready", CreatedUnix: s.now().Unix()}
	if err := s.store.PutTrajectory(record); err != nil {
		return domain.TrajectoryRecord{}, err
	}
	s.cache.Remember(record)
	return record, nil
}

func (s *TrajectoryService) Get(id string) (domain.TrajectoryRecord, error) {
	if record, ok := s.cache.Get(id); ok {
		return record, nil
	}
	return s.store.GetTrajectory(id)
}

func (s *TrajectoryService) List(filter domain.TrajectoryFilter) ([]domain.TrajectoryRecord, error) {
	return s.store.ListTrajectories(filter)
}

func (s *TrajectoryService) Archive(id string) error {
	record, err := s.store.GetTrajectory(id)
	if err != nil {
		return err
	}
	record.Status = "archived"
	return s.store.PutTrajectory(record)
}

func (s *TrajectoryService) CacheSize() int { return s.cache.Size() }

func (s *TrajectoryService) ClearCache() { s.cache.Clear() }

func (s *TrajectoryService) ChangeView(id string, mode domain.ViewMode) (domain.TrajectoryRecord, error) {
	record, err := s.Get(id)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	updated, err := record.WithView(mode)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	if err := s.store.PutTrajectory(updated); err != nil {
		return domain.TrajectoryRecord{}, err
	}
	s.cache.Remember(updated)
	return updated, nil
}

func (s *TrajectoryService) Remove(id string) error {
	if _, err := s.store.GetTrajectory(id); err != nil {
		return err
	}
	s.cache.Clear()
	return s.store.DeleteTrajectory(id)
}
