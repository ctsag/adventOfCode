package schematic

import (
	"adventOfCode/common/fileops"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type schematicValue struct {
	num string
	row int
	col int
}

func CalculateTotal(path string, reader fileops.ReadableFile) (int, error) {
	schematic, err := extractSchematic(path, reader)
	if err != nil {
		return -1, err
	}

	values := extractValues(schematic)
	total := 0

	for _, value := range values {
		if isAdjacentToSymbols(value, schematic) {
			num, _ := strconv.Atoi(value.num)
			total += num
			fmt.Printf("added %s on row %d col %d\n", value.num, value.row+1, value.col+1)
		} else {
			fmt.Printf("%s on row %d col %d is false\n", value.num, value.row+1, value.col+1)
		}
	}

	return total, nil
}

func extractSchematic(path string, reader fileops.ReadableFile) ([][]byte, error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	var schematic [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		schematic = append(schematic, line)
	}

	return schematic, nil
}

func extractValues(schematic [][]byte) []schematicValue {
	var values []schematicValue

	for rowIdx, row := range schematic {
		var num string
		var rowFound, colFound int

		for colIdx, col := range row {
			isDigit := unicode.IsDigit(rune(col))
			isNewNum := num == "" && isDigit
			isRowEnding := colIdx == len(row)-1
			hasNumEnded := num != "" && (!isDigit || isRowEnding)

			if isDigit {
				num += string(col)
			}

			if isNewNum {
				rowFound = rowIdx
				colFound = colIdx
			}

			if hasNumEnded {
				values = append(values, schematicValue{num, rowFound, colFound})
				num = ""
			}
		}
	}

	return values
}

func isAdjacentToSymbols(value schematicValue, schematic [][]byte) bool {
	const symbols = "#$%&*+-/=@"

	rowStart := max(value.row-1, 0)
	rowEnd := min(value.row+1, len(schematic)-1)

	colStart := max(value.col-1, 0)
	colEnd := min(value.col+len(value.num), len(schematic[0])-1)

	found := false
	for i := rowStart; i <= rowEnd; i++ {
		if !found {
			for j := colStart; j <= colEnd; j++ {
				char := string(schematic[i][j])
				if strings.Contains(symbols, char) {
					found = true
					break
				}
			}
		} else {
			break
		}
	}

	return found
}
