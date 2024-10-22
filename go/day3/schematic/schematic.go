package schematic

import (
	"adventOfCode/common/fileops"
	"bufio"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type schematicValue struct {
	num string
	row int
	col int
}

type gear struct {
	row int
	col int
}

func CalculateTotals(path string, reader fileops.ReadableFile) (schematicValueTotal int, gearRatioTotal int, errorMsg error) {
	schematic, err := extractSchematic(path, reader)
	if err != nil {
		return -1, -1, err
	}

	values, gears := extractTokens(schematic)
	valueTotal := 0
	ratioTotal := 0

	for _, value := range values {
		if isAdjacentToSymbols(value, schematic) {
			num, _ := strconv.Atoi(value.num)
			valueTotal += num
		}
	}

	for _, gear := range gears {
		ratioTotal += getGearRatioOrZero(gear, schematic)
	}

	return valueTotal, ratioTotal, nil
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

func extractTokens(schematic [][]byte) ([]schematicValue, []gear) {
	var values []schematicValue
	var gears []gear

	for rowIdx, row := range schematic {
		var num string
		var rowFound, colFound int

		for colIdx, col := range row {
			isGear := col == '*'
			isDigit := unicode.IsDigit(rune(col))
			isNewNum := num == "" && isDigit
			isRowEnding := colIdx == len(row)-1
			hasNumEnded := num != "" && (!isDigit || isRowEnding)

			if isGear {
				gears = append(gears, gear{rowIdx, colIdx})
			}

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

	return values, gears
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

func getGearRatioOrZero(coords gear, schematic [][]byte) int {
	rowStart := max(coords.row-1, 0)
	rowEnd := min(coords.row+1, len(schematic)-1)

	colStart := max(coords.col-1, 0)
	colEnd := min(coords.col+1, len(schematic[0])-1)

	nums := make(map[schematicValue]bool)

	for i := rowStart; i <= rowEnd; i++ {
		for j := colStart; j <= colEnd; j++ {
			char := schematic[i][j]
			if unicode.IsDigit(rune(char)) {
				value, _ := getValueFromCoordinates(i, j, schematic)
				nums[value] = true
			}
		}
	}

	if len(nums) != 2 {
		return 0
	}

	ratio := 1
	for key := range nums {
		num, _ := strconv.Atoi(key.num)
		ratio *= num
	}

	return ratio
}

func getValueFromCoordinates(row int, col int, schematic [][]byte) (schematicValue, error) {
	var toTheLeft, toTheRight string
	var startingCol int
	var err error = nil

	for i := col; i >= 0; i-- {
		if !unicode.IsDigit(rune(schematic[row][i])) {
			break
		}
		toTheLeft = string(schematic[row][i:col])
		startingCol = i
	}

	for i := col; i < len(schematic[row]); i++ {
		if !unicode.IsDigit(rune(schematic[row][i])) {
			break
		}
		toTheRight = string(schematic[row][col : i+1])
	}

	value := toTheLeft + toTheRight

	if len(value) == 0 {
		err = errors.New("no value found at provided coordinates")
	}

	return schematicValue{value, row, startingCol}, err
}
