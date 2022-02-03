package main

import "sort"

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
