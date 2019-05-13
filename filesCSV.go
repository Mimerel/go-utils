package go_utils

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type CSVFileStructure struct {
	File                 string
	Separator            string
	Output               interface{}
	Titles               []string
	ColumnTitle          bool
	LoggerInfo           func(string, ...interface{})
	LoggerError          func(string, ...interface{})
	Hook                 func(...interface{}) (err error)
	HookEvery            int
	HookArgs             interface{}
	HookResetOutput      bool
	CounterDisplay       int
	fileContent          *os.File
	scannedRowsCount     int64
	scannedRowDetails    string
	destinationStructure reflect.Value
	TitleDB              []StructureMatchWithCSV
	Debug                bool
}

type StructureMatchWithCSV struct {
	Index          int
	CSVTitle       string
	StructureTitle string
}

func NewCSVFile() *CSVFileStructure {
	return &CSVFileStructure{}
}

/**
Initialize values
*/
func (f *CSVFileStructure) init() (err error) {
	// init variables
	f.scannedRowsCount = 0
	if f.LoggerInfo == nil {
		f.LoggerInfo = DefaultLogOutput
	}
	if f.LoggerError == nil {
		f.LoggerError = DefaultLogOutput
	}
	if f.Separator == "" {
		f.Separator = ";"
	}
	return nil
}

/**
Extract titles and tags from output structure
*/
func (f *CSVFileStructure) extractTitlesAndTagsFromStructure() (err error) {

	f.TitleDB = []StructureMatchWithCSV{}
	// Analysing Structure details
	elements := reflect.TypeOf(f.Output).Elem().Elem()
	f.destinationStructure = reflect.New(elements).Elem()
	for i := 0; i < f.destinationStructure.NumField(); i++ {
		f.TitleDB = append(f.TitleDB, StructureMatchWithCSV{
			Index:          i,
			CSVTitle:       f.destinationStructure.Type().Field(i).Tag.Get("csv"),
			StructureTitle: f.destinationStructure.Type().Field(i).Name,
		})
	}
	return nil
}

/**
Check if output structure is an Array
*/
func (f *CSVFileStructure) check() (err error) {
	if StringInArray(string(reflect.TypeOf(f.Output).Kind()), []string{"Slice", "Array"}) {
		return fmt.Errorf("Destination variable is not an Array")
	}
	return nil
}

/*
Main method that converts the CSV to a flat structure
*/
func (f *CSVFileStructure) UnmarshalCSV() (err error) {
	err = f.check()
	if err != nil {
		return err
	}

	err = f.init()
	if err != nil {
		return err
	}

	err = f.extractTitlesAndTagsFromStructure()
	if err != nil {
		return err
	}

	// Open file to unmarshal
	f.fileContent, err = os.Open(f.File)
	if err != nil {
		f.LoggerError("%v", err)
		return err
	}
	defer f.fileContent.Close()

	// Starting to scan file
	err = f.scanFile()
	if err != nil {
		return err
	}

	f.LoggerInfo("Unmarchal process completed")
	return nil
}

/**
Method that reads line by line the csv file
*/
func (f *CSVFileStructure) scanFile() (err error) {
	scanner := bufio.NewScanner(f.fileContent)
	counter := 0
	for scanner.Scan() {
		f.scannedRowDetails = scanner.Text()
		f.scannedRowsCount++
		counter++

		if f.scannedRowsCount == 1 {
			if f.ColumnTitle {
				f.Titles = strings.Split(f.scannedRowDetails, f.Separator)
				f.LoggerInfo("Found titles")
				continue
			} else {
				f.LoggerInfo("Using titles inputed by user")
			}
		}

		// Extracting Row content to store it in destination variable
		err = f.extractDataFromRow()
		if err != nil {
			f.LoggerError("Error deserializing line : %s", f.scannedRowDetails)
			return err
		}
		if counter >= f.HookEvery {
			err = f.Hook(f.Output, f.scannedRowsCount, f.HookArgs)
			if err != nil {
				return err
			}
			if f.HookResetOutput {
				// Reset the output slice to empty
				v := reflect.ValueOf(f.Output)
				v.Elem().Set(reflect.MakeSlice(v.Type().Elem(), 0, v.Elem().Cap()))
			}
			counter = 0
		}
	}

	if err := scanner.Err(); err != nil {
		f.LoggerError("Error scanning file")
		return err
	}
	return nil
}

/**
Method that extracts the data for a Row, stores it in a structure and appends the output array
*/
func (f *CSVFileStructure) extractDataFromRow() (err error) {

	parts, err := splitRowValues(f.scannedRowDetails, f.Separator, f.Debug)
	if f.Debug {
		f.LoggerInfo("Rows : %s", f.scannedRowDetails)
		f.LoggerInfo("parts : %v", parts)
	}
	dbase := reflect.ValueOf(f.Output).Elem()
	for k, val := range parts {
		if k >= len(f.TitleDB) {
			return fmt.Errorf("More parts than columns in table/n cols: %v row: %v parts: %v", f.TitleDB, f.scannedRowDetails, parts)
		}
		index, err := f.getFieldIndex(f.Titles[k])
		if err != nil {
			f.LoggerInfo("title list : %s", f.TitleDB)
			f.LoggerInfo("field value : %s", val)
			f.LoggerInfo("field index : %v", k)
			return err
		}

		if index >= len(f.TitleDB) {
			return fmt.Errorf("Index out of range %v\n, for record %v", index, parts)
		}
		f.destinationStructure.Field(index).SetString(val)
	}

	dbase.Set(reflect.Append(dbase, f.destinationStructure))
	return nil
}

/**
Searches for the struct field corresponding to the csv column title
*/
func (f *CSVFileStructure) getFieldIndex(tagName string) (index int, err error) {
	for _, v := range f.TitleDB {
		if v.CSVTitle == tagName {
			return v.Index, nil
		}
	}
	return index, fmt.Errorf("Unable to find corresponding field")
}
