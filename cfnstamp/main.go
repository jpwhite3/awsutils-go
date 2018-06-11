package main

import (
	"log"
	"os"
	"github.com/akamensky/argparse"
	"github.com/Jeffail/gabs"
	"fmt"
	"io/ioutil"
	"io"
	"crypto/md5"
)

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {

	// Create new parser object
	parser := argparse.NewParser("cfnstamp", "A tool for stamping Cloudformation templates with meta data")

	// Create file flag
	inputFile := parser.File("f", "file", os.O_RDWR, 0644, &argparse.Options{Required: true, Help: "Path to Cloudformation template"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		fmt.Print(parser.Usage(err))
	} else {
		raw, err := ioutil.ReadAll(inputFile)
		check(err)

		// Read existing file contents and parse as JSON
		jsonstring := string(raw)
		jsonParsed, err := gabs.ParseJSON([]byte(jsonstring))
		check(err)

		// Store then remove the metadata section of the template if it exists
		metadataPath := "Metadata"
		var version float64
		var md5sum string
		if jsonParsed.Exists(metadataPath) {
			version, _ = jsonParsed.Path("Metadata.version").Data().(float64)
			md5sum, _ = jsonParsed.Path("Metadata.md5").Data().(string)
			jsonParsed.Delete(metadataPath)
		}

		// Calculate the MD5 of the remaining template
		hash := md5.New()
		io.WriteString(hash, jsonParsed.String())

		// Calculate new version and MD5sum
		newVersion := version + 1
		newSum := fmt.Sprintf("%x", hash.Sum(nil))

		// If MD5 is different, update template
		if md5sum != newSum {
			// Add the updated metadata back into the original template
			jsonParsed, _ = gabs.ParseJSON([]byte(jsonstring))
			jsonParsed.Set(newSum, "Metadata", "md5")
			jsonParsed.Set(newVersion, "Metadata", "version")

			// Truncate the input file and seek to the beginning
			inputFile.Truncate(0)
			inputFile.Seek(0, 0)

			// Encode and write back into the original file
			inputFile.WriteString(jsonParsed.StringIndent("", "    "))
		}

	}

}
