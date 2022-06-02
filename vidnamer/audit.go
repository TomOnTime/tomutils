package main

import (
	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
	"golang.org/x/exp/maps"
)

//  auditKeywords returns a list of invalid keywords in use.
func auditKeywords(inventory []filminventory.Film, permittedKeywords []string) []string {
	var result []string

	for _, film := range inventory {
		existing := film.Keywords
		invalid := difference(existing, permittedKeywords)
		result = append(result, invalid...)

	}

	return removeDuplicateValues(result)
}

//  auditTags returns a list of invalid keywords in use.
func auditTags(inventory []filminventory.Film, permittedTags []string) []string {
	var result []string

	for _, film := range inventory {
		existing := maps.Keys(film.Tags)
		invalid := difference(existing, permittedTags)
		result = append(result, invalid...)
	}

	return removeDuplicateValues(result)
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
