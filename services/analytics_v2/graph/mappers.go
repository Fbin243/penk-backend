package graph

import (
	"time"

	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/services/analytics_v2/graph/model"
)

func mapModelToDto(snapshot *analyticsdb.Snapshot) *model.Snapshot {
	return &model.Snapshot{
		ID:        snapshot.ID.Hex(),
		Timestamp: snapshot.Timestamp.Format(time.RFC3339),
		Character: &model.Character{
			ID: snapshot.Character.ID.Hex(),
		},
	}
}
