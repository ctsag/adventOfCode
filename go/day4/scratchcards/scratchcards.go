package scratchcards

import (
	"adventOfCode/common/fileops"
	"bufio"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type scratchcard struct {
	id      int
	winning []int
	drawn   []int
}

func CalculateTotals(path string, reader fileops.ReadableFile) (score int, count int, errorMsg error) {
	scratchcards, err := extractScratchcards(path, reader)
	if err != nil {
		return -1, -1, err
	}

	score = 0
	for _, card := range scratchcards {
		score += determineScore(card)
	}

	bonusScratchcards := addBonusCards(scratchcards)
	count = len(bonusScratchcards)

	return score, count, nil
}

func extractScratchcards(path string, reader fileops.ReadableFile) ([]scratchcard, error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	var scratchcards []scratchcard

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		scratchcards = append(scratchcards, cardFromLine(line))
	}

	return scratchcards, nil
}

func cardFromLine(line string) scratchcard {
	regEx := regexp.MustCompile(`Card\s+(\d+):\s+([\d\s]+?)\s*\|\s*([\d\s]+)`)
	match := regEx.FindStringSubmatch(line)

	id, _ := strconv.Atoi(match[1])
	winning := toIntArray(match[2])
	drawn := toIntArray(match[3])

	return scratchcard{id, winning, drawn}
}

func toIntArray(line string) []int {
	parts := strings.Fields(line)

	nums := make([]int, len(parts))
	for i, num := range parts {
		nums[i], _ = strconv.Atoi(num)
	}

	return nums
}

func determineScore(card scratchcard) int {
	score := 0

	for _, winning := range card.winning {
		if slices.Contains(card.drawn, winning) {
			if score == 0 {
				score = 1
			} else {
				score *= 2
			}
		}
	}

	return score
}

func addBonusCards(cards []scratchcard) []scratchcard {
	for i := 0; i < len(cards); i++ {
		baseIdx := max(cards[i].id-1, 0)

		score := countMatches(cards[baseIdx])

		startIdx := baseIdx + 1
		endIdx := min(baseIdx+1+score, len(cards))

		cards = append(cards, cards[startIdx:endIdx]...)
	}

	return cards
}

func countMatches(card scratchcard) int {
	score := 0

	for _, winning := range card.winning {
		if slices.Contains(card.drawn, winning) {
			score++
		}
	}

	return score
}
