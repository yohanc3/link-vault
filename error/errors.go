package errors

type CustomError struct {
	ErrorType 			 ErrorType
	InternalMessage  string
	ExternalMessage  string
}

// Implement the error interface for CustomError
func (e CustomError) Error() string {
	return e.ExternalMessage
}

type ErrorType string

// Predefined error types
const (
	InvalidCommand ErrorType = "InvalidCommand"
	InvalidTags ErrorType = "InvalidTags"
	InvalidUrl ErrorType = "InvalidURL"
	MissingTags ErrorType = "MissingTags"
)

// Predefined error instances
var (
	InvalidCommandError = CustomError{
		ErrorType:       InvalidCommand,
		InternalMessage: "The command is not recognized",
		ExternalMessage: "The command issued is invalid. Please check the available commands and try again!",
	}
	InvalidTagsError = CustomError{
		ErrorType: InvalidTags,
		InternalMessage: "Tags were not given or are unrecognized",
		ExternalMessage: "No tags were recognized:(",
	}
	InvalidUrlError = CustomError{
		ErrorType: InvalidUrl,
		InternalMessage: "The provided URL is not valid",
		ExternalMessage: "Your URL seems invalid...",
	}
	MissingTagsError = CustomError{
		ErrorType: MissingTags,
		InternalMessage: "Tags were not provided",
		ExternalMessage: "You also need to pass in tags that describe the type of content for this post",
	}
)

func NewCustomError(errorType ErrorType, internalMessage string, externalMessage string) CustomError {
	return CustomError{
		ErrorType:			 errorType,
		InternalMessage: internalMessage,
		ExternalMessage: externalMessage,
	}
}