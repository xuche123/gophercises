package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvPtr := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	intPtr := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvPtr)

	if err != nil {
		fmt.Println("Failed to open provided CSV file")
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fileReader := csv.NewReader(file)
	questions, err := fileReader.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	exitSignal := make(chan bool)

	go func() {
		timer := time.NewTimer(time.Duration(*intPtr) * time.Second)
		<-timer.C
		exitSignal <- true
	}()

	var ans string
	score := 0

	go func() {
		for _, element := range questions {

			question := element[0]
			fmt.Printf("%s = ", question)
			_, err := fmt.Scanf("%s\n", &ans)
			if err != nil {
				fmt.Println("Answer should be a valid integer.")
			}

			if ans == element[1] {
				score++
			}
		}
	}()

	<-exitSignal
	fmt.Printf("\nYou scored %d out of %d\n", score, len(questions))
	os.Exit(0)
}
