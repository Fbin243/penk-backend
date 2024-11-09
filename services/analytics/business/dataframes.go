package business

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

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
	if len(ap.CapturedRecords) == 0 {
		return ap.AnalyticResults
	}

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

	// Process the captured records for the OVERALL section
	if lo.Contains(ap.AnalyticSections, model.AnalyticSectionOverall) {
		// Total active days
		totalFocusedDays := dfCaptureRecords.Nrow()

		// Total focused time
		totalFocusedTime := dfCaptureRecords.Col("total_focused_time").Sum()

		// Best & Current streak (count continuous days)
		bestStreak := 0
		currentStreak := 0
		timeStamps := lo.Map(ap.CapturedRecords, func(record model.CapturedRecord, index int) time.Time {
			return record.Timestamp
		})
		for i := 0; i < len(timeStamps); i++ {
			if i == 0 || timeStamps[i].Sub(timeStamps[i-1]) == 24*time.Hour {
				currentStreak++
			} else {
				bestStreak = int(math.Max(float64(bestStreak), float64(currentStreak)))
				currentStreak = 1
			}
		}
		bestStreak = int(math.Max(float64(bestStreak), float64(currentStreak)))

		ap.AnalyticResults["totalFocusedTime"] = totalFocusedTime
		ap.AnalyticResults["totalFocusedDays"] = totalFocusedDays
		ap.AnalyticResults["avarageFocusedTime"] = totalFocusedTime / float64(totalFocusedDays)
		ap.AnalyticResults["bestStreak"] = bestStreak
		ap.AnalyticResults["currentStreak"] = currentStreak
		ap.AnalyticResults["startDate"] = timeStamps[0].Format(time.DateOnly)
		ap.AnalyticResults["endDate"] = timeStamps[len(timeStamps)-1].Format(time.DateOnly)
	}

	// Process the captured records for the DISTRIBUTION section
	if lo.Contains(ap.AnalyticSections, model.AnalyticSectionDistribution) {
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
	if lo.Contains(ap.AnalyticSections, model.AnalyticSectionTimeline) {
		dfa := &DataframeAggregator{}
		dfa.SetAnalyticResults(ap.AnalyticResults)
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
		maxTotalFocusedTime := dfCaptureRecords.Col("total_focused_time").Max()
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

		// Fill the last redundant days with -1
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
			fmt.Print(value + " ")
			if _, ok := curResultMap[value]; !ok {
				curResultMap[value] = make(map[string]interface{})
			}
			curResultMap = curResultMap[value].(map[string]interface{})
		}
		fmt.Println()
		curResultMap["time"], _ = sumDF.Col(dfa.sumField).Elem(i).Int()
	}

	dfa.df = sumDF

	return dfa
}
