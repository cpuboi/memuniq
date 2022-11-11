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
	var b64FilePath = flag.String("f", "~/.cache/bloomfilter.b64", "Location of bloomfilter file")
	//var b64FilePath string = "/dev/shm/bloomfilter.b64"
	flag.Parse()
	var bf, err = bloom.NewFilter(10, 0.1)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(string(*b64FilePath)); err == nil { // If file exists, load it.
		//data, err := os.ReadFile(b64FilePath)
		file, err := os.ReadFile(*b64FilePath)
		if err != nil {
			panic(err)
		}
		//bf, err = bloom.Unmarshal(data)
		bf, err = bloom.Unmarshal(file)
		if err != nil {
			panic(err)
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
	logger.Println("Bloomfilter duplicates:", hit, "Unique:", i) // Show statistics
	//TODO: if can write to path, write:
	os.WriteFile(*b64FilePath, bf.Marshal(), 0644)
}
