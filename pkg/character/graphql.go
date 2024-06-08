package character

import (
	"fmt"
	"slices"

	"github.com/graphql-go/graphql"
)

var localcharacter []*CharacterData = createCharacterArrays()

var characterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Character",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"user_id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"total_focused_time": &graphql.Field{
			Type: graphql.Float,
		},
		"custom_metrics": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "CustomMetrics",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.String,
					},
					"character_id": &graphql.Field{
						Type: graphql.String,
					},
					"type": &graphql.Field{
						Type: graphql.String,
					},
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"value": &graphql.Field{
						Type: graphql.String,
					},
				},
			})),
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"character": &graphql.Field{
			Type:        characterType,
			Description: "Get a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				i := slices.IndexFunc(localcharacter, func(c *CharacterData) bool {
					return c.ID == id
				})

				if i == -1 {
					return nil, fmt.Errorf("character not found")
				}

				return localcharacter[i], nil
			},
		},
		"characters": &graphql.Field{
			Type:        graphql.NewList(characterType),
			Description: "Get all characters",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return localcharacter, nil
			},
		},
		// "testRequiredAuth": &graphql.Field{
		// 	Type:        graphql.String,
		// 	Description: "Give me your jwt, I'll greet you",
		// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// 		authProfile, err := auth.GetProfileByContext(params.Context)
		// 		if err != nil {
		// 			return nil, err
		// 		}

		// 		return "Hello from " + authProfile.Email, nil
		// 	},
		// },
	},
})

var newCustomMetricInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CusTomMetricInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"character_id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"value": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

func updatedCharacter(id string, user_id string, name string, tags []string, totalFocusTime float64, customMetricsDatas []CustomMetricData) CharacterData {
	return CharacterData{
		ID:               id,
		UserID:           user_id,
		Name:             name,
		Tags:             tags,
		TotalFocusedTime: totalFocusTime,
		CustomMetrics:    customMetricsDatas,
	}
}

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"updateCharacter": &graphql.Field{
			Type:        characterType,
			Description: "Update a character by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"user_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"total_focus_time": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"custome_metrics": &graphql.ArgumentConfig{
					Type: graphql.NewList(newCustomMetricInput),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				userID := params.Args["user_id"].(string)
				name := params.Args["name"].(string)
				tags := params.Args["tags"].([]string)
				totalFocusTime := params.Args["total_focus_time"].(float64)
				customMetricsInput := params.Args["custom_metrics"].([]interface{})

				var customMetricDatas []CustomMetricData
				for _, cm := range customMetricsInput {
					cmMap, _ := cm.(map[string]interface{})
					metricID, _ := cmMap["id"].(string)
					metricCharacterID, _ := cmMap["character_id"].(string)
					metricType, _ := cmMap["type"].(string)
					metricName, _ := cmMap["name"].(string)
					metricValue, _ := cmMap["value"].(string)

					customMetric := CustomMetricData{
						ID:          metricID,
						CharacterID: metricCharacterID,
						Type:        metricType,
						Name:        metricName,
						Value:       metricValue,
					}
					customMetricDatas = append(customMetricDatas, customMetric)
				}

				return updatedCharacter(id, userID, name, tags, totalFocusTime, customMetricDatas), nil
			},
		},
	},
})

var CharacterSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
