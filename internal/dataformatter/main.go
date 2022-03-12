package dataformatter

import (
	"fmt"
	"sort"

	"github.com/jomei/notionapi"
)

type ExistingRow struct {
	index int
	data  []string
}

func GenerateNotionFields(notionProperties notionapi.PropertyConfigs) ([]string, []string) {
	properties := []string{}
	for key := range notionProperties {
		properties = append(properties, key)
	}
	sort.Strings(properties)
	fields := []string{}
	for _, key := range properties {
		structure, complex := readNotionPropertyStructure(key, notionProperties[key])
		if complex {
			for _, suffix := range structure {
				fields = append(fields, key+"----"+suffix)
			}
		} else {
			fields = append(fields, key)
		}
	}
	return properties, fields
}

func SplitNotionData(rows [][]string, existingIds [][]interface{}) ([]ExistingRow, [][]string) {
	ids := make(map[string]int)
	for index, id := range existingIds {
		if len(id) > 0 {
			ids[id[0].(string)] = index
		}
	}
	fmt.Println(ids)
	existingRows := []ExistingRow{}
	newRows := [][]string{}
	for _, row := range rows {
		index, exists := ids[row[0]]
		if exists {
			existingRows = append(existingRows, ExistingRow{
				index: index,
				data:  row,
			})
		} else {
			newRows = append(newRows, row)
		}
	}
	return existingRows, newRows
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
