package processing

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

const outputUrl = "/Users/peterbishop/Development/upgraded-telegram/processing/output/"
const inputUrl = "/Users/peterbishop/Development/upgraded-telegram/processing/input/"

func OpenFile(fileName string) ([][]string, error) {
	file, err := os.Open(inputUrl + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	return records, nil
}

func GenerateTripData() bool {

	if _, err := os.Stat(outputUrl + "tasks.go"); err == nil {
		fmt.Println("Trip File already exists, skipping generation.")
		return true
	}

	records, err := OpenFile("trips.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	var trips []Trip
	for i, row := range records {
		if i == 0 {
			continue
		}

		var directionID int
		fmt.Sscanf(row[4], "%d", &directionID)
		blockID := strings.TrimSpace(row[5])

		trips = append(trips, Trip{
			RouteID:      row[0],
			ServiceID:    row[1],
			TripID:       row[2],
			TripHeadsign: row[3],
			DirectionID:  directionID,
			BlockID:      blockID,
			ShapeID:      row[6],
		})
	}

	outputFile := fmt.Sprintf(outputUrl + "trips.go")
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating Go file:", err)
		return false
	}
	defer file.Close()

	fmt.Fprintln(file, "package output")
	fmt.Fprintln(file, "import \"probable-system/main.go/processing\"")
	fmt.Fprintln(file, "var Trips = []processing.Trip{")
	for _, trip := range trips {
		fmt.Fprintf(file, "\t{RouteID: \"%s\", ServiceID: \"%s\", TripID: \"%s\", TripHeadsign: \"%s\", DirectionID: %d, BlockID: \"%s\", ShapeID: \"%s\"},\n",
			trip.RouteID, trip.ServiceID, trip.TripID, trip.TripHeadsign, trip.DirectionID, trip.BlockID, trip.ShapeID)
	}
	fmt.Fprintln(file, "}")
	fmt.Println("Go file successfully saved to", outputFile)
	return true
}
