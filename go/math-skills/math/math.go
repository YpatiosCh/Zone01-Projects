package math

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// readData reads numbers from a file, returning a slice of float64
func ReadData(filename string) ([]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var numbers []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number in file: %v", err)
		}
		numbers = append(numbers, num)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("no valid numbers found in file")
	}

	return numbers, nil
}

// calculateAverage computes the arithmetic mean
func calculateAverage(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

// calculateMedian computes the median value
func calculateMedian(numbers []float64) float64 {
	// Create a copy to avoid modifying the original slice
	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	sort.Float64s(sorted)

	length := len(sorted)
	if length%2 == 0 {
		return (sorted[length/2-1] + sorted[length/2]) / 2
	}
	return sorted[length/2]
}

// calculateVariance computes the variance
func calculateVariance(numbers []float64, average float64) float64 {
	sumSquaredDiff := 0.0
	for _, num := range numbers {
		diff := num - average
		sumSquaredDiff += diff * diff
	}
	return sumSquaredDiff / float64(len(numbers))
}

// calculateStdDeviation computes the standard deviation
func calculateStdDeviation(variance float64) float64 {
	return math.Sqrt(variance)
}

func GetCalculations(numbers []float64) {
	// Calculate statistics
	avg := calculateAverage(numbers)
	med := calculateMedian(numbers)
	variance := calculateVariance(numbers, avg)
	stdDev := calculateStdDeviation(variance)

	// Print results rounded to integers
	fmt.Printf("Average: %d\n", int(math.Round(avg)))
	fmt.Printf("Median: %d\n", int(math.Round(med)))
	fmt.Printf("Variance: %d\n", int(math.Round(variance)))
	fmt.Printf("Standard Deviation: %d\n", int(math.Round(stdDev)))
}
