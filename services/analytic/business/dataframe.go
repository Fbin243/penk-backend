package business

import (
	"math"
	"time"

	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/entity"

	"github.com/go-gota/gota/dataframe"
	"github.com/samber/lo"
)

type FilterType int

const (
	FilterTypeUser FilterType = iota
	FilterTypeCharacter
)

type AnalyticsProcessor struct {
	AnalyticSections []entity.AnalyticSection
	AnalyticResults  map[string]any
	CapturedRecords  []entity.CapturedRecord
	CharacterID      string
	StartTime        *time.Time
	EndTime          *time.Time
}

// ProcessCapturedRecords processes the captured records and returns the analytic results
func (ap *AnalyticsProcessor) ProcessCapturedRecords() (map[string]any, error) {
	capturedRecordsNum := len(ap.CapturedRecords)
	dfCapturedRecords := dataframe.LoadStructs(ap.CapturedRecords)
	// PrintDF(&dfCapturedRecords)

	// ---------------------------- OVERALL ----------------------------------
	if capturedRecordsNum > 0 && lo.Contains(ap.AnalyticSections, entity.AnalyticSectionOverall) {
		// Total active days (count distinct days)
		uniqTimeStampStrings := lo.Uniq(lo.Map(ap.CapturedRecords, func(record entity.CapturedRecord, index int) string {
			return record.Date
		}))

		totalFocusedDays := len(uniqTimeStampStrings)
		totalFocusedTime := dfCapturedRecords.Col("time").Sum()

		// Best & Current streak (count continuous days)
		timeStamps := lo.Map(uniqTimeStampStrings, func(timeStampString string, index int) time.Time {
			timeStamp, _ := time.Parse(time.DateOnly, timeStampString)
			return timeStamp
		})
		timeStamps = append(timeStamps, time.Time{})

		bestStreak := 1
		currentStreak := 1
		bestStreakStartDate := timeStamps[0]
		currentStreakStartDate := timeStamps[0]
		for i := 1; i < len(timeStamps); i++ {
			if timeStamps[i].Sub(timeStamps[i-1]) == 24*time.Hour {
				currentStreak++
				continue
			}
			if currentStreak > bestStreak {
				bestStreak = currentStreak
				bestStreakStartDate = currentStreakStartDate
			}
			if i != len(timeStamps)-1 {
				currentStreak = 1
				currentStreakStartDate = timeStamps[i]
			}
		}

		ap.AnalyticResults["overall"] = make(map[string]any)
		analyticResultsOverall := ap.AnalyticResults["overall"].(map[string]any)
		analyticResultsOverall["totalFocusedTime"] = totalFocusedTime
		analyticResultsOverall["totalFocusedDays"] = totalFocusedDays
		analyticResultsOverall["averageFocusedTime"] = totalFocusedTime / float64(totalFocusedDays)
		analyticResultsOverall["bestStreak"] = bestStreak
		analyticResultsOverall["bestStreakStartDate"] = bestStreakStartDate
		analyticResultsOverall["bestStreakEndDate"] = bestStreakStartDate.AddDate(0, 0, bestStreak-1)

		currentStreakEndDate := currentStreakStartDate.AddDate(0, 0, currentStreak-1)
		if utils.Now().Sub(currentStreakEndDate) > 24*time.Hour {
			analyticResultsOverall["currentStreak"] = 0
			analyticResultsOverall["currentStreakStartDate"] = nil
			analyticResultsOverall["currentStreakEndDate"] = nil
		} else {
			analyticResultsOverall["currentStreak"] = currentStreak
			analyticResultsOverall["currentStreakStartDate"] = currentStreakStartDate
			analyticResultsOverall["currentStreakEndDate"] = currentStreakEndDate
		}
	}

	// ---------------------------- DISTRIBUTION ----------------------------------
	if capturedRecordsNum > 0 && lo.Contains(ap.AnalyticSections, entity.AnalyticSectionDistribution) {
		totalFocusedTime := dfCapturedRecords.Col("time").Sum()
		ap.AnalyticResults["distribution"] = make(map[string]any)
		dfDistribution := dfCapturedRecords.GroupBy("character_id", "category_id").Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"time"})
		// Calculate the percentage of the total focused time
		for i := range dfDistribution.Nrow() {
			categoryID := dfDistribution.Col("category_id").Elem(i).String()
			categoryTime := dfDistribution.Col("time_SUM").Elem(i).Float()
			ap.AnalyticResults["distribution"].(map[string]any)[categoryID] = categoryTime
		}

		ap.AnalyticResults["distribution"].(map[string]any)["totalFocusedTime"] = totalFocusedTime
	}

	// ---------------------------- TIMELINE ----------------------------------
	if capturedRecordsNum > 0 && lo.Contains(ap.AnalyticSections, entity.AnalyticSectionTimeline) {
		dfa := &DataframeAggregator{}
		ap.AnalyticResults["timeline"] = make(map[string]any)
		analyticResultsTimeline := ap.AnalyticResults["timeline"].(map[string]any)
		dfa.SetAnalyticResults(analyticResultsTimeline)
		dfa.SetDataframe(dfCapturedRecords).SetSumField("time").
			SetGroupByArray([]string{"character_id", "category_id", "year", "month", "week", "day"}).Aggregate().
			SetGroupByArray([]string{"character_id", "category_id", "year", "month", "week"}).Aggregate().
			SetGroupByArray([]string{"character_id", "category_id", "year", "month"}).Aggregate().
			SetGroupByArray([]string{"character_id", "category_id", "year"}).Aggregate().
			SetGroupByArray([]string{"character_id", "category_id"}).Aggregate().
			SetGroupByArray([]string{"character_id"}).Aggregate()
	}

	// ---------------------------- FREQUENCY ----------------------------------
	if lo.Contains(ap.AnalyticSections, entity.AnalyticSectionFrequency) {
		if ap.StartTime == nil || ap.EndTime == nil {
			return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "Start and end time are required for frequency analysis")
		}

		maxTime := float64(0)
		dfFrequency := dfCapturedRecords.GroupBy("character_id", "year", "month", "week", "day", "date").Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"time"})
		timeCol := dfFrequency.Col("time_SUM")
		if timeCol.Error() == nil {
			maxTime = timeCol.Max()
		}
		unitRange := math.Ceil(maxTime / NUMBER_OF_FREQUENCY_RANGE)

		// Calculate number of days between the start and end time
		numberOfDays := ap.EndTime.Sub(*ap.StartTime).Hours()/float64(NUMBER_OF_HOUR_IN_A_DAY) + 1
		startWeekDay := int(ap.StartTime.Weekday())
		endWeekDay := int(ap.EndTime.Weekday())
		numberOfDays += float64(startWeekDay)

		// Make the two-dimensional array to store the focused time
		row := int(math.Ceil(numberOfDays / float64(NUMBER_OF_DAYS_IN_A_WEEK)))
		column := NUMBER_OF_DAYS_IN_A_WEEK
		frequencyMatrix := make([][]int, row)
		for i := range row {
			frequencyMatrix[i] = make([]int, column)
			for j := range column {
				// Fill the days in previous year with -1
				if i == 0 && j < startWeekDay {
					frequencyMatrix[i][j] = -1
				} else {
					frequencyMatrix[i][j] = 0
				}
			}
		}

		// Fill the matrix with the focused time
		for i := range dfFrequency.Nrow() {
			date, _ := time.Parse(time.DateOnly, dfFrequency.Col("date").Elem(i).String())
			time, _ := dfFrequency.Col("time_SUM").Elem(i).Int()
			dayIndex := math.Floor(date.Sub(*ap.StartTime).Hours()/float64(NUMBER_OF_HOUR_IN_A_DAY)) + float64(startWeekDay)
			dayRow := int(math.Floor(dayIndex / float64(NUMBER_OF_DAYS_IN_A_WEEK)))
			dayCol := int(math.Mod(dayIndex, float64(NUMBER_OF_DAYS_IN_A_WEEK)))
			frequencyMatrix[dayRow][dayCol] = int(math.Ceil(float64(time) / unitRange))
		}

		// Fill the last redundant days in the next year with -1
		for j := endWeekDay + 1; j < NUMBER_OF_DAYS_IN_A_WEEK; j++ {
			frequencyMatrix[row-1][j] = -1
		}

		ap.AnalyticResults["frequency"] = frequencyMatrix
	}

	return ap.AnalyticResults, nil
}
