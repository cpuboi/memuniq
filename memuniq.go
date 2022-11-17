package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	bloom "github.com/EverythingMe/inbloom/go/inbloom" // Vendor directory overrides github path.
	"github.com/cpuboi/memuniq/tools"
)

// Load the good old input arguments
func getArgs() (string, bool, bool, bool, int, float64, bool) {
	// This should handle alternative cache locations..
	var filterFilePath = flag.String("f", os.ExpandEnv("$HOME/.cache/bloomfilter.bin"), "Location of bloomfilter file")
	var showInfo = flag.Bool("i", false, "Show information about processed lines")
	var beVerbose = flag.Bool("v", false, "Show verbose information")
	var newBloomfilter = flag.Bool("n", false, "Create a new filter and delete the old")
	var sizeBloomfilter = flag.Int("s", 1_000_000, "Size of bloomfilter before major collissions occur")
	var percentageHitrate = flag.Float64("p", 0.001, "Approximate error rate percentage, default 0.001%")
	var abortIfFileMissing = flag.Bool("a", false, "Abort process if the filter file does not exist")
	flag.Parse()
	return *filterFilePath, *showInfo, *beVerbose, *newBloomfilter, *sizeBloomfilter, *percentageHitrate, *abortIfFileMissing
}

func main() {
	logger := log.New(os.Stderr, "", 0)
	scanner := bufio.NewScanner(os.Stdin) // Scanner returns lines one by one
	i := 0
	hit := 0

	filterFilePath, showInfo, beVerbose, newBloomfilter, sizeBloomfilter, percentageHitrate, abortIfFileMissing := getArgs()
	var bf, err = bloom.NewFilter(sizeBloomfilter, percentageHitrate) // Create initial bloomfilter
	if err != nil {
		panic(err)
	}

	// Check that bloom filter path is writeable
	if !tools.CheckFilterPath(filterFilePath) {
		logger.Fatal("- Could not write to", filterFilePath)
	}

	// Test the abort function
	if abortIfFileMissing {
		if _, err := os.Stat(filterFilePath); err != nil {
			logger.Fatal("- Filter does not exist:", filterFilePath)
		}
	}
	// If user did not select create new filter, try and find old filter and load it, otherwise continue with the newly generated filter.
	if !newBloomfilter {
		if _, err := os.Stat(string(filterFilePath)); err == nil { // If file exists, load it.
			file, err := os.ReadFile(filterFilePath)
			if err != nil {
				panic(err)
			}
			bf, err = bloom.Unmarshal(file)
			if err != nil {
				panic(err)
			}
			if beVerbose {
				logger.Println("- Loading bloomfilter file", string(filterFilePath))
			}
		} else {
			if beVerbose {
				logger.Println("- Did not find bloomfilter", string(filterFilePath), "creating new")
			}
		}
	}

	// Read lines from stdin
	if beVerbose {
		fmt.Printf("- Printing lines not seen by bloomfilter\n\n")
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
		logger.Println("\n- Statistics \n\t- Duplicates:", hit, "\n\t- Unique:", i) // Show statistics
	}

	saveFile, err := os.Create(filterFilePath)
	if err != nil {
		logger.Println("- Could not create bloomfilter file", filterFilePath)
		panic(err)
	}
	defer saveFile.Close() // Make sure to close file.

	bytesWritten, err := saveFile.Write(bf.Marshal())
	if err != nil {
		logger.Println("- Could not save bloomfilter file", filterFilePath)
		panic(err)
	}
	if beVerbose {
		logger.Printf("\n- Saved %d bytes in file %s\n", bytesWritten, filterFilePath)
	}
}
