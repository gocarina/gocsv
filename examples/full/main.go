package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type Test struct {
	TestName   string  `csv:"Test Name"`
	TestDate   string  `csv:"Test Date"`
	Batch      string  `csv:"Batch"`
	SA         int     `csv:"Student Appeared"`
	TMode      string  `csv:"Test Mode"`
	Physics    int     `csv:"Physics (180)"`
	Chemistry  int     `csv:"Chemistry (180)"`
	Biology    int     `csv:"Biology (360)"`
	Total      int     `csv:"Total Marks (720)"`
	Percentage float64 `csv:"Percentage (%)"`
	Percentile float64 `csv:"Percentile"`
	TestRank   string  `csv:"Test Rank"`
	AIR        string  `csv:"AIR"`
}

func main() {
	// Data representing the Student Performance Report
	tests := []Test{
		{
			TestName:   "MAJOR TEST (M36604460)",
			TestDate:   "18 Dec-24",
			Batch:      "MEL6B",
			SA:         5546,
			TMode:      "OFFLINE",
			Physics:    150,
			Chemistry:  149,
			Biology:    338,
			Total:      637,
			Percentage: 88.47,
			Percentile: 93.23,
			TestRank:   "-",
			AIR:        "-",
		},
		{
			TestName:   "MAJOR TEST (M36013940)",
			TestDate:   "13 Dec-24",
			Batch:      "MEL6B",
			SA:         5938,
			TMode:      "OFFLINE",
			Physics:    144,
			Chemistry:  152,
			Biology:    331,
			Total:      627,
			Percentage: 87.08,
			Percentile: 91.23,
			TestRank:   "-",
			AIR:        "-",
		},
		{
			TestName:   "MAJOR TEST (M35516516)",
			TestDate:   "08 Dec-24",
			Batch:      "MEL6B",
			SA:         6455,
			TMode:      "OFFLINE",
			Physics:    126,
			Chemistry:  152,
			Biology:    355,
			Total:      633,
			Percentage: 87.92,
			Percentile: 96.75,
			TestRank:   "-",
			AIR:        "-",
		},
		{
			TestName:   "MAJOR TEST (M35049023)",
			TestDate:   "03 Dec-24",
			Batch:      "MEL6B",
			SA:         6178,
			TMode:      "OFFLINE",
			Physics:    165,
			Chemistry:  170,
			Biology:    355,
			Total:      690,
			Percentage: 95.83,
			Percentile: 99.04,
			TestRank:   "-",
			AIR:        "-",
		},
	}

	// Create the CSV file
	file, err := os.Create("student_performance_report.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Write to CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{"Test Name", "Test Date", "Batch", "Student Appeared", "Test Mode", "Physics (180)", "Chemistry (180)", "Biology (360)", "Total Marks (720)", "Percentage (%)", "Percentile", "Test Rank", "AIR"}
	if err := writer.Write(headers); err != nil {
		log.Fatalf("Failed to write headers: %v", err)
	}

	// Write test data
	for _, test := range tests {
		row := []string{
			test.TestName,
			test.TestDate,
			test.Batch,
			formatInt(test.SA),
			test.TMode,
			formatInt(test.Physics),
			formatInt(test.Chemistry),
			formatInt(test.Biology),
			formatInt(test.Total),
			formatFloat(test.Percentage),
			formatFloat(test.Percentile),
			test.TestRank,
			test.AIR,
		}
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row: %v", err)
		}
	}

	log.Println("CSV successfully created: student_performance_report.csv")
}

// Helper function to format integers
func formatInt(value int) string {
	return fmt.Sprintf("%d", value)
}

// Helper function to format floats
func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
