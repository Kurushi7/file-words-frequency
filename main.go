package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("data/GoLang_Text.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Problem closing the file")
		}
	}(file)

	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)

	var formattedLine string
	post := PairList{}
	finalList := PairList{}

	for scanner.Scan() {
		params := url.Values{}

		// replace dots and commas in the string with space
		formattedLine = strings.Replace(scanner.Text(), ".", " ", -1)
		formattedLine = strings.Replace(formattedLine, ",", " ", -1)

		params.Add("text", formattedLine)
		params.Add("fetchAll", "true")

		resp, err := http.PostForm("http://localhost:8085/frequency", params)
		if err != nil {
			log.Printf("Request Failed: %s", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		// Unmarshal result
		err = json.Unmarshal(body, &post)
		if err != nil {
			log.Printf("Reading body failed: %s", err)
			return
		}

		//merge the finalList to the new lines coming in
		finalList = mergeAndRecount(finalList, post)
		closeError := resp.Body.Close()
		if err != nil {
			fmt.Println("Error closing the scanning channel:" + closeError.Error())
		}
	}

	// if the scan of lines ever go wrong
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// the PairList will be not be sorted by value when its come out
	// so we have to do it manually
	sortedList := sortList(finalList)

	if len(sortedList) > 10 {
		fmt.Println(sortedList[:10])
	} else {
		fmt.Println(sortedList)
	}

}
