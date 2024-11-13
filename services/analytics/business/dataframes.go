package business

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"tenkhours/pkg/utils"
	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/analytics/repo"

	"github.com/go-gota/gota/dataframe"
	"github.com/samber/lo"
)

type AnalyticsProcessor struct {
	AnalyticSections []model.AnalyticSection
	AnalyticResults  map[string]interface{}
	CapturedRecords  []model.CapturedRecord
	FilterType       FilterType
	StartTime        time.Time
	EndTime          time.Time
}

// ProcessCapturedRecords processes the captured records and returns the analytic results
func (ap *AnalyticsProcessor) ProcessCapturedRecords() map[string]interface{} {
	numberOfCapturedRecords := len(ap.CapturedRecords)
	dfCaptureRecordCustomMetricsData := make([]repo.DFCapturedRecordCustomMetric, 0)
	dfCaptureRecordsData := lo.Map(ap.CapturedRecords, func(record model.CapturedRecord, index int) repo.DFCapturedRecord {
		recordDay := record.Timestamp.Day()
		recordWeek := math.Min(math.Ceil(float64(recordDay)/NUMBER_OF_DAYS_IN_A_WEEK), NUMBER_OF_WEEKS_IN_A_MONTH)

		for _, metric := range record.CustomMetrics {
			dfCaptureRecordCustomMetricsData = append(dfCaptureRecordCustomMetricsData, repo.DFCapturedRecordCustomMetric{
				ID:       record.ID.Hex(),
				MetricID: metric.ID.Hex(),
				Time:     int(metric.Time),
			})
		}

		return repo.DFCapturedRecord{
			ID:               record.ID.Hex(),
			ProfileID:        record.Metadata.ProfileID.Hex(),
			CharacterID:      record.Metadata.CharacterID.Hex(),
			Year:             record.Timestamp.Year(),
			Month:            int(record.Timestamp.Month()),
			Week:             int(recordWeek),
			Day:              recordDay,
			TotalFocusedTime: int(record.TotalFocusedTime),
		}
	})

	dfCaptureRecords := dataframe.LoadStructs(dfCaptureRecordsData)
	dfCaptureRecordMetrics := dataframe.LoadStructs(dfCaptureRecordCustomMetricsData)
	dfInnerJoin := dfCaptureRecords.InnerJoin(dfCaptureRecordMetrics, "id")
	// Print the dfInnerJoin to debug if needed
	// jsonOutput, _ := json.MarshalIndent(dfInnerJoin.Records(), "", "  ")
	// fmt.Println(string(jsonOutput))

	// Process the captured records for the OVERALL section
	if numberOfCapturedRecords > 0 && lo.Contains(ap.AnalyticSections, model.AnalyticSectionOverall) {
		// Total active days
		totalFocusedDays := dfCaptureRecords.Nrow()

		// Total focused time
		totalFocusedTime := dfCaptureRecords.Col("total_focused_time").Sum()

		// Best & Current streak (count continuous days)
		timeStamps := lo.Map(ap.CapturedRecords, func(record model.CapturedRecord, index int) time.Time {
			return record.Timestamp
		})
		bestStreak := 1
		currentStreak := 1
		bestStreakStartDate := timeStamps[0]
		bestStreakEndDate := timeStamps[0]
		currentStreakStartDate := timeStamps[0]
		currentStreakEndDate := timeStamps[0]
		for i := 1; i < len(timeStamps); i++ {
			if utils.ResetTimeToBeginningOfDay(timeStamps[i]).Sub(utils.ResetTimeToBeginningOfDay(timeStamps[i-1])) == 24*time.Hour {
				currentStreak++
				currentStreakEndDate = timeStamps[i]
			} else {
				if currentStreak > bestStreak {
					bestStreak = currentStreak
					bestStreakStartDate = currentStreakStartDate
					bestStreakEndDate = currentStreakEndDate
				}
				currentStreak = 1
				currentStreakStartDate = timeStamps[i]
			}
		}

		if currentStreak > bestStreak {
			bestStreak = currentStreak
			bestStreakStartDate = currentStreakStartDate
			bestStreakEndDate = currentStreakEndDate
		}

		ap.AnalyticResults["overall"] = make(map[string]interface{})
		analyticResultsOverall := ap.AnalyticResults["overall"].(map[string]interface{})
		analyticResultsOverall["totalFocusedTime"] = totalFocusedTime
		analyticResultsOverall["totalFocusedDays"] = totalFocusedDays
		analyticResultsOverall["averageFocusedTime"] = totalFocusedTime / float64(totalFocusedDays)
		analyticResultsOverall["bestStreak"] = bestStreak
		analyticResultsOverall["currentStreak"] = currentStreak
		analyticResultsOverall["bestStreakStartDate"] = bestStreakStartDate
		analyticResultsOverall["bestStreakEndDate"] = bestStreakEndDate
		analyticResultsOverall["currentStreakStartDate"] = currentStreakStartDate
		analyticResultsOverall["currentStreakEndDate"] = currentStreakEndDate
	}

	// Process the captured records for the DISTRIBUTION section
	if numberOfCapturedRecords > 0 && lo.Contains(ap.AnalyticSections, model.AnalyticSectionDistribution) {
		// Total focused time
		totalFocusedTime := dfCaptureRecords.Col("total_focused_time").Sum()
		ap.AnalyticResults["distribution"] = make(map[string]interface{})
		if ap.FilterType == FilterTypeUser {
			distriRecords := dfCaptureRecords.GroupBy("profile_id", "character_id").Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"total_focused_time"})
			// Calculate the percentage of the total focused time
			for i := 0; i < distriRecords.Nrow(); i++ {
				characterID := distriRecords.Col("character_id").Elem(i).String()
				ap.AnalyticResults["distribution"].(map[string]interface{})[characterID] = distriRecords.Col("total_focused_time_SUM").Elem(i).Float()
			}
		} else {
			distriRecords := dfInnerJoin.GroupBy("character_id", "metric_id").Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"time"})
			nonMetricTime := totalFocusedTime
			// Calculate the percentage of the total focused time
			for i := 0; i < distriRecords.Nrow(); i++ {
				metricID := distriRecords.Col("metric_id").Elem(i).String()
				metricTime := distriRecords.Col("time_SUM").Elem(i).Float()
				ap.AnalyticResults["distribution"].(map[string]interface{})[metricID] = metricTime
				nonMetricTime -= metricTime
			}
			ap.AnalyticResults["distribution"].(map[string]interface{})["other"] = nonMetricTime
		}

		ap.AnalyticResults["distribution"].(map[string]interface{})["totalFocusedTime"] = totalFocusedTime
	}

	// Process the captured records for the TIMELINE section
	if numberOfCapturedRecords > 0 && lo.Contains(ap.AnalyticSections, model.AnalyticSectionTimeline) {
		dfa := &DataframeAggregator{}
		ap.AnalyticResults["timeline"] = make(map[string]interface{})
		analyticResultsTimeline := ap.AnalyticResults["timeline"].(map[string]interface{})
		dfa.SetAnalyticResults(analyticResultsTimeline)
		switch ap.FilterType {
		case FilterTypeUser:
			dfa.SetDataframe(dfCaptureRecords).SetSumField("total_focused_time").
				SetGroupByArray([]string{"profile_id", "character_id", "year", "month", "week", "day"}).Aggregate().
				SetGroupByArray([]string{"profile_id", "character_id", "year", "month", "week"}).Aggregate().
				SetGroupByArray([]string{"profile_id", "character_id", "year", "month"}).Aggregate().
				SetGroupByArray([]string{"profile_id", "character_id", "year"}).Aggregate().
				SetGroupByArray([]string{"profile_id", "character_id"}).Aggregate().
				SetGroupByArray([]string{"profile_id"}).Aggregate()

		case FilterTypeCharacter:
			dfa.SetDataframe(dfInnerJoin).SetSumField("time").
				SetGroupByArray([]string{"character_id", "metric_id", "year", "month", "week", "day"}).Aggregate().
				SetGroupByArray([]string{"character_id", "metric_id", "year", "month", "week"}).Aggregate().
				SetGroupByArray([]string{"character_id", "metric_id", "year", "month"}).Aggregate().
				SetGroupByArray([]string{"character_id", "metric_id", "year"}).Aggregate().
				SetGroupByArray([]string{"character_id", "metric_id"}).Aggregate().
				SetGroupByArray([]string{"character_id"}).Aggregate()
		}

		jsonOutput, _ := json.MarshalIndent(ap.AnalyticResults, "", "  ")
		fmt.Println(string(jsonOutput))
	}

	// Process the captured records for the FREQUENCY section
	if lo.Contains(ap.AnalyticSections, model.AnalyticSectionFrequency) {
		maxTotalFocusedTime := float64(0)
		totalFocusedTimeCol := dfCaptureRecords.Col("total_focused_time")
		if totalFocusedTimeCol.Error() == nil {
			maxTotalFocusedTime = totalFocusedTimeCol.Max()
		}

		unitRange := math.Ceil(maxTotalFocusedTime / NUMBER_OF_FREQUENCY_RANGE)

		// Calculate number of days between the start and end time
		numberOfDays := math.Ceil(ap.EndTime.Sub(ap.StartTime).Hours()/float64(NUMBER_OF_HOUR_IN_A_DAY)) + 1
		startWeekDay := int(ap.StartTime.Weekday())
		endWeekDay := int(ap.EndTime.Weekday())
		numberOfDays += float64(startWeekDay)

		// Make the two-dimensional array to store the focused time
		row := int(math.Ceil(numberOfDays / float64(NUMBER_OF_DAYS_IN_A_WEEK)))
		column := NUMBER_OF_DAYS_IN_A_WEEK
		frequencyMatrix := make([][]int, row)
		for i := 0; i < row; i++ {
			frequencyMatrix[i] = make([]int, column)
			for j := 0; j < column; j++ {
				// Fill the days in previous year with -1
				if i == 0 && j < startWeekDay {
					frequencyMatrix[i][j] = -1
				} else {
					frequencyMatrix[i][j] = 0
				}
			}
		}

		// Fill the matrix with the focused time
		for i := 0; i < len(ap.CapturedRecords); i++ {
			record := ap.CapturedRecords[i]
			dayIndex := math.Floor(record.Timestamp.Sub(ap.StartTime).Hours()/float64(NUMBER_OF_HOUR_IN_A_DAY)) + float64(startWeekDay)
			dayRow := int(math.Floor(dayIndex / float64(NUMBER_OF_DAYS_IN_A_WEEK)))
			dayCol := int(math.Mod(dayIndex, float64(NUMBER_OF_DAYS_IN_A_WEEK)))
			frequencyMatrix[dayRow][dayCol] = int(math.Ceil(float64(record.TotalFocusedTime) / unitRange))
		}

		// Fill the last redundant days in the next year with -1
		for j := endWeekDay + 1; j < NUMBER_OF_DAYS_IN_A_WEEK; j++ {
			frequencyMatrix[row-1][j] = -1
		}

		ap.AnalyticResults["frequency"] = frequencyMatrix
	}

	return ap.AnalyticResults
}

type DataframeAggregator struct {
	df              dataframe.DataFrame
	analyticResults map[string]interface{}
	groupByArr      []string
	sumField        string
}

func (dfa *DataframeAggregator) SetGroupByArray(groupByArr []string) *DataframeAggregator {
	dfa.groupByArr = groupByArr
	return dfa
}

func (dfa *DataframeAggregator) SetDataframe(df dataframe.DataFrame) *DataframeAggregator {
	dfa.df = df
	return dfa
}

func (dfa *DataframeAggregator) SetAnalyticResults(analyticResults map[string]interface{}) *DataframeAggregator {
	dfa.analyticResults = analyticResults
	return dfa
}

func (dfa *DataframeAggregator) SetSumField(sumField string) *DataframeAggregator {
	dfa.sumField = sumField
	return dfa
}

func (dfa *DataframeAggregator) Aggregate() *DataframeAggregator {
	// Group by the fields
	sumDF := dfa.df.GroupBy(dfa.groupByArr...).Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{dfa.sumField})

	// Rename the columns
	sumDF = sumDF.Rename(dfa.sumField, dfa.sumField+"_SUM")

	// Add the result to the analytic results
	for i := 0; i < sumDF.Nrow(); i++ {
		// Don't track the zero values
		if value, _ := sumDF.Col(dfa.sumField).Elem(i).Int(); value == 0 {
			continue
		}

		curResultMap := dfa.analyticResults
		for _, field := range dfa.groupByArr {
			value := sumDF.Col(field).Elem(i).String()
			if _, ok := curResultMap[value]; !ok {
				curResultMap[value] = make(map[string]interface{})
			}
			curResultMap = curResultMap[value].(map[string]interface{})
		}
		curResultMap["time"], _ = sumDF.Col(dfa.sumField).Elem(i).Int()
	}

	dfa.df = sumDF

	return dfa
}
