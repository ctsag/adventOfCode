package coordinates

import (
	"adventOfCode/common/fileops"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	line = replaceWordsWithDigits(line)

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

func replaceWordsWithDigits(line string) string {
	nums := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for i := 0; i < len(line); i++ {
		for _, num := range nums {
			if strings.HasPrefix(line[i:], num) {
				numIdx, _ := arrayIndexOf(nums, num)
				replaced := nonDestructiveReplace(num, numIdx+1)
				line = strings.Replace(line, num, replaced, 1)
				line = replaceWordsWithDigits(line)
			}
		}
	}

	return line
}

func arrayIndexOf(arr []string, search string) (int, error) {
	for idx, val := range arr {
		if val == search {
			return idx, nil
		}
	}

	return -1, errors.New("index not found")
}

func nonDestructiveReplace(val string, num int) string {
	replacedVal := fmt.Sprintf("%s%d%s", val[0:1], num, val[1:])

	return replacedVal
}

func isDigitOtherwiseZero(input uint8) int {
	char := rune(input)

	if unicode.IsDigit(char) {
		digit, _ := strconv.Atoi(string(char))
		return digit
	}

	return 0
}
