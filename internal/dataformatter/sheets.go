package dataformatter

import (
	"strconv"

	"google.golang.org/api/sheets/v4"
)

func FormatToRowValueRange(values []string) []interface{} {
	valueRange := make([]interface{}, len(values))
	for i, v := range values {
		valueRange[i] = v
	}
	return valueRange
}

func FormatToValueRange(values [][]string) [][]interface{} {
	valueRange := make([][]interface{}, len(values))
	for i, v := range values {
		valueRange[i] = make([]interface{}, len(v))
		for j, val := range v {
			valueRange[i][j] = val
		}
	}
	return valueRange
}

func FormatValueRangeBatch(tab string, rows []ExistingRow) []*sheets.ValueRange {
	ranges := make([]*sheets.ValueRange, len(rows))
	for i, row := range rows {
		rangeIndex := strconv.Itoa(row.index + 1)
		ranges[i] = &sheets.ValueRange{
			Range:  tab + "!A" + rangeIndex + ":ZZ" + rangeIndex,
			Values: [][]interface{}{FormatToRowValueRange(row.data)},
		}
	}
	return ranges
}

func TabExists(search string, sheets []*sheets.Sheet) bool {
	for _, sheets := range sheets {
		if search == sheets.Properties.Title {
			return true
		}
	}
	return false
}
