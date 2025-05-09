package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// replaceTextInIdentifiedLine searches for lines containing `lineIdentifier`. In such lines,
// it replaces the first occurrence of `textToReplace` with `replacementText`.
// This function is kept internal to the cmd package.
func replaceTextInIdentifiedLine(filePath string, lineIdentifier string, textToReplace string, replacementText string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	identifierFound := false
	replacementMade := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, lineIdentifier) {
			identifierFound = true
			modifiedLine := strings.Replace(line, textToReplace, replacementText, 1) // Replace only the first occurrence
			if modifiedLine != line {
				lines = append(lines, modifiedLine)
				replacementMade = true
			} else {
				// Line contained identifier, but no replacement of textToReplace occurred
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	if !identifierFound {
		return fmt.Errorf("identifier '%s' not found in %s", lineIdentifier, filePath)
	}
	if !replacementMade {
		return fmt.Errorf("text '%s' not found for replacement (or no change needed) in lines containing '%s' in %s", textToReplace, lineIdentifier, filePath)
	}

	// Write changes back to the file
	file.Seek(0, 0)  // Go to the beginning of the file
	file.Truncate(0) // Clear the file content
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

// Execute defines and parses flags, then calls the file editing logic.
func Execute() {
	dshFilePath := flag.String("dshfile", "", "Path to the .dsh file to edit (required).")
	lineIdent := flag.String("lineid", "", "Substring to identify the target line(s) for editing (required).")
	oldText := flag.String("oldtext", "", "The exact text to be replaced within the identified line(s) (required).")
	newText := flag.String("newtext", "", "The new text to substitute for oldtext (can be empty to delete oldtext).")

	flag.Parse()

	// Perform validation directly after parsing flags
	if *dshFilePath == "" {
		fmt.Fprintln(os.Stderr, "Error: -dshfile path cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}
	if *lineIdent == "" {
		fmt.Fprintln(os.Stderr, "Error: -lineid cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}
	if *oldText == "" {
		fmt.Fprintln(os.Stderr, "Error: -oldtext cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

	err := replaceTextInIdentifiedLine(*dshFilePath, *lineIdent, *oldText, *newText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error editing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully processed %s: replaced '%s' with '%s' in lines containing '%s'.\n", *dshFilePath, *oldText, *newText, *lineIdent)
}
