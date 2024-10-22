package business

import (
	"encoding/json"
	"fmt"
	"math"

	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/analytics/repo"

	"github.com/go-gota/gota/dataframe"
)

var processCapturedRecords = func(capturedRecords []model.CapturedRecord) map[string]interface{} {
	dfCaptureRecordsData := make([]repo.DFCapturedRecord, 0)
	dfCaptureRecordCustomMetricsData := make([]repo.DFCapturedRecordCustomMetric, 0)

	for _, record := range capturedRecords {
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
	}

	dfCaptureRecords := dataframe.LoadStructs(dfCaptureRecordsData)
	dfCaptureRecordMetrics := dataframe.LoadStructs(dfCaptureRecordCustomMetricsData)
	dfInnerJoin := dfCaptureRecords.InnerJoin(dfCaptureRecordMetrics, "id")

	analyticsResult := map[string]interface{}{}

	fmt.Println("\n================= DAILY =================")
	dfMetricDaily := aggregateAnalyticsData(dfInnerJoin, &analyticsResult, []string{"character_id", "metric_id", "year", "month", "week", "day"}, "time")

	fmt.Println("\n================= WEEKLY =================")
	dfMetricWeekly := aggregateAnalyticsData(dfMetricDaily, &analyticsResult, []string{"character_id", "metric_id", "year", "month", "week"}, "time")

	fmt.Println("\n================= MONTHLY =================")
	dfMetricMonthly := aggregateAnalyticsData(dfMetricWeekly, &analyticsResult, []string{"character_id", "metric_id", "year", "month"}, "time")

	fmt.Println("\n================= YEARLY =================")
	dfMetricYearly := aggregateAnalyticsData(dfMetricMonthly, &analyticsResult, []string{"character_id", "metric_id", "year"}, "time")

	fmt.Println("\n================= METRIC =================")
	aggregateAnalyticsData(dfMetricYearly, &analyticsResult, []string{"character_id", "metric_id"}, "time")

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
		curResultMap := *resultMap
		for _, field := range groupByArr {
			value := sumDF.Col(field).Elem(i).String()
			fmt.Print(value + " ")
			if field == "metric_id" {
				if _, ok := curResultMap["metrics"]; !ok {
					curResultMap["metrics"] = make(map[string]interface{})
				}
				curResultMap = curResultMap["metrics"].(map[string]interface{})
			}

			if _, ok := curResultMap[value]; !ok {
				curResultMap[value] = make(map[string]interface{})
			}
			curResultMap = curResultMap[value].(map[string]interface{})
		}
		fmt.Print(sumDF.Col(sumField).Elem(i).Int())
		fmt.Println()
		curResultMap["value"], _ = sumDF.Col(sumField).Elem(i).Int()
	}

	return sumDF
}
