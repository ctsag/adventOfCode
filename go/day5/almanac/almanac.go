package almanac

import (
	"adventOfCode/common/fileops"
	"bufio"
	"strconv"
	"strings"
)

type rangeMap struct {
	destinationRangeStart int
	sourceRangeStart      int
	sourceRangeLength     int
}

type almanac struct {
	seeds        []int
	soils        []rangeMap
	fertilizers  []rangeMap
	waters       []rangeMap
	lights       []rangeMap
	temperatures []rangeMap
	humidities   []rangeMap
	locations    []rangeMap
}

type seed struct {
	id          int
	soil        int
	fertilizer  int
	water       int
	light       int
	temperature int
	humidity    int
	location    int
}

func DetermineLowestLocation(path string, reader fileops.ReadableFile) (location int, errorMsg error) {
	alm, err := extractAlmanac(path, reader)
	if err != nil {
		return -1, err
	}

	lowest := -1

	for _, seedId := range alm.seeds {
		s := extractSeed(alm, seedId)

		if lowest == -1 {
			lowest = s.location
		} else {
			lowest = min(lowest, s.location)
		}
	}

	return lowest, nil
}

func extractAlmanac(path string, reader fileops.ReadableFile) (almanac, error) {
	file, err := fileops.OpenFile(path, reader)
	if err != nil {
		return almanac{}, err
	}

	defer func() {
		_ = fileops.CloseFile(file)
	}()

	var alm almanac
	var currentRange *[]rangeMap

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		split := strings.Split(line, ":")

		if len(split) == 2 {
			switch split[0] {
			case "seeds":
				alm.seeds = toIntArray(split[1])
			case "seed-to-soil map":
				currentRange = &alm.soils
			case "soil-to-fertilizer map":
				currentRange = &alm.fertilizers
			case "fertilizer-to-water map":
				currentRange = &alm.waters
			case "water-to-light map":
				currentRange = &alm.lights
			case "light-to-temperature map":
				currentRange = &alm.temperatures
			case "temperature-to-humidity map":
				currentRange = &alm.humidities
			case "humidity-to-location map":
				currentRange = &alm.locations
			}

			continue
		}

		rangeValues := toIntArray(line)

		if currentRange != nil {
			*currentRange = append(*currentRange, rangeMap{rangeValues[0], rangeValues[1], rangeValues[2]})
		}
	}

	return alm, nil
}

func toIntArray(line string) []int {
	fields := strings.Fields(line)
	var intArray []int

	for i := 0; i < len(fields); i++ {
		intValue, _ := strconv.Atoi(fields[i])
		intArray = append(intArray, intValue)
	}

	return intArray
}

func extractSeed(alm almanac, id int) seed {
	s := seed{}

	s.id = id
	s.soil = nextValue(alm.soils, s.id)
	s.fertilizer = nextValue(alm.fertilizers, s.soil)
	s.water = nextValue(alm.waters, s.fertilizer)
	s.light = nextValue(alm.lights, s.water)
	s.temperature = nextValue(alm.temperatures, s.light)
	s.humidity = nextValue(alm.humidities, s.temperature)
	s.location = nextValue(alm.locations, s.humidity)

	return s
}

func nextValue(rmaps []rangeMap, id int) int {

	for i := 0; i < len(rmaps); i++ {
		destinationRangeEnd := rmaps[i].destinationRangeStart + rmaps[i].sourceRangeLength
		offset := id - rmaps[i].sourceRangeStart
		value := offset + rmaps[i].destinationRangeStart

		if offset < 0 {
			continue
		}

		if value >= rmaps[i].destinationRangeStart && value <= destinationRangeEnd {
			return value
		}
	}

	return id
}
