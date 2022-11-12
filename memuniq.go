package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	bloom "github.com/EverythingMe/inbloom/go/inbloom"
)

func main() {

	logger := log.New(os.Stderr, "", 0)
	scanner := bufio.NewScanner(os.Stdin) // Scanner returns lines one by one
	i := 0
	hit := 0

	// This section should handle several cache locations automatically
	var b64FilePath = flag.String("f", os.ExpandEnv("$HOME/.cache/bloomfilter.b64"), "Location of bloomfilter file")
	var showStatistics = flag.Bool("s", false, "Show statistics of processed lines")
	var beVerbose = flag.Bool("v", false, "Show verbose information")
	//var b64FilePath string = "/dev/shm/bloomfilter.b64"
	flag.Parse()
	var bf, err = bloom.NewFilter(10, 0.1)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(string(*b64FilePath)); err == nil { // If file exists, load it.
		file, err := os.ReadFile(*b64FilePath)
		if err != nil {
			panic(err)
		}
		bf, err = bloom.Unmarshal(file)
		if err != nil {
			panic(err)
		}
		if *beVerbose {
			fmt.Println("Loading bloomfilter file", string(*b64FilePath))
		}
	} else { // Create new bloomfilter
		bf, err = bloom.NewFilter(1_000_000, 0.001)
		if err != nil {
			panic(err)
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
	if *showStatistics {
		logger.Println("Bloomfilter duplicates:", hit, "Unique:", i) // Show statistics
	}

	//TODO: if can write to path, write:
	// Can onle write to file if it already exists

	//os.WriteFile(*b64FilePath, bf.Marshal(), 0644)

	saveFile, err := os.Create(*b64FilePath)
	if err != nil {
		logger.Println("Could not create bloomfilter file", *b64FilePath)
		panic(err)
	}
	defer saveFile.Close() // Make sure to close file.

	bytesWritten, err := saveFile.Write(bf.Marshal())
	if err != nil {
		logger.Println("Could not save bloomfilter file", *b64FilePath)
		panic(err)
	}
	if *beVerbose {
		logger.Printf("Saved %d bytes in file %s\n", bytesWritten, *b64FilePath)
	}
}
