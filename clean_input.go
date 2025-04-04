package main

import (
	"slices"
	"strings"
)

func checkForProfane(b string) string {
	splited := strings.Split(b, " ")
	final := make([]string, 0)
	forbidden := []string{"kerfuffle", "sharbert", "fornax"}
	for _, s := range splited {
		l := strings.ToLower(s)
		isForbiddenWord := slices.Contains(forbidden, l)
		if isForbiddenWord {
			final = append(final, "****")
		} else {
			final = append(final, s)
		}
	}
	return strings.Join(final, " ")
}
