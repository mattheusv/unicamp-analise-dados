package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 12, 31, 23, 59, 59, 59, time.UTC)

	sales := generateSales(startDate, endDate)
	cpu := generateCPUFromSales(sales)
	disk := generateDiskFromSales(sales)
	memory := generateMemoryFromSales(sales)

	if err := createCSV("sales.csv", sales); err != nil {
		log.Fatal(err)
	}

	if err := createCSV("cpu.csv", cpu); err != nil {
		log.Fatal(err)
	}

	if err := createCSV("disk.csv", disk); err != nil {
		log.Fatal(err)
	}

	if err := createCSV("memory.csv", memory); err != nil {
		log.Fatal(err)
	}

	// for i, ts := range sales {
	// 	fmt.Fprint(io.Discard, ts.ts, ts.value)
	// 	fmt.Println(ts.ts.Format("2006-01-02 03:04:05"), ts.value, cpu[i].value)
	// }
}

func createCSV(csvName string, timeseries []timeserie) error {
	buf := bytes.NewBufferString("timestamp,value\n")

	for _, ts := range timeseries {
		fmt.Fprintf(buf, "%s,%d\n", ts.ts.Format("2006-01-02 15:04:05"), ts.value)
	}

	f, err := os.Create(csvName)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

type timeserie struct {
	ts    time.Time
	value int
}

func generateCPUFromSales(sales []timeserie) []timeserie {
	timeseries := make([]timeserie, 0)

	cost := 10
	for _, ts := range sales {
		value := (ts.value * cost) / 100

		timeseries = append(timeseries, timeserie{
			ts:    ts.ts,
			value: value,
		})
	}

	return timeseries
}

func generateDiskFromSales(sales []timeserie) []timeserie {
	timeseries := make([]timeserie, 0)

	cost := 4
	for _, ts := range sales {
		value := (ts.value / cost) / 100

		timeseries = append(timeseries, timeserie{
			ts:    ts.ts,
			value: value,
		})
	}

	return timeseries
}

func generateMemoryFromSales(sales []timeserie) []timeserie {
	timeseries := make([]timeserie, 0)

	cost := 2
	var lastTS time.Time
	for _, ts := range sales {
		if ts.ts.Month() <= 10 && (ts.ts.Month() > lastTS.Month()) {
			cost = cost * 2
		}

		value := (ts.value * cost) / 100

		timeseries = append(timeseries, timeserie{
			ts:    ts.ts,
			value: value,
		})

		lastTS = ts.ts
	}

	return timeseries
}

func generateSales(start, end time.Time) []timeserie {
	return generateSeries(start, end, 900, 3)
}

func generateSeries(start, end time.Time, max, freq int) []timeserie {
	timeseries := make([]timeserie, 0)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d

		// n is a number that grows exponentially with the month.
		n := (int(date.Month()) + 1) * 2

		// Iterate over minutes based on the given freq. If freq is 3 it means
		// that the series generated will be between 3 mintutes.
		for i := 0; i < 1440; i++ {
			if i != 0 {
				d2 := date.Add(time.Minute * time.Duration(freq))
				if (d2.Year() > date.Year()) || (d2.Month() > date.Month()) || (d2.Day() > date.Day()) {
					// Date after sum is in the future, break the loop and go
					// to the next date.
					break
				}
				date = d2
			}

			value := n*(i+1) + time.Now().Nanosecond()
			// Make sure that the value don't execed the maximum value and also
			// that is not negative.
			value = int(math.Abs(float64(value)))
			value = truncate(value, max)

			timeseries = append(timeseries, timeserie{
				ts:    date,
				value: value,
			})
		}

	}

	return timeseries
}

func truncate(v, max int) int {
	for {
		if v <= max {
			return v
		}
		v = v / 2
	}
}
