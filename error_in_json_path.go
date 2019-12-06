package main

import (
	"log"
)

func error_in_json_path(json_validationerror JSONValidationError, document string) {

	log.Fatal("\"", json_validationerror.errortype, "\"-error in json at object\n  ", json_validationerror.field, ": ", json_validationerror.description) // throw error

	///TODO improve error handling
}
