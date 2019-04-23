package go_utils

import (
	"fmt"
	"reflect"
	"strconv"
)

/**
Method that extracts the data for a Row, stores it in a structure and appends the output array
*/
func ExtractDataFromRowToStructure(output interface{}, rows []string, cols []string, seperator string, debug bool) (err error) {

	elements := reflect.TypeOf(output).Elem().Elem()
	destinationStructure := reflect.New(elements).Elem()

	titleDB, err := extractNamesAndTagsFromStructure(destinationStructure)

	for _, row := range rows {

		parts, err := splitRowValues(row, seperator, debug)
		if err != nil {
			return err
		}
		dbase := reflect.ValueOf(output).Elem()
		for k, val := range parts {
			index, err := getFieldIndex(cols[k], titleDB)
			if err != nil {
				return err
			}
			switch destinationStructure.Field(index).Kind() {
			case reflect.String:
				destinationStructure.Field(index).SetString(val)
			case reflect.Int64:
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				destinationStructure.Field(index).SetInt(valInt)
			case reflect.Bool:
				valBool, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				destinationStructure.Field(index).SetBool(valBool)
			default:
				destinationStructure.Field(index).SetString(val)
			}
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
	fmt.Printf("looking for field %s\n", tagName)
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

	return data, nil
}
