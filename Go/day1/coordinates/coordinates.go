package coordinates

import (
	"bufio"
	"day1/fileops"
	"strconv"
	"unicode"
)

func CalculateTotal(path string, reader fileops.ReadableFile) (int, error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return -1, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	scanner := bufio.NewScanner(file)

	calibrationValuesTotal := 0

	for scanner.Scan() {
		line := scanner.Text()
		calibrationValuesTotal += combineFirstAndLastDigit(line)
	}

	return calibrationValuesTotal, nil
}

func combineFirstAndLastDigit(line string) int {
	firstDigit, lastDigit := 0, 0

	for i := 0; i < len(line); i++ {
		firstDigit = isDigitOtherwiseZero(line[i])

		if firstDigit > 0 {
			break
		}
	}

	for i := len(line) - 1; i >= 0; i-- {
		lastDigit = isDigitOtherwiseZero(line[i])

		if lastDigit > 0 {
			break
		}
	}

	return firstDigit*10 + lastDigit
}

func isDigitOtherwiseZero(input uint8) int {
	char := rune(input)

	if unicode.IsDigit(char) {
		digit, _ := strconv.Atoi(string(char))
		return digit
	}

	return 0
}
