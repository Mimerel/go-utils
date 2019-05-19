package go_utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	if params.Debug {
		fmt.Printf("reflect.TypeOf(output): %v\n", reflect.TypeOf(output))
		fmt.Printf("reflect.TypeOf(output).Elem(): %v\n", reflect.TypeOf(output).Elem())
		fmt.Printf("reflect.TypeOf(output).Elem().Elem(): %v\n", reflect.TypeOf(output).Elem().Elem())
	}
	
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
				if val == "" {
					val = "False"
				}
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



type MariaDBConfiguration struct {
	LoggerInfo  func(string, ...interface{})
	LoggerError func(string, ...interface{})
	User        string
	Password    string
	Database    string
	IP          string
	Port        string
	DB          *sql.DB
	Seperator   string
	WhereClause string
	Table       string
	DataType    interface{}
}

type StructureDetails struct {
	Index     int
	FieldTag  string
	FieldName string
}

type SelectResponse struct {
	Columns   []string
	Rows      []string
	Seperator string
}

func NewMariaDB() *MariaDBConfiguration {
	return &MariaDBConfiguration{}
}

const (
	Low_Priority = "LOW_PRIORITY"
	Normal       = ""
)

/**
Initialize values
*/
func (c *MariaDBConfiguration) init() (err error) {
	// init variables
	if c.IP == "" {
		c.IP = "localhost"
	}
	if c.LoggerInfo == nil {
		c.LoggerInfo = DefaultLogOutput
	}
	if c.LoggerError == nil {
		c.LoggerError = DefaultLogOutput
	}
	if c.Seperator == "" {
		c.Seperator = ";"
	}
	return nil
}

func (c *MariaDBConfiguration) connectMariaDb() (err error) {
	c.init()
	c.DB, err = sql.Open("mysql", c.User+":"+c.Password+"@tcp("+c.IP+":"+c.Port+")/"+c.Database)
	if err != nil {
		c.LoggerError("Unable to create connexion to MariaDb")
		return err
	}

	return nil
}

func (c *MariaDBConfiguration) DecryptStructureAndData(data interface{}) (columns string, values string, err error) {

	var valuesBuilder strings.Builder
	var columnsBuilder strings.Builder
	titleDB := []StructureDetails{}

	// Analysing Structure details
	elements := reflect.TypeOf(data).Elem()
	structureModel := reflect.New(elements).Elem()

	for i := 0; i < structureModel.NumField(); i++ {
		if structureModel.Type().Field(i).Tag.Get("csv") != "" {
			titleDB = append(titleDB, StructureDetails{
				Index:     i,
				FieldTag:  structureModel.Type().Field(i).Tag.Get("csv"),
				FieldName: structureModel.Type().Field(i).Name,
			})
		}
	}

	_, _ = fmt.Fprintf(&columnsBuilder, "%s", "(")
	for k, v := range titleDB {
		if k != 0 {
			_, _ = fmt.Fprintf(&columnsBuilder, "%s", ",")
		}
		_, _ = fmt.Fprintf(&columnsBuilder, "`%s`", v.FieldTag)
	}
	_, _ = fmt.Fprintf(&columnsBuilder, "%s", ")")

	columns = columnsBuilder.String()

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(data)
		for i := 0; i < v.Len(); i++ {
			var subValuesBuilder strings.Builder
			if i != 0 {
				_, _ = fmt.Fprintf(&valuesBuilder, "%s", ",")
			}
			_, _ = fmt.Fprintf(&subValuesBuilder, "%s", "(")
			for k1, v1 := range titleDB {
				// Finding correct name of column corresponding to field
				if k1 != 0 {
					_, _ = fmt.Fprintf(&subValuesBuilder, "%s", ",")
				}
				// Change output depending on type of field to import
				switch v.Index(i).Field(v1.Index).Kind() {
				case reflect.String:
					valueString := strings.Replace(v.Index(i).Field(v1.Index).String(), "\\", "", -1)
					valueString = strings.Replace(valueString, "\"", "", -1)
					_, _ = fmt.Fprintf(&subValuesBuilder, "%s%s%s", "\"", valueString, "\"")
				case reflect.Int64:
					_, _ = fmt.Fprintf(&subValuesBuilder, "%s", strconv.FormatInt(v.Index(i).Field(v1.Index).Int(), 10))
				case reflect.Bool:
					_, _ = fmt.Fprintf(&subValuesBuilder, "%s", strconv.FormatBool(v.Index(i).Field(v1.Index).Bool()))
				default:
					_, _ = fmt.Fprintf(&subValuesBuilder, "%s", v.Index(i).Field(v1.Index).String())
				}

			}
			_, _ = fmt.Fprintf(&subValuesBuilder, "%s", ")")
			_, _ = fmt.Fprintf(&valuesBuilder, "%s", subValuesBuilder.String())
		}

	}
	values = valuesBuilder.String()
	return columns, values, nil
}

func (c *MariaDBConfiguration) Replace(priority string, table string, col string, val string) (err error) {
	err = c.connectMariaDb()
	if err != nil {
		c.LoggerError("Unable to connect to database")
		return err
	}
	defer c.DB.Close()
	sqlRequest := "REPLACE " + priority + " INTO " + table + " " + col + " VALUES " + val
	request, err := c.DB.Prepare(sqlRequest)
	if err != nil {
		c.LoggerError("Unable to prepare Replace request")
		return err
	}

	_, err = request.Exec()
	if err != nil {
		c.LoggerError("Unable to execure Replace request")
		return err
	}

	return nil
}

func (c *MariaDBConfiguration) Insert(ignore bool, table string, col string, val string) (err error) {
	err = c.connectMariaDb()
	if err != nil {
		c.LoggerError("Unable to connect to database")
		return err
	}
	defer c.DB.Close()
	var ignoreValue = ""
	if ignore {
		ignoreValue = "IGNORE"
	}
	sqlRequest := "INSERT " + ignoreValue + " INTO " + table + " " + col + " VALUES " + val
	request, err := c.DB.Prepare(sqlRequest)
	if err != nil {
		c.LoggerError("Unable to prepare insert request %s", sqlRequest)
		return err
	}

	_, err = request.Exec()
	if err != nil {
		c.LoggerError("Unable to execure insert request")
		return err
	}

	return nil
}

func (c *MariaDBConfiguration) Request(requestString string) (err error) {
	err = c.connectMariaDb()
	if err != nil {
		c.LoggerError("Unable to connect to database")
		return err
	}
	defer c.DB.Close()

	request, err := c.DB.Prepare(requestString)
	if err != nil {
		c.LoggerError("Unable to prepare request with string %s", requestString)
		return err
	}

	_, err = request.Exec()
	if err != nil {
		c.LoggerError("Unable to execure request")
		return err
	}

	return nil
}

func (c *MariaDBConfiguration) Select(requestString string) (response SelectResponse, err error) {
	err = c.connectMariaDb()
	if err != nil {
		c.LoggerError("Unable to connect to database")
		return response, err
	}
	defer c.DB.Close()

	// Execute the query
	rows, err := c.DB.Query(requestString)
	if err != nil {
		c.LoggerError("Failed to run Query : %+v", err)
		return response, err
	}

	// Get column names
	response.Columns, err = rows.Columns()
	if err != nil {
		c.LoggerError("Failed to get columns : %+v", err)
		return response, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(response.Columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	response.Seperator = c.Seperator
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			c.LoggerError("Failed to get row : %+v", err)
			return response, err
		}
		var rowBuilder strings.Builder
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if i != 0 {
				_, _ = fmt.Fprintf(&rowBuilder, "%s", c.Seperator)
			}
			if col == nil {
				_, _ = fmt.Fprintf(&rowBuilder, "%s", "")
			} else {
				_, _ = fmt.Fprintf(&rowBuilder, "%s", string(col))
			}
		}
		response.Rows = append(response.Rows, rowBuilder.String())
	}
	if err = rows.Err(); err != nil {
		c.LoggerError("Failed to run Query : %+v", err)
		return response, err
	}

	return response, nil
}

func SearchInTable(c *MariaDBConfiguration) (data interface{}, err error) {
	err = c.connectMariaDb()
	if err != nil {
		c.LoggerError("Unable to connect to database")
		return data, err
	}
	defer c.DB.Close()
	request := "SELECT * FROM " + c.Table + " WHERE " + c.WhereClause
	c.LoggerInfo("Sending request to database %s", request)
	resp, err := c.Select(request)
	if err != nil {
		c.LoggerError("Unable to launch select request : %v, %s", err, request)
		return data, err
	}
	c.LoggerInfo("Extracting data from response")
	params := new(ExtractDataOptions)
	params.Rows = resp.Rows
	params.Cols = resp.Columns
	params.Seperator = resp.Seperator
	params.Debug = false
	params.RemoveEndSpace = true
	params.RemoveStartSpace = true
	params.RemoveDoubleSpaces = true
	err = ExtractDataFromRowToStructure(c.DataType, *params)
	if err != nil {
		c.LoggerError("Unable to deserialize response : %v", err)
		return data, err
	}

	c.LoggerInfo("Extracting ended of data from response")
	return c.DataType, nil
}
