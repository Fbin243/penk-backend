package character

import (
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db Database

func init() {
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	} else {
		log.Println("Connected to MongoDB!")
	}
}

var characterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Character",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
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
				character, err := db.GetCharacterByID(id)
				if err != nil {
					return nil, fmt.Errorf("character not found: %v", err)
				}

				return character, nil
			},
		},
		"characters": &graphql.Field{
			Type:        graphql.NewList(characterType),
			Description: "Get all characters",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				characters, err := db.GetAllCharacters()
				if err != nil {
					return nil, fmt.Errorf("error fetching characters: %v", err)
				}
				return characters, nil
			},
		},
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

func updatedCharacter(user_id string, name string, tags []string, totalFocusTime float64, customMetricsDatas []CustomMetricData) CharacterData {
	return CharacterData{
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
		"createCharacter": &graphql.Field{
			Type:        characterType,
			Description: "Create a character",
			Args: graphql.FieldConfigArgument{
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
				"custom_metrics": &graphql.ArgumentConfig{
					Type: graphql.NewList(newCustomMetricInput),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userID := params.Args["userID"].(string)
				name := params.Args["name"].(string)
				var tags []string
				tagsInterface := params.Args["tags"].([]interface{})
				tags = make([]string, len(tagsInterface))
				for i, tag := range tagsInterface {
					tags[i] = tag.(string)
				}
				totalFocusTime := params.Args["total_focus_time"].(float64)
				customMetricsInput := params.Args["custom_metrics"].([]interface{})

				var customMetricDatas []CustomMetricData
				for _, cm := range customMetricsInput {
					cmMap, _ := cm.(map[string]interface{})
					metricID, _ := cmMap["id"].(primitive.ObjectID)
					metricCharacterID, _ := cmMap["character_id"].(primitive.ObjectID)
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

				character := updatedCharacter(userID, name, tags, totalFocusTime, customMetricDatas)
				err := db.InsertCharacter(character)
				if err != nil {
					return nil, fmt.Errorf("failed to update character: %v", err)
				}

				return character, nil
			},
		},
		"updateCharacter": &graphql.Field{
			Type:        characterType,
			Description: "Update a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"user_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"total_focus_time": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"custom_metrics": &graphql.ArgumentConfig{
					Type: graphql.NewList(newCustomMetricInput),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				character, err := db.GetCharacterByID(id)
				if err != nil {
					return nil, fmt.Errorf("character not found: %v", err)
				}

				if userID, ok := params.Args["user_id"].(string); ok {
					character.UserID = userID
				}
				if name, ok := params.Args["name"].(string); ok {
					character.Name = name
				}
				if tagsInterface, ok := params.Args["tags"].([]interface{}); ok {
					tags := make([]string, len(tagsInterface))
					for i, tag := range tagsInterface {
						tags[i] = tag.(string)
					}
					character.Tags = tags
				}
				if totalFocusTime, ok := params.Args["total_focus_time"].(float64); ok {
					character.TotalFocusedTime = totalFocusTime
				}
				if customMetricsInput, ok := params.Args["custom_metrics"].([]interface{}); ok {
					var customMetricDatas []CustomMetricData
					for _, cm := range customMetricsInput {
						cmMap, _ := cm.(map[string]interface{})
						metricID, _ := cmMap["id"].(primitive.ObjectID)
						metricCharacterID, _ := cmMap["character_id"].(primitive.ObjectID)
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
					character.CustomMetrics = customMetricDatas
				}

				err = db.UpdateCharacter(id, *character)
				if err != nil {
					return nil, fmt.Errorf("failed to update character: %v", err)
				}

				return character, nil
			},
		},
		"deleteCharacter": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				err := db.DeleteCharacter(id)
				if err != nil {
					return nil, fmt.Errorf("failed to delete character: %v", err)
				}

				return true, nil
			},
		},
	},
})

var CharacterSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
