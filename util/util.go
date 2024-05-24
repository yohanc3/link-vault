package util

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	. "github.com/yohanc3/link-vault/config"
	. "github.com/yohanc3/link-vault/error"
)

// parseFindCommand
// parseFindCommand function parses the input string and returns a slice of strings
// that come after the "{'PREFIX'}help" command.
// Returns []string for tags and error
func ParseFindCommand(input string) ([]string, error) {
	// Split the input string into fields (words)
	words := strings.Fields(input)

	for i, word := range words {
		if word == BOT_PREFIX+"find" {
			// Return the slice of words that come after "+help"
			tags := words[i+1:]
			
			if len(tags) == 0 {
				fmt.Println((*InvalidTagsError).LogMessage)
				return []string{}, InvalidTagsError
			}

			return tags, nil
		}
	}
	// if "+find" is not found, return an empty slice, error and true to signal this error shouldn't stop the program
	return []string{}, InvalidCommandError

}

//Returns the url, a slice of tags, and an error if any
func ParseSaveCommand(input string) (string, []string, error){

	words := strings.Fields(input)

	if words[0] != BOT_PREFIX+"save"{
		return "", []string{}, InvalidCommandError
	}

	if len(words) == 1 {
		return "", []string{}, MissingUrlError
	}
	
	url := words[1]
	
	if !isValidURL(url){
		return "", []string{}, InvalidUrlError
	}

	//if there's only the command and the url with no additional tags
	if len(words) == 2 {
		return "", []string{}, MissingTagsError
	}

	return url, words[2:], nil

}

// isValidURL checks if the given string is a valid URL.
func isValidURL(str string) bool {
	parsedURL, err := url.Parse(str)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func FormatArrayToString(arr []string) string {
	// Determine the maximum width of the elements
	maxWidth := 0
	for _, str := range arr {
		if len(str) > maxWidth {
			maxWidth = len(str)
		}
	}

	var builder strings.Builder

	// Loop through the array and format each row
	for i := 0; i < len(arr); i += 4 {
		end := i + 4
		if end > len(arr) {
			end = len(arr)
		}
		row := arr[i:end]

		for j, item := range row {
			if j > 0 {
				builder.WriteString(" ") // Add space between columns
			}
			builder.WriteString(fmt.Sprintf("%-*s", maxWidth, item))
		}
		builder.WriteString("\n") // Add a newline at the end of each row
	}

	return builder.String()
}

// Function to merge two slices and remove duplicates
func MergeSlices(slice1, slice2 []string) []string {
	// Create a map to track unique elements
	uniqueElements := make(map[string]bool)
	// Create a slice to hold the result
	result := []string{}

	// Add elements from the first slice to the map and result slice
	for _, v := range slice1 {
			if !uniqueElements[v] {
					uniqueElements[v] = true
					result = append(result, v)
			}
	}

	// Add elements from the second slice to the map and result slice
	for _, v := range slice2 {
			if !uniqueElements[v] {
					uniqueElements[v] = true
					result = append(result, v)
			}
	}

	return result
}

func CurrentTime() string {
	now := time.Now()
	formattedTime := now.Format("01/02/2006 15:04")
	return formattedTime
}
