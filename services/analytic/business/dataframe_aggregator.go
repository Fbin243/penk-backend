package business

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-gota/gota/dataframe"
)

type DataframeAggregator struct {
	df              dataframe.DataFrame
	analyticResults map[string]any
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

func (dfa *DataframeAggregator) SetAnalyticResults(analyticResults map[string]any) *DataframeAggregator {
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
	for i := range sumDF.Nrow() {
		// Don't track the zero values
		if value, _ := sumDF.Col(dfa.sumField).Elem(i).Int(); value == 0 {
			continue
		}

		curResultMap := dfa.analyticResults
		for _, field := range dfa.groupByArr {
			value := sumDF.Col(field).Elem(i).String()
			if _, ok := curResultMap[value]; !ok {
				curResultMap[value] = make(map[string]any)
			}
			curResultMap = curResultMap[value].(map[string]any)
		}
		curResultMap["time"], _ = sumDF.Col(dfa.sumField).Elem(i).Int()
	}

	dfa.df = sumDF

	return dfa
}

func PrintDF(df *dataframe.DataFrame) {
	records := df.Records()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)

	for _, record := range records[0:] {
		fmt.Fprintln(w, strings.Join(record, "\t")+"\t")
	}

	w.Flush()
}
