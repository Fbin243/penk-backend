package graph

import (
	"tenkhours/services/analytics/graph/model"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/repo"
)

func MapToSnapshotDto(snapshot *analyticsRepo.Snapshot) *model.Snapshot {
	character := snapshot.Character
	customMetrics := make([]model.SnapshotCustomMetric, 0)
	for _, customMetric := range character.CustomMetrics {
		customMetrics = append(customMetrics, *MapToCustomMetricDto(&customMetric))
	}

	return &model.Snapshot{
		ID:        snapshot.ID,
		Timestamp: snapshot.Timestamp,
		Character: model.SnapshotCharacter{
			ID:               character.ID,
			ProfileID:        character.ProfileID,
			Name:             character.Name,
			Gender:           character.Gender,
			Tags:             character.Tags,
			TotalFocusedTime: character.TotalFocusedTime,
			CustomMetrics:    customMetrics,
		},
		Description: snapshot.Description,
	}
}

func MapToCustomMetricDto(customMetrics *repo.CustomMetric) *model.SnapshotCustomMetric {
	props := make([]model.SnapshotMetricProperty, 0)
	for _, prop := range customMetrics.Properties {
		props = append(props, *MapToMetricPropertyDto(&prop))
	}

	return &model.SnapshotCustomMetric{
		ID:          customMetrics.ID,
		Name:        customMetrics.Name,
		Description: &customMetrics.Description,
		Time:        customMetrics.Time,
		Style: model.SnapshotMetricStyle{
			Color: &customMetrics.Style.Color,
			Icon:  &customMetrics.Style.Icon,
		},
		Properties: props,
	}
}

func MapToMetricPropertyDto(properties *repo.MetricProperty) *model.SnapshotMetricProperty {
	return &model.SnapshotMetricProperty{
		ID:    properties.ID,
		Name:  properties.Name,
		Type:  model.MetricPropertyType(properties.Type.String()),
		Value: properties.Value,
		Unit:  &properties.Unit,
	}
}
