package errors

import (
	. "github.com/yohanc3/link-vault/config"
)

//Error struct
type Error struct {
	UserMessage string
	LogMessage string
}

// Implement the error interface for Error
func (e *Error) Error() string {
	return e.LogMessage
}

// Predefined error instances
var (
	InvalidCommandError = &Error{
		LogMessage: "invalid command",
		UserMessage: "The command issued is invalid. Please check the available commands and try again!",
	}
	InvalidTagsError = &Error{
		LogMessage: "tags were not found",
		UserMessage: "No tags were recognized:( \nExample: \n> " + BOT_PREFIX + "find cinema tech sports",
	}
	InvalidUrlError = &Error{
		LogMessage: "invalid url",
		UserMessage: "Your URL seems invalid...",
	}
	MissingTagsError = &Error{
		LogMessage: "missing tags",
		UserMessage: "You are missing then tags! \nExample: \n > "+BOT_PREFIX+"save https://example.com/ cinema tech soccer",
	}
	MissingUrlError = &Error{
		LogMessage: "missing url",
		UserMessage: "You are missing the url! \nExample: \n "+BOT_PREFIX+"save https://example.com/ cinema tech soccer",
	}
	//Avoid using LogMessage when using GenericError
	GenericError = &Error{
		LogMessage: "Something went wrong...Try again later:(",
		UserMessage: "Something went wrong",
	}
)

//Create new personalized error
func NewError(userMessage, logMessage string) error {
	return &Error{
			UserMessage: userMessage,
			LogMessage:  logMessage,
	}
}

