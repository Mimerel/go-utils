package go_utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/**
Structure that groups the main data required to do the extraction to Structure
Rows = rows received from mysql
cols = name of columns from mysql
seperator = seperator used by mysql
*/
type ExtractDataOptions struct {
	Rows               []string
	Cols               []string
	Seperator          string
	Debug              bool
	RemoveDoubleSpaces bool
	RemoveEndSpace     bool
	RemoveStartSpace   bool
}

/**
Method that extracts the data for a Row, stores it in a structure and appends the output array
*/
func ExtractDataFromRowToStructure(output interface{}, params ExtractDataOptions) (err error) {

	elements := reflect.TypeOf(output).Elem().Elem()
	destinationStructure := reflect.New(elements).Elem()

	titleDB, err := extractNamesAndTagsFromStructure(destinationStructure)
	if params.Debug {
		fmt.Printf("Title dbase\n")
		fmt.Printf("----\n")
		for _, v := range titleDB {
			fmt.Printf("Index: %v\n", v.Index)
			fmt.Printf("csv Title: %v\n", v.CSVTitle)
			fmt.Printf("str Title: %v\n", v.StructureTitle)
			fmt.Printf("----")
		}
	}
	for _, row := range params.Rows {

		parts, err := splitRowValues(row, params.Seperator, params.Debug)
		if err != nil {
			return err
		}
		dbase := reflect.ValueOf(output).Elem()
		for k, val := range parts {
			index, err := getFieldIndex(params.Cols[k], titleDB)
			if err != nil {
				return err
			}
			switch destinationStructure.Field(index).Kind() {
			case reflect.String:
				destinationStructure.Field(index).SetString(transforedString(params, val))
			case reflect.Int64:
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				destinationStructure.Field(index).SetInt(valInt)
			case reflect.Int:
				valInt, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				destinationStructure.Field(index).SetCap(valInt)
			case reflect.Bool:
				valBool, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				destinationStructure.Field(index).SetBool(valBool)
			default:
				destinationStructure.Field(index).SetString(transforedString(params, val))
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
	return index, fmt.Errorf("Unable to find corresponding field %s\n", tagName)
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
func transforedString(params ExtractDataOptions, value string) (result string) {
	result = value
	if params.RemoveDoubleSpaces {
		for strings.Index(result, "  ") != -1 {
			result = strings.Replace(result, "  ", " ", -1)
			fmt.Printf("value : %s", strings.Replace(result, " ", "€", -1))
		}
	}
	if params.RemoveStartSpace {
		result = strings.TrimLeft(result, " ")
	}
	if params.RemoveEndSpace {
		result = strings.TrimRight(result, " ")
	}
	return result
}
