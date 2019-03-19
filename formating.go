package formating

import (
	"fmt"
	"strconv"
)

type TableData [][]string

func New() (t *TableData) {
	t = &TableData{}
	return t
}

func (table *TableData) Add(value ...string) {
	*table = append(*table, value)
}

func (table *TableData) AddArray(value []string) {
	*table = append(*table, value)
}

func (table TableData) ArrayToFormatedText() string {
	table, nbrOfColumns := table.createSquareTable()
	columnSizes := table.calculateSizeOfColumns(nbrOfColumns)
	totalWidth := calculateTotalWidth(columnSizes, nbrOfColumns)
	text := "\n"
	for row := 0; row < len(table); row++ {
		for column := 0; column < len(table[row]); column++ {
			text += fmt.Sprintf("%-"+strconv.Itoa(columnSizes[column])+"s", table[row][column]) + " "

		}
		text += "\n"
		if row == 0 {
			for i := 0; i < totalWidth; i++ {
				text += "-"
			}
			text += "\n"
		}
	}
	return text
}

func calculateTotalWidth(columnSizes []int, amount int) int {
	totalColumnsWidth := 0
	for _, col := range columnSizes {
		totalColumnsWidth += col
	}
	return totalColumnsWidth + amount - 1
}

func (table TableData) createSquareTable() ([][]string, int) {
	maxColumns := 0
	for _, row := range table {
		if len(row) > maxColumns {
			maxColumns = len(row)
		}
	}
	for i := 0; i < len(table); i++ {
		if len(table[i]) < maxColumns {
			for j := 1; j <= maxColumns; j++ {
				table[i] = append(table[i], "")
			}
		}
	}
	return table, maxColumns
}

func (table TableData) calculateSizeOfColumns(nbrOfColumns int) (columnSize []int) {
	for column := 0; column < nbrOfColumns; column++ {
		maxSizeOfColumn := 0
		for row := 0; row < len(table); row++ {
			rowSize := len(table[row][column])
			if maxSizeOfColumn < rowSize {
				maxSizeOfColumn = rowSize
			}
		}
		columnSize = append(columnSize, maxSizeOfColumn)
	}
	return columnSize
}
