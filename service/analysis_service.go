package service

import (
	"fmt"

	"bridge-trajectory/analysis"
	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

type AnalysisService struct {
	store *store.Store
}

func NewAnalysisService(database *store.Store) *AnalysisService {
	return &AnalysisService{store: database}
}

func (s *AnalysisService) Analyze(id string) (analysis.Report, error) {
	record, err := s.store.GetTrajectory(id)
	if err != nil {
		return analysis.Report{}, err
	}
	return analysis.Analyze(record)
}

func (s *AnalysisService) AnalyzeAll(filter domain.TrajectoryFilter) ([]analysis.Report, error) {
	items, err := s.store.ListTrajectories(filter)
	if err != nil {
		return nil, err
	}
	reports := make([]analysis.Report, 0, len(items))
	for _, item := range items {
		report, reportErr := analysis.Analyze(item)
		if reportErr != nil {
			return nil, fmt.Errorf("analyze %s: %w", item.ID, reportErr)
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (s *AnalysisService) Compare(firstID, secondID string) (analysis.Comparison, error) {
	first, err := s.Analyze(firstID)
	if err != nil {
		return analysis.Comparison{}, err
	}
	second, err := s.Analyze(secondID)
	if err != nil {
		return analysis.Comparison{}, err
	}
	return analysis.Compare(first, second), nil
}

func (s *AnalysisService) Review(id string) (string, error) {
	report, err := s.Analyze(id)
	if err != nil {
		return "", err
	}
	return analysis.RiskBand(report.RiskScore), nil
}

func (s *AnalysisService) CompareLatest(filter domain.TrajectoryFilter) (analysis.Comparison, error) {
	reports, err := s.AnalyzeAll(filter)
	if err != nil {
		return analysis.Comparison{}, err
	}
	if len(reports) < 2 {
		return analysis.Comparison{}, fmt.Errorf("at least two trajectories are required")
	}
	return analysis.Compare(reports[len(reports)-2], reports[len(reports)-1]), nil
}
