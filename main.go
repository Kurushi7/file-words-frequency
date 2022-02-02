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
	"sort"
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

		resp, err := http.PostForm("http://localhost:8080/", params)
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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sortedList := sortList(finalList)

	if len(sortedList) > 10 {
		fmt.Println(sortedList[:10])
	} else {
		fmt.Println(sortedList)
	}

}

/*
	Takes array of struct and compares with the final struct being built.
	Fills any word which is not in the final array
	Updates the occurrences param by adding the count from the passed array of struct to that of the final array
*/
func mergeAndRecount(finalList PairList, currentPost PairList) PairList {

	if len(finalList) == 0 {
		finalList = currentPost
	} else {
		i := len(finalList) + 1
		found := false
		for _, value1 := range currentPost {
			found = false
			for index, value2 := range finalList {
				if value1.Text == value2.Text {
					found = true
					finalList[index].Occurence = value2.Occurence + value1.Occurence
				}
				i++
			}
			if !found {
				keyPair := Pair{value1.Text, value1.Occurence}
				finalList = append(finalList, keyPair)
			}
		}
	}

	return finalList
}

func sortList(finalList PairList) PairList {

	sort.Slice(finalList, func(i, j int) bool {
		return finalList[i].Occurence > finalList[j].Occurence
	})
	return finalList
}
