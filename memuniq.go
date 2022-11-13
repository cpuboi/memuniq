package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	bloom "github.com/EverythingMe/inbloom/go/inbloom"
)

// Load the good old input arguments
func getArgs() (string, bool, bool, bool, int, float64) {

	// This should handle alternative cache locations..
	var b64FilePath = flag.String("f", os.ExpandEnv("$HOME/.cache/bloomfilter.b64"), "Location of bloomfilter file")
	var showInfo = flag.Bool("i", false, "Show statistics of processed lines")
	var beVerbose = flag.Bool("v", false, "Show verbose information")
	var newBloomfilter = flag.Bool("n", false, "Create a new filter and delete the old.")
	var sizeBloomfilter = flag.Int("s", 1_000_000, "Size of bloomfilter before major collissions occur, default: 1_000_000")
	var percentageHitrate = flag.Float64("p", 0.001, "Approximate error rate percentage, default 0.001%")
	flag.Parse()
	return *b64FilePath, *showInfo, *beVerbose, *newBloomfilter, *sizeBloomfilter, *percentageHitrate
}

func main() {

	logger := log.New(os.Stderr, "", 0)
	scanner := bufio.NewScanner(os.Stdin) // Scanner returns lines one by one
	i := 0
	hit := 0

	b64FilePath, showInfo, beVerbose, newBloomfilter, sizeBloomfilter, percentageHitrate := getArgs()
	var bf, err = bloom.NewFilter(sizeBloomfilter, percentageHitrate) // Create initial bloomfilter
	if err != nil {
		panic(err)
	}

	// If user did not select new filter, try and find old filter and load it, otherwise continue with the newly generated filter.
	if !newBloomfilter {
		if _, err := os.Stat(string(b64FilePath)); err == nil { // If file exists, load it.
			file, err := os.ReadFile(b64FilePath)
			if err != nil {
				panic(err)
			}
			bf, err = bloom.Unmarshal(file)
			if err != nil {
				panic(err)
			}
			if beVerbose {
				logger.Println("Loading bloomfilter file", string(b64FilePath))
			}
		} else {
			if beVerbose {
				logger.Println("Did not find bloomfilter", string(b64FilePath), "creating new")
			}
		}
	}

	for scanner.Scan() { // For every line in input:
		txtLine := scanner.Text()

		if bf.Contains(txtLine) {
			hit++
		} else {
			fmt.Println(txtLine)
			i++
			bf.Add(txtLine)
		}
	}

	if showInfo {
		logger.Println("Bloomfilter duplicates:", hit, "Unique:", i) // Show statistics
	}

	//TODO: if can write to path, write:
	// Can onle write to file if it already exists
	saveFile, err := os.Create(b64FilePath)
	if err != nil {
		logger.Println("Could not create bloomfilter file", b64FilePath)
		panic(err)
	}
	defer saveFile.Close() // Make sure to close file.

	bytesWritten, err := saveFile.Write(bf.Marshal())
	if err != nil {
		logger.Println("Could not save bloomfilter file", b64FilePath)
		panic(err)
	}
	if beVerbose {
		logger.Printf("Saved %d bytes in file %s\n", bytesWritten, b64FilePath)
	}
}
