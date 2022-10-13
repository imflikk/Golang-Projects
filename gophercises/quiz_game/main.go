package main

import (
	"fmt"
	"flag"
	"encoding/csv"
	"os"
	"log"
	"strconv"
)


func main() {

	csvPtr := flag.String("csv", "problems.csv", "CSV file containing quiz problems")

	flag.Parse()

	var guess int
	var counter int

	filePath := *csvPtr
	correct := 0

	f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()

	if err != nil {
	    log.Fatal("Unable to parse file as CSV for " + filePath, err)
	}

    for i := 0; i < len(records); i++ {

    	answer, _ := strconv.Atoi(records[i][1])
    	counter = i + 1
		fmt.Printf("Problem #%d: %+v = ", counter, records[i][0])
		fmt.Scanf("%d", &guess)

		if guess == answer {
			correct++
		}
    }

    fmt.Printf("You guessed %d out of %d correctly.\n", correct, counter)
}

