package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	pw "github.com/pschlump/pw3"
)

func main() {
	Pw := pw.NewParseWords()
	// Pw.SetOptions("C", false, false)  -- these are the default opsions, so skip

	scanner := bufio.NewScanner(os.Stdin)
	line_no := 0
	for scanner.Scan() {
		line_no++
		line := scanner.Text()

		Pw.SetLine(line)
		words := Pw.GetWords()

		fmt.Printf("Words len=%d: %s\n", len(words), words)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
