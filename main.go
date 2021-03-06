package main

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	dDir = "/data"
	rDir = "/result"
)

var (
	results          []string
	resultsExtracted []dataExtract
)

type dataExtract struct {
	date     string
	nickname string
	data     string
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Get Work Directory
	wd, err := os.Getwd()
	checkError(err)

	// Read Directory First
	log.Println("Directory:", wd+dDir+"/donate")
	files, err := ioutil.ReadDir(wd + dDir + "/donate")
	checkError(err)

	flen := len(files)
	log.Println("Files Count:", flen)

	if flen <= 0 {
		checkError(errors.New("no such files"))
	}

	// Read all files
	for _, f := range files {
		f, err := ioutil.ReadFile(wd + dDir + "/donate/" + f.Name())
		checkError(err)
		sf := strings.Split(string(f), "\n")
		results = append(results, sf...)
	}

	log.Println("Result Count:", len(results))

	// Extract to struct
	err = extract(results)
	checkError(err)
	log.Println("Data Extracted")

	// Convert to csv
	err = convertToCsv(resultsExtracted)
	checkError(err)
	log.Println("Csv Successfully Created")
}

func extract(d []string) error {
	if len(d) == 0 {
		return errors.New("data can not be null")
	}

	// loop data to be extracted
	for _, ld := range d {
		if len(ld) == 0 {
			continue
		}

		extLd := strings.Split(ld, " ")
		ext := dataExtract{
			date:     strings.Join(extLd[0:2], " "),
			nickname: extLd[3],
			data:     strings.Join(extLd[5:], " "),
		}

		resultsExtracted = append(resultsExtracted, ext)
	}

	return nil
}

func convertToCsv(d []dataExtract) error {
	wd, err := os.Getwd()
	checkError(err)

	// Create csv result file
	file, err := os.Create(wd + rDir + "/result.csv")
	checkError(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to csv
	h := []string{"Date", "Nickname", "Value"}
	err = writer.Write(h)
	checkError(err)

	// Store data to csv
	for _, value := range d {
		dt := []string{value.date, value.nickname, value.data}
		err := writer.Write(dt)
		checkError(err)
	}
	return nil
}
