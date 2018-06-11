package jsonlint

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	// Create new parser object
	parser := argparse.NewParser("jsonlint", "A simple JSON linting and formatting tool")

	// Create file flag
	inputFile := parser.File("f", "file", os.O_RDWR, 0644, &argparse.Options{Required: true, Help: "Path to JSON file"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		fmt.Print(parser.Usage(err))
	} else {
		raw, err := ioutil.ReadAll(inputFile)
		check(err)

		jsonstring := string(raw)
		var js json.RawMessage
		err = json.Unmarshal([]byte(jsonstring), &js)
		check(err)

		// Truncate the input file and seek to the beginning
		inputFile.Truncate(0)
		inputFile.Seek(0, 0)

		//encoder := json.NewEncoder(os.Stdout)
		encoder := json.NewEncoder(inputFile)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "    ")

		err = encoder.Encode(js)
		check(err)
	}

}
