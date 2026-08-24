package service

import (
	"fmt"
	"time"

	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

type BridgeService struct {
	store *store.Store
	now   func() time.Time
}

func NewBridgeService(database *store.Store) *BridgeService {
	return &BridgeService{store: database, now: func() time.Time { return time.Unix(0, 0) }}
}

func (s *BridgeService) RegisterBridge(id, name string, length, deckHeight float64) (domain.BridgeRecord, error) {
	bridge := domain.NewBridge(id, name, length, deckHeight, s.now())
	if err := bridge.Validate(); err != nil {
		return domain.BridgeRecord{}, err
	}
	if _, err := s.store.GetBridge(id); err == nil {
		return domain.BridgeRecord{}, fmt.Errorf("bridge %q already exists", id)
	}
	if err := s.store.PutBridge(bridge); err != nil {
		return domain.BridgeRecord{}, err
	}
	return bridge, nil
}

func (s *BridgeService) ListBridges() ([]domain.BridgeRecord, error) {
	return s.store.ListBridges()
}

func (s *BridgeService) SaveScenario(id, bridgeID string, wind, amplitude, step, duration float64, description string) (domain.WindScenarioRecord, error) {
	if _, err := s.store.GetBridge(bridgeID); err != nil {
		return domain.WindScenarioRecord{}, fmt.Errorf("bridge %q does not exist", bridgeID)
	}
	scenario := domain.NewScenario(id, bridgeID, wind, amplitude, step, duration, description)
	if err := scenario.Validate(); err != nil {
		return domain.WindScenarioRecord{}, err
	}
	if err := s.store.PutScenario(scenario); err != nil {
		return domain.WindScenarioRecord{}, err
	}
	return scenario, nil
}

func (s *BridgeService) ListScenarios(bridgeID string) ([]domain.WindScenarioRecord, error) {
	return s.store.ListScenarios(bridgeID)
}

func (s *BridgeService) SetView(userID string, mode domain.ViewMode) (domain.ViewPreferenceRecord, error) {
	if userID == "" {
		return domain.ViewPreferenceRecord{}, fmt.Errorf("user id is required")
	}
	preference := domain.ViewPreferenceRecord{UserID: userID, Mode: domain.DefaultView(mode), Updated: s.now().Unix()}
	if err := s.store.PutViewPreference(preference); err != nil {
		return domain.ViewPreferenceRecord{}, err
	}
	return preference, nil
}

func (s *BridgeService) GetView(userID string) (domain.ViewMode, error) {
	preference, err := s.store.GetViewPreference(userID)
	if err != nil {
		if err == store.ErrNotFound {
			return domain.ViewTop, nil
		}
		return domain.ViewTop, err
	}
	return domain.DefaultView(preference.Mode), nil
}
