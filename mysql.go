package go_utils

import (
	"fmt"
	"reflect"
)


/**
Method that extracts the data for a Row, stores it in a structure and appends the output array
*/
func ExtractDataFromRowToStructure(output interface{}, rows []string, cols []string, seperator string) (err error) {

	elements := reflect.TypeOf(output).Elem().Elem()
	destinationStructure := reflect.New(elements).Elem()

	titleDB, err := extractNamesAndTagsFromStructure(destinationStructure)

	fmt.Printf("starting row analysis\n")

	for _, row := range rows {

		parts, err := splitRowValues(row, seperator)
		if err != nil {
			return err
		}
		dbase := reflect.ValueOf(output).Elem()
		for k, val := range parts {
			index, err := getFieldIndex(cols[k], titleDB)
			if err != nil {
				return err
			}
			destinationStructure.Field(index).SetString(val)
		}
		dbase.Set(reflect.Append(dbase, destinationStructure))
	}

	return nil
}

/**
Searches for the struct field corresponding to the csv column title
*/
func getFieldIndex(tagName string, titleDB []StructureMatchWithCSV) (index int, err error) {
	for _, v := range titleDB {
		if v.CSVTitle == tagName {
			return v.Index, nil
		}
	}
	return index, fmt.Errorf("Unable to find corresponding field")
}


func extractNamesAndTagsFromStructure(destinationStructure reflect.Value) (data []StructureMatchWithCSV, err error) {
	for i := 0; i < destinationStructure.NumField(); i++ {
		data = append(data, StructureMatchWithCSV{
			Index:          i,
			CSVTitle:       destinationStructure.Type().Field(i).Tag.Get("csv"),
			StructureTitle: destinationStructure.Type().Field(i).Name,
		})
	}
	fmt.Printf("titles: %v\n", data)


	return data,nil
}