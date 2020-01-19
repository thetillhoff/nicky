package main

/// the required imports
import (
	"encoding/json"
	"fmt"
	"log"
)

/// global variables
var (
	globalSourcepath         string // define global variable for source file path
	globalSchemapath         string // define global variable for own schema path
	globalOfficialschemapath string // define global variable for official schema path
	globalSchemaOnly         bool   // define global variable for working on schema only
	globalSourceOnly         bool   // define global variable for working on source only
	globalSchemaisjson       bool   // define global variable for disabling conversion and unfolding
	globalDisableUnfolding   bool   // define global variable for disabling unfolding feature
	globalDebug              bool   // initialize variable for debugging mode with false
)

/// main function of nicky
/// For more details on what this program does see ./readme.md.
func main() {

	//
	// Next comes the loading, conversion and verification of the document-file
	//

	globalSourcepath, globalSchemapath, globalOfficialschemapath, globalSchemaOnly, globalSourceOnly, globalSchemaisjson, globalDisableUnfolding, globalDebug = readCliArguments() // read and parse cli arguments and check for possible errors
	printCliArguments()                                                                                                                                                            // output of values on debug mode

	documentjson := ""                         // declaration is required out of if-statement, as it is needed for further processing
	documentyaml := ""                         // declaration is required out of if-statement for error tracing at the end
	if !globalSchemaOnly || globalSourceOnly { // if not only the schema should be processed OR if only the source should be processed, then also the conversion and verification should be made
		documentyaml = loadFile(globalSourcepath) // load document.yaml
		debuglog("content of "+globalSourcepath, documentyaml)

		documentjson = convertYamlDocumentToJSON(documentyaml) // convert and thus validate the document content (-> validation here means it is valid yaml)
		debuglog("content of jsonified "+globalSourcepath, documentjson)
	}

	if !globalSourceOnly { // if not only source (-> schema) is provided, schema and/or schema validation needs to be handled

		//
		// The yaml-formatted document is now converted to json, which means, it contains valid json after conversion.
		// Next comes the loading, conversion and verification of the schema-file
		//

		schemajson := ""       // declaration is required out of if-statement, as it is needed for further processing
		schemayaml := ""       // declaration is required out of if-statement for error tracing
		foldedSchemayaml := "" // declaration is required out of if-statement for error tracing

		if !globalSchemaisjson { // if schema is not provided in json -> schema is provided in yaml, then do the unfolding
			foldedSchemayaml = loadFile(globalSchemapath) // load schema.yaml
			debuglog("content of provided "+globalSchemapath, foldedSchemayaml)

			schemajson = convertYamlDocumentToJSON(foldedSchemayaml) // convert schema from yaml to json and validate it // not used any further if no validation error occurs

			//
			// The yaml-formatted schema was now converted to json, which means, it contains valid json after conversion.
			// Next there can be two options:
			//   The schema is provided in json, which means there is nothing much todo concerning conversion and unfolding of the features (or "schema extensions") introduced by this software.
			//   The schema is provided in yaml, which means conversion must happen and unfolding should happen if enabled.
			//

			if !globalDisableUnfolding { // if unfolding is not disabled
				schemayaml = unfoldSchema(foldedSchemayaml) // unfold own schemayaml-functions
				debuglog("unfolded "+globalSchemapath, schemayaml)
			}

			schemajson = convertYamlDocumentToJSON(schemayaml) // convert unfolded schema from yaml to json and validate it
			debuglog("content of jsonified "+globalSchemapath, schemajson)

		} else { // schema is provided in json -> don't do unfolding
			schemajson = loadFile(globalSchemapath) // load schema.json from file
			debuglog("content of "+globalSchemapath, schemajson)

			// making sure the provided json schema file contains valid json
			var validationvariable map[string]interface{}                  // not further used
			err := json.Unmarshal([]byte(schemajson), &validationvariable) // parse filecontent from json to map // map is not further used, but required for method call
			if err != nil {                                                // if error is thrown
				log.Fatal("error: ", err) // throw error
			}
		}

		//
		// The schema is now formatted in json.
		// Next comes either
		//   the verification whether the schema is a valid json schema or
		//   this is skipped and the next step starts
		//

		if globalOfficialschemapath != "" { // if officialschema is set

			jsonValidationerrors := validateJSON(schemajson, globalOfficialschemapath) // validate schema against officialschema

			if jsonValidationerrors != nil { // if validation is unsuccessful
				// Schema is not a valid json schema.
				// If schema was provided in json, display line number of error in json file
				// else (schema was provided in yaml) display line number of error in yaml file
				for _, jsonValidationerror := range jsonValidationerrors { // for each validationerror
					if globalSchemaisjson {
						// schema was provided in json
						errorInJSONPath(jsonValidationerror, schemajson)
					} else {
						// schema was provided in yaml
						errorInYAMLPath(jsonValidationerror, schemayaml)
					}
				}
			} else { // everything is fine
				debuglog("", "jsonified (unfolded) schema was successfully validated against provided metaschema at "+globalOfficialschemapath)
			}
		}

		//
		// The yaml-formatted schema was now unfolded and
		// again converted to json and
		// [optionally] it was also verified against the official json-meta-schema, which means, it contains a valid json-schema after conversion.
		// Next comes either
		//   the validation of the document against the custom schema or
		//   if only a schema is provided, output the (unfolded) jsonified schema
		//

		if !globalSchemaOnly { // if only schema is provided and processed, don't do the validation

			jsonValidationerrors := validateJSON(documentjson, schemajson) // validate document against schema

			if jsonValidationerrors != nil { // if validation is unsuccessful
				// document is not valid by schema definitions
				for _, jsonValidationerror := range jsonValidationerrors { // for each validationerror
					errorInYAMLPath(jsonValidationerror, documentyaml)
				}
			} else { // everything is fine
				debuglog("", globalSourcepath+" was successfully validated against provided "+globalSchemapath)
				fmt.Println("verification successful")
			}
		} else { // instead
			fmt.Println(schemajson) // just print the jsonified schema
		}
	} else { // instead, when if only source (-> no schema) is provided, schema and/or schema validation doesn't need to be handled

		//
		// The yaml-formatted document is now converted to json, which means, it contains valid json after conversion.
		// Next comes the creation of an example schema out of the provided sourcefile
		//

		generatedschemayaml := generateFoldedSchema(documentyaml) // generate example schema
		fmt.Println(generatedschemayaml)                          // just print the generated schema
	}
}
