package gsvalidation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError is the struct returned by some functions of this package to indicate that
// one or more errors were found in a validation process.
type HttpSuggestionError struct {
	Status uint
	Message string
}

// ResponseType is used to inform in a function output params of this package what
// type of response is being returned to the caller.
type ResponseType int

const (
	OK_RESPONSE ResponseType = 0
	ERR_RESPONSE ResponseType = 1
	DECODING_ERROR ResponseType = 2
	VALIDATION_ERROR ResponseType = 3
)

// validate is designed to be thread-safe and used as a singleton instance. It caches
// information about every struct and validations, in essence only parsing your validation
// tags once per struct type. Using multiple instances neglects the benefit of caching.
var validate *validator.Validate

// Initialization
func init() {
	validate = validator.New()
}

// DecodeJSONResponseBody wraps ioutil.ReadAll and json.Unmarshal functions to extract an
// http.Response's body and unmarshal it into either an ok response or an error response
// from an external API.
func DecodeJSONResponseBody(r *http.Response, dstOk interface{}, dstErr interface{}) (respType ResponseType, err error) {

	// Reads http.Response's body as a slice of bytes
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return DECODING_ERROR, fmt.Errorf("error retrieving response: %s", err.Error())
	}

	// Tries to unmarshal the slice of bytes as the destination Ok response. If no error is 
	// found, it means the decoding was successful.
	err = json.Unmarshal(data, dstOk)
	if err == nil {

		// Validates the response before returning it to the caller.
		err = Struct(dstOk)
		if err != nil {
			return VALIDATION_ERROR, err
		}

		// Returns OK_RESPONSE notifying the caller the ok response is available in dstOk.
		return OK_RESPONSE, nil
	}

	// Save the error message got so far.
	errMsg := err.Error()

	// If the caller only wanted to parse OK response.
	if dstErr == nil {
		return DECODING_ERROR, fmt.Errorf("error unmarshalling ok response: %s", errMsg)
	}

	// Tries to unmarshal the slice of bytes as the destination Err response. If no error is
	// found, it means the decoding was successful.
	err = json.Unmarshal(data, dstErr)
	if err == nil {

		// Validates the response before returning it to the caller.
		err = Struct(dstErr)
		if err != nil {
			return VALIDATION_ERROR, err
		}

		// Returns ERR_RESPONSE notifying the caller the error response is available in dstErr.
		return ERR_RESPONSE, nil
	}

	// Both tries to unmarshal the data failed. Notify the caller.
	return DECODING_ERROR, fmt.Errorf("error retrieving response. dstOk unmarshal error: %s; dstErr unmarshal error: %s", errMsg, err.Error())
}

// DecodeJSONRequestBody wraps json.Decoder.Decode() function, adding validation and error
// checking. If an error is found, an HttpSuggestionError is returned, suggesting that the
// caller returns this error to the client.
func DecodeJSONRequestBody(r *http.Request, dst interface{}) *HttpSuggestionError {

	// Checks if the content-type header is set. If it's not, request is discarted, even
	// if the body contains a well-formed JSON.
	if r.Header.Get("Content-Type") == "" {
		return &HttpSuggestionError{
			Status: http.StatusUnsupportedMediaType,
			Message: "Content-Type header is missing",
		}
	}

	// Checks if the content-type header has the value application/json.
	// Note that the check works even if the client includes additional charset or boundary
	// information in the header.
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return &HttpSuggestionError{
			Status: http.StatusUnsupportedMediaType,
			Message: "Content-Type header is not application/json",
		}
	}

	// Enforces a maximum read of 1 MB from the request body. Additionaly, A request 
	// body larger than that will now result in Decode() returning an EOF error.
	// Then, we can detect if the request was larger than 1 MB by checking if N <= 0.
	maxBytes := int64(1 << 20) // 1 MB
	limitedReader := &io.LimitedReader{R: r.Body, N: maxBytes}

	// Setup the decoder.
	dec := json.NewDecoder(limitedReader)

	err := dec.Decode(&dst)
	if err != nil {
		
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message which
		// interpolates the location of the problem to make it easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &HttpSuggestionError{
				Status: http.StatusBadRequest,
				Message: msg,
			}

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error for
		// syntax errors in the JSON.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &HttpSuggestionError{
				Status: http.StatusBadRequest,
				Message: msg,
			}

		// Catch any type errors, like trying to assign a string in the JSON request body
		// to a int field in our struct. We can interpolate the relevant field name and
		// position into the error message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)", 
				unmarshalTypeError.Field, unmarshalTypeError.Offset,
			)
			return &HttpSuggestionError{
				Status: http.StatusBadRequest,
				Message: msg,
			}

		// Catch the error caused by the request body being too large. The error in this case
		// would be EOF so this validation should be prior to the EOF catch.
		case limitedReader.N <= 0:
			return &HttpSuggestionError{
				Status: http.StatusRequestEntityTooLarge,
				Message: "Request body must not be larger than 1 MB",
			}

		// An io.EOF error is returned by Decode() if the request body is empty.
		case errors.Is(err, io.EOF):
			return &HttpSuggestionError{
				Status: http.StatusBadRequest,
				Message: "Request body must not be empty",
			}

		// Otherwise default to sending a 500 server error.
		default:
			return &HttpSuggestionError{
				Status: http.StatusInternalServerError,
				Message: "Unrecognized error when decoding request JSON body",
			}
		}
	}
	// Finally, call RequestStruct to run all validations on it.
	return RequestStruct(dst)
}

// RequestStruct wraps validator.Validate.Struct function. Validates a struct's exposed
// fields, and automatically validates nested structs, unless otherwise specified.

// If an error is found, an HttpSuggestionError is returned, suggesting that the caller 
// returns this error to the client.
func RequestStruct(v interface{}) *HttpSuggestionError {
	err := validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return &HttpSuggestionError{
				Status: http.StatusInternalServerError,
				Message: "Could not validate request body. Please contact support team",
			}
		}
		return &HttpSuggestionError{
			Status: http.StatusBadRequest,
			Message: FlatErrors(err.(validator.ValidationErrors)).Error(),
		}
	}
	return nil
}

// Struct wraps validator.Validate.Struct function. Validates a struct's exposed fields
// and automatically validates nested structs, unless otherwise specified.

// Returns error found, if any, containing all messages joined by a semicolon.
func Struct(v interface{}) error {
	err := validate.Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid validation found while validating client response")
		}
		return FlatErrors(err.(validator.ValidationErrors))
	}
	return nil
}

// Var wraps validator.Validate.Var function. Validates a single field against the
// requirements given by the tag.

// Returns error found, if any, containing all messages joined by a semicolon.
func Var(v interface{}, tag string) error {
	err := validate.Var(v, tag)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid validation found while validating field")
		}
		return FlatErrors(err.(validator.ValidationErrors))
	}
	return nil
}

// FlatErrors receives a slice of validator.FieldError (alias: ValidationErrors) and
// returns an error built with all retrieved messages joined by a semicolon.
func FlatErrors(e validator.ValidationErrors) error {
	msg := ""
	for _, err := range e {
		errMsg := err.Error()
		if msg == "" {
			msg = errMsg
		} else {
			msg = fmt.Sprintf("%s; %s", msg, errMsg)
		}
	}
	return errors.New(msg)
}