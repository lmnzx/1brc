package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const measurementsPath = "../measurements.txt"

type Measurement struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func main() {
	dataFile, err := os.Open(measurementsPath)
	if err != nil {
		panic(err)
	}

	defer dataFile.Close()

	measurements := make(map[string]*Measurement)

	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		rawString := fileScanner.Text()
		stationName, termperatureStr, found := strings.Cut(rawString, ";")
		if !found {
			continue
		}

		termperature, err := strconv.ParseFloat(termperatureStr, 32)
		if err != nil {
			panic(err)
		}

		measurement := measurements[stationName]
		if measurement == nil {
			measurements[stationName] = &Measurement{
				Min:   termperature,
				Max:   termperature,
				Sum:   termperature,
				Count: 1,
			}
		} else {
			measurement.Max = max(measurement.Max, termperature)
			measurement.Min = min(measurement.Min, termperature)
			measurement.Sum += termperature
			measurement.Count += 1
		}
	}

	printResults(measurements)
}

func printResults(results map[string]*Measurement) {
	// sort by station name
	stationNames := make([]string, 0, len(results))
	for stationName := range results {
		stationNames = append(stationNames, stationName)
	}

	slices.Sort(stationNames)

	fmt.Print("{")
	for idx, stationName := range stationNames {
		measurement := results[stationName]
		mean := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f", stationName, measurement.Min, mean, measurement.Max)
		if idx < len(stationNames)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")
}
