package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Anomalies struct {
	Year  int
	Month string
	Mean  float64
}

func main() {
	local := true
	var nasagta_csv string
	var data interface{}
	var err error

	if local {
		nasagta_csv = "C:/Users/mhoff/Downloads/GLB.Ts+dSST.csv"
		data, err = os.Open(nasagta_csv)
		if err != nil {
			panic(err)
		}

	} else {
		nasagta_csv = "https://data.giss.nasa.gov/gistemp/tabledata_v4/GLB.Ts+dSST.csv"
		resp, err := http.Get(nasagta_csv)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		data = resp.Body
	}
	fmt.Println("Reading CSV file from: ", nasagta_csv)

	// Create a CSV reader
	readr := csv.NewReader(data.(interface{ Read([]byte) (int, error) }))
	readr.FieldsPerRecord = -1 // disable FilesPerRecord test

	// Read all CSV records
	records, err := readr.ReadAll()
	if err != nil {
		panic(err)
	}

	// Initialize anomalies struct
	var anomalies []Anomalies

	for i := 2; i < len(records); i++ {
		year, err := strconv.Atoi(records[i][0]) // Convert year to int
		if err != nil {
			panic(err)
		}

		for j := 1; j < 13; j++ {
			mean, err := strconv.ParseFloat(records[i][j], 64) // Convert mean to float
			if err != nil {
				break
			}

			anomalies = append(anomalies, Anomalies{
				Year:  year,
				Month: records[1][j], // Month name from the first row
				Mean:  mean,
			})
		}
	}

	fmt.Println("Anomalies data:")
	fmt.Println("Year, Month, Mean")
	for _, anomaly := range anomalies {
		fmt.Printf("%d, %s, %.2f\n", anomaly.Year, anomaly.Month, anomaly.Mean)
	}

}
