package apierrors

import "fmt"

// error codes
const (
	// default
	ERR_NOT_DEFINED 	     ErrorCode = -1
	// business errors
	USER_NOT_FOUND  	     ErrorCode = 1001
	EXTERNAL_API_ERROR       ErrorCode = 1002
	// client errors
	OPERATION_NOT_DEFINED    ErrorCode = 2001
	INVALID_ARGUMENT 	     ErrorCode = 2002
	// server errors
	INTERNAL_SERVER_ERROR    ErrorCode = 5001
	IO_FILE_ERROR 		     ErrorCode = 5002
	JSON_PARSING_ERROR 	     ErrorCode = 5003
	READER_ERROR             ErrorCode = 5004
	CLIENT_NOT_DEFINED       ErrorCode = 5005
	DTO_MAPPING_ERROR        ErrorCode = 5010
	// HTTP errors
	HTTP_CONNECTION_ERROR    ErrorCode = 6001
	RESPONSE_UNMARSHAL_ERROR ErrorCode = 6002
)

var errors = map[ErrorCode]ErrorMessage {
	// default
	ERR_NOT_DEFINED: {
		Label: "ERR_NOT_DEFINED", 
		Message: "Unrecognized error encountered. Contact support team"},
	// business errors
	USER_NOT_FOUND: {
		Label: "USER_NOT_FOUND", 
		Message: "User could not be found"},
	EXTERNAL_API_ERROR: {
		Label: "EXTERNAL_API_ERROR",
		Message: "External API returned an error",
	},
	// client errors
	OPERATION_NOT_DEFINED: {
		Label: "OPERATION_NOT_DEFINED",
		Message: "Invoked path/operation is not defined",
	},
	INVALID_ARGUMENT: {
		Label: "INVALID_ARGUMENT",
		Message: "Client specified invalid request parameter",
	},
	// server errors
	INTERNAL_SERVER_ERROR: {
		Label: "INTERNAL_SERVER_ERROR",
		Message: "Internal error found. Please contact support team",
	},
	IO_FILE_ERROR: {
		Label: "IO_FILE_ERROR",
		Message: "Input/Output error while reading a file",
	},
	JSON_PARSING_ERROR: {
		Label: "JSON_PARSING_ERROR",
		Message: "Error while parsing a JSON content",
	},
	READER_ERROR: {
		Label: "READER_ERROR",
		Message: "Error while reading Body from Request",
	},
	CLIENT_NOT_DEFINED: {
		Label: "CLIENT_NOT_DEFINED",
		Message: "External API client not defined",
	},
	DTO_MAPPING_ERROR: {
		Label: "DTO_MAPPING_ERROR",
		Message: "Failed mapping model to DTO",
	},
	// HTTP errors
	HTTP_CONNECTION_ERROR: {
		Label: "HTTP_CONNECTION_ERROR",
		Message: "Connection error while attempting http connection",
	},
	RESPONSE_UNMARSHAL_ERROR: {
		Label: "RESPONSE_UNMARSHAL_ERROR",
		Message: "Could not unmarshal response body received from external api",
	},
}

type ErrorCode int

type ErrorMessage struct {
	Label 	string
	Message string
}

// the serializable error structure
type Error struct {
	Code 	ErrorCode `json:"code"`
	Label 	string 	  `json:"label"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return e.ToString()
}

func (e *Error) ToString() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Label, e.Message)
}

func New(ec ErrorCode) *Error {
	errMsg := errors[ERR_NOT_DEFINED]
	if em, ok := errors[ec]; ok {
		errMsg = em
	}
	return &Error{
		Code: ec,
		Label: errMsg.Label,
		Message: errMsg.Message,
	}
}

func NewWithMsg(ec ErrorCode, m string) *Error {
	errMsg := errors[ERR_NOT_DEFINED]
	if em, ok := errors[ec]; ok {
		errMsg = em
	}
	return &Error{
		Code: ec,
		Label: errMsg.Label,
		Message: m,
	}
}