package graph

import (
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"
	"tenkhours/services/analytics/graph/model"
)

func MapToSnapshotDto(snapshot *analyticsdb.Snapshot) *model.Snapshot {
	character := snapshot.Character
	customMetrics := make([]model.CustomMetricData, 0)
	for _, customMetric := range character.CustomMetrics {
		customMetrics = append(customMetrics, *MapToCustomMetricDto(&customMetric))
	}

	return &model.Snapshot{
		ID:        snapshot.ID,
		Timestamp: snapshot.Timestamp,
		Character: &model.CharacterData{
			ID:               character.ID,
			ProfileID:        character.ProfileID,
			Name:             character.Name,
			Gender:           character.Gender,
			Tags:             character.Tags,
			TotalFocusedTime: int(character.TotalFocusedTime),
			CustomMetrics:    customMetrics,
		},
	}
}

func MapToCustomMetricDto(customMetrics *coredb.CustomMetric) *model.CustomMetricData {
	props := make([]model.MetricPropertyData, 0)
	for _, prop := range customMetrics.Properties {
		props = append(props, *MapToMetricPropertyDto(&prop))
	}

	return &model.CustomMetricData{
		ID:          customMetrics.ID,
		Name:        customMetrics.Name,
		Description: &customMetrics.Description,
		Time:        int(customMetrics.Time),
		Style: &model.MetricStyleData{
			Color: &customMetrics.Style.Color,
			Icon:  &customMetrics.Style.Icon,
		},
		Properties: props,
	}
}

func MapToMetricPropertyDto(properties *coredb.MetricProperty) *model.MetricPropertyData {
	return &model.MetricPropertyData{
		ID:    properties.ID,
		Name:  properties.Name,
		Type:  properties.Type,
		Value: properties.Value,
		Unit:  &properties.Unit,
	}
}
