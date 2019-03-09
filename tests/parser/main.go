package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var verboseFlag = flag.Bool("v", false, "verbose output")
var sleepFLag = flag.Int64("s", 1, "sleep when EOF in milliseconds")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Printf("usage: %s <file>\n", os.Args[0])
		return
	}

	file, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(0, io.SeekEnd)
	parser := csv.NewReader(file)
	parser.FieldsPerRecord -= 1
	parser.TrimLeadingSpace = true
	parser.ReuseRecord = true
	parser.Comma = ' '
	parser.LazyQuotes = true

	for {
		record, err := parser.Read()
		if err == io.EOF {
			if *verboseFlag {
				log.Println("EOF")
			}
			time.Sleep(time.Millisecond * time.Duration(*sleepFLag))
			continue
		}
		if *verboseFlag {
			log.Println(record)
		}
	}
}
