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

// Process the captured records after filtering them to get the analytics results
var processCapturedRecords = func(filterType FilterType, capturedRecords []model.CapturedRecord) map[string]interface{} {
	analyticsResult := map[string]interface{}{}
	if len(capturedRecords) == 0 {
		return analyticsResult
	}

	dfCaptureRecordsData := make([]repo.DFCapturedRecord, 0)
	dfCaptureRecordCustomMetricsData := make([]repo.DFCapturedRecordCustomMetric, 0)

	uniqTimestamps := lo.Map(capturedRecords, func(record model.CapturedRecord, index int) time.Time {
		recordDay := record.Timestamp.Day()
		recordWeek := math.Min(math.Ceil(float64(recordDay)/NUMBER_OF_DAYS_IN_A_WEEK), NUMBER_OF_WEEKS_IN_A_MONTH)

		dfCaptureRecordsData = append(dfCaptureRecordsData, repo.DFCapturedRecord{
			ID:               record.ID.Hex(),
			CharacterID:      record.Metadata.CharacterID.Hex(),
			Year:             record.Timestamp.Year(),
			Month:            int(record.Timestamp.Month()),
			Week:             int(recordWeek),
			Day:              recordDay,
			TotalFocusedTime: int(record.TotalFocusedTime),
		})

		for _, metric := range record.CustomMetrics {
			dfCaptureRecordCustomMetricsData = append(dfCaptureRecordCustomMetricsData, repo.DFCapturedRecordCustomMetric{
				ID:       record.ID.Hex(),
				MetricID: metric.ID.Hex(),
				Time:     int(metric.Time),
			})
		}

		return record.Timestamp
	})

	// Total focused days
	analyticsResult["days"] = len(lo.Uniq(lo.Map(uniqTimestamps, func(timestamp time.Time, index int) string {
		return timestamp.Format(time.DateOnly)
	})))

	// Best & Current streak (count continuous days)
	bestStreak := 0
	currentStreak := 0
	for i := 0; i < len(uniqTimestamps); i++ {
		fmt.Print(uniqTimestamps[i].Format(time.RFC3339) + "\n")
		if i == 0 || uniqTimestamps[i].Sub(uniqTimestamps[i-1]) == 24*time.Hour {
			currentStreak++
		} else {
			bestStreak = int(math.Max(float64(bestStreak), float64(currentStreak)))
			currentStreak = 1
		}
	}
	bestStreak = int(math.Max(float64(bestStreak), float64(currentStreak)))

	analyticsResult["bestStreak"] = bestStreak
	analyticsResult["currentStreak"] = currentStreak
	analyticsResult["startDay"] = uniqTimestamps[0].Format(time.DateOnly)
	analyticsResult["endDay"] = uniqTimestamps[len(uniqTimestamps)-1].Format(time.DateOnly)

	dfCaptureRecords := dataframe.LoadStructs(dfCaptureRecordsData)
	dfCaptureRecordMetrics := dataframe.LoadStructs(dfCaptureRecordCustomMetricsData)
	dfInnerJoin := dfCaptureRecords.InnerJoin(dfCaptureRecordMetrics, "id")

	switch filterType {
	case FilterTypeUser:
		fmt.Println("\n================= CHARACTER DAILY =================")
		dfCharacterDaily := aggregateAnalyticsData(dfCaptureRecords, &analyticsResult, []string{"character_id", "year", "month", "week", "day"}, "total_focused_time")

		fmt.Println("\n================= CHARACTER WEEKLY =================")
		dfCharacterWeekly := aggregateAnalyticsData(dfCharacterDaily, &analyticsResult, []string{"character_id", "year", "month", "week"}, "total_focused_time")

		fmt.Println("\n================= CHARACTER MONTHLY =================")
		dfCharacterMonthly := aggregateAnalyticsData(dfCharacterWeekly, &analyticsResult, []string{"character_id", "year", "month"}, "total_focused_time")

		fmt.Println("\n================= CHARACTER YEARLY =================")
		dfCharacterYearly := aggregateAnalyticsData(dfCharacterMonthly, &analyticsResult, []string{"character_id", "year"}, "total_focused_time")

		fmt.Println("\n================= CHARACTER =================")
		aggregateAnalyticsData(dfCharacterYearly, &analyticsResult, []string{"character_id"}, "total_focused_time")

	case FilterTypeCharacter:
		fmt.Println("\n================= DAILY =================")
		dfMetricDaily := aggregateAnalyticsData(dfInnerJoin, &analyticsResult, []string{"metric_id", "year", "month", "week", "day"}, "time")

		fmt.Println("\n================= WEEKLY =================")
		dfMetricWeekly := aggregateAnalyticsData(dfMetricDaily, &analyticsResult, []string{"metric_id", "year", "month", "week"}, "time")

		fmt.Println("\n================= MONTHLY =================")
		dfMetricMonthly := aggregateAnalyticsData(dfMetricWeekly, &analyticsResult, []string{"metric_id", "year", "month"}, "time")

		fmt.Println("\n================= YEARLY =================")
		dfMetricYearly := aggregateAnalyticsData(dfMetricMonthly, &analyticsResult, []string{"metric_id", "year"}, "time")

		fmt.Println("\n================= METRIC =================")
		aggregateAnalyticsData(dfMetricYearly, &analyticsResult, []string{"metric_id"}, "time")
	}

	jsonOutput, _ := json.MarshalIndent(analyticsResult, "", "  ")
	fmt.Println(string(jsonOutput))

	return analyticsResult
}

var aggregateAnalyticsData = func(df dataframe.DataFrame, resultMap *map[string]interface{}, groupByArr []string, sumField string) dataframe.DataFrame {
	// Group by the fields
	sumDF := df.GroupBy(groupByArr...).Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{sumField})

	// Rename the column
	sumDF = sumDF.Rename(sumField, sumField+"_SUM")

	// Add the result to the result map
	for i := 0; i < sumDF.Nrow(); i++ {
		// Don't track the zero values
		if value, _ := sumDF.Col(sumField).Elem(i).Int(); value == 0 {
			continue
		}

		curResultMap := *resultMap
		for _, field := range groupByArr {
			value := sumDF.Col(field).Elem(i).String()
			fmt.Print(value + " ")
			if _, ok := curResultMap[value]; !ok {
				curResultMap[value] = make(map[string]interface{})
			}
			curResultMap = curResultMap[value].(map[string]interface{})
		}
		fmt.Println()
		curResultMap["time"], _ = sumDF.Col(sumField).Elem(i).Int()
	}

	return sumDF
}
