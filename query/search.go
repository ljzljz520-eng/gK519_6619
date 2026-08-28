package query

import (
	"sort"
	"strings"

	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

type Catalog struct {
	store *store.Store
}

func NewCatalog(database *store.Store) *Catalog { return &Catalog{store: database} }

func (c *Catalog) SearchBridges(term string) ([]domain.BridgeRecord, error) {
	bridges, err := c.store.ListBridges()
	if err != nil {
		return nil, err
	}
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return bridges, nil
	}
	filtered := make([]domain.BridgeRecord, 0, len(bridges))
	for _, bridge := range bridges {
		if strings.Contains(strings.ToLower(bridge.ID), term) || strings.Contains(strings.ToLower(bridge.Name), term) {
			filtered = append(filtered, bridge)
		}
	}
	return filtered, nil
}

func (c *Catalog) SearchTrajectories(filter domain.TrajectoryFilter) ([]domain.TrajectoryRecord, error) {
	items, err := c.store.ListTrajectories(filter)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].CreatedUnix < items[j].CreatedUnix })
	return items, nil
}

func (c *Catalog) ScenariosFor(bridgeID string) ([]domain.WindScenarioRecord, error) {
	return c.store.ListScenarios(bridgeID)
}

func (c *Catalog) EventsFor(subject string) ([]domain.AuditEvent, error) {
	return c.store.ListEvents(subject)
}

func (c *Catalog) StoreEvent(event domain.AuditEvent) error {
	return c.store.AppendEvent(event)
}
