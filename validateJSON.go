package main

import (
	"log"

	"github.com/xeipuuv/gojsonschema"
)

func validateJSON(document string, schema string) []JSONValidationError {

	var schemaLoader gojsonschema.JSONLoader            // initialize schemaLoader before if, so it can be used afterwards
	if schema[0:4] == "file" || schema[0:4] == "http" { // if schemalocation is provided via uri path
		schemaLoader = gojsonschema.NewReferenceLoader(schema) // use ReferenceLoader
	} else if schema[:1] == "/" || schema[:2] == "./" || schema[:3] == "../" { // if schema is provided via unix path
		schemaLoader = gojsonschema.NewStringLoader(loadFile(schema)) //
	} else { // else schema is provided as string
		schemaLoader = gojsonschema.NewStringLoader(schema) // use StringLoader
	}

	documentLoader := gojsonschema.NewStringLoader(document)           // load document from string
	result, err := gojsonschema.Validate(schemaLoader, documentLoader) // validate document against schema
	if err != nil {                                                    // if error is thrown
		debuglog("document", document)
		debuglog("schema or -path", schema)
		log.Fatal("error: ", err) // throw error
	}
	if result.Valid() { // if validation was successful
		//debuglog("", "document is valid JSON")
		return nil // return 'true' (no error means success)
	} else { // if validation was unsuccesful
		debuglog("", "document is not valid JSON. See errors below")
		var jsonValidationerrors []JSONValidationError
		for _, validationerror := range result.Errors() { // for each error
			jsonValidationerrors = append(jsonValidationerrors, JSONValidationError{errortype: validationerror.Type(), description: validationerror.Description(), field: validationerror.Field()}) // add error to local error list. The instance is created in this line, thus the length
		}
		return jsonValidationerrors // return 'false' (errors mean failure)
	}
}

// For better error handling they are passed around in this type.
type JSONValidationError struct {
	errortype   string      // contains a type from https://github.com/xeipuuv/gojsonschema#working-with-errors
	description interface{} // contains the error itself
	field       string      // contains the location where the error happened
}
