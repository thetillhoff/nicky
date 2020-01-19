package main

import (
	"log"
)

func errorInJSONPath(jsonValidationerror JSONValidationError, document string) {

	log.Fatal("\"", jsonValidationerror.errortype, "\"-error in json at object\n  ", jsonValidationerror.field, ": ", jsonValidationerror.description) // throw error

	///TODO improve error handling
}
