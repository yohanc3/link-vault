package util

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"wisdom-hoard/config"
	customErr "wisdom-hoard/error"
)

const PREFIX string = config.BOT_PREFIX

// parseHelpCommand
// parseHelpCommand function parses the input string and returns a slice of strings
// that come after the "{'PREFIX'}help" command.
// Returns []string for tags, error, and if error is not critical
func ParseFindCommand(input string) ([]string, error, bool) {
	// Split the input string into fields (words)
	words := strings.Fields(input)

	err := errors.New(customErr.InvalidTagsError.ExternalMessage + "\nExample: \n> " + PREFIX + "find cinema tech sports")
	
	for i, word := range words {
		if word == PREFIX+"find" {
			// Return the slice of words that come after "+help"
			tags := words[i+1:]

			if len(tags) == 0 {
				return []string{}, err, true
			}

			return tags, nil, true
		}
	}
	// if "+find" is not found, return an empty slice, error and true to signal this error shouldn't stop the program
	return []string{}, errors.New(customErr.InvalidCommandError.ExternalMessage), true

}

//Returns the url, a slice of tags, error if any, and true if it is not a critical error or false if it is a critical error
func ParseSaveCommand(input string) (string, []string, error, bool){

	words := strings.Fields(input)

	if words[0] != PREFIX+"save"{
		return "", []string{}, errors.New(customErr.InvalidCommandError.ExternalMessage), true
	}

	if len(words) == 1 {
		PREFIX := config.BOT_PREFIX
		MissingTags := customErr.NewCustomError(customErr.InvalidCommand, "", "You are missing the url and tags! \nExample: \n > "+PREFIX+"save https://example.com/ cinema tech soccer")
		return "", []string{}, errors.New(MissingTags.ExternalMessage), true
	}
	
	url := words[1]
	
	if !isValidURL(url){
		InvalidUrlError := errors.New(customErr.InvalidUrlError.ExternalMessage)
		return "", []string{}, InvalidUrlError, true
	}

	//if there's only the command and the url with no additional tags
	if len(words) == 2 {
		return "", []string{}, errors.New(customErr.MissingTagsError.ExternalMessage), true
	}

	return url, words[2:], nil, true

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