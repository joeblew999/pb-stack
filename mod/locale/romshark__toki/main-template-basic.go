package main

import (
	"fmt"

	"TOKI_NAME/tokibundle"

	"golang.org/x/text/language"
)

func main() {
	// Get a localized reader for British English.
	// Toki will automatically select the most appropriate translation catalog available.
	reader, _ := tokibundle.Match(language.BritishEnglish)

	// This comment describes the text below and is included in the translator context.
	fmt.Println(reader.String(`{"Framework"} is powerful yet easy to use!`, "Toki"))

	// This comment describes the text below and is included in the translator context.
	fmt.Println(reader.String(`{"Framework"} is powerful yet easy to use also `, "Toki"))
}
