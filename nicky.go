package main

/// the required imports
import (
	"encoding/json"
	"fmt"
	"log"
)

/// global variables
var (
	global_sourcepath         string // define global variable for source file path
	global_schemapath         string // define global variable for own schema path
	global_officialschemapath string // define global variable for official schema path
	global_schema_only        bool   // define global variable for working on schema only
	global_source_only        bool   // define global variable for working on source only
	global_schemaisjson       bool   // define global variable for disabling conversion and unfolding
	global_disable_unfolding  bool   // define global variable for disabling unfolding feature
	global_debug              bool   // initialize variable for debugging mode with false
)

/// main function of nicky
/// For more details on what this program does see ./readme.md.
func main() {

	//
	// Next comes the loading, conversion and verification of the document-file
	//

	global_sourcepath, global_schemapath, global_officialschemapath, global_schema_only, global_source_only, global_schemaisjson, global_disable_unfolding, global_debug = read_cli_arguments() // read and parse cli arguments and check for possible errors
	print_cli_arguments()                                                                                                                                                                       // output of values on debug mode

	documentjson := ""                             // declaration is required out of if-statement, as it is needed for further processing
	documentyaml := ""                             // declaration is required out of if-statement for error tracing at the end
	if !global_schema_only || global_source_only { // if not only the schema should be processed OR if only the source should be processed, then also the conversion and verification should be made
		documentyaml = load_file(global_sourcepath) // load document.yaml
		debuglog("content of "+global_sourcepath, documentyaml)

		documentjson = convert_yaml_document_to_json(documentyaml) // convert and thus validate the document content (-> validation here means it is valid yaml)
		debuglog("content of jsonified "+global_sourcepath, documentjson)
	}

	if !global_source_only { // if not only source (-> schema) is provided, schema and/or schema validation needs to be handled

		//
		// The yaml-formatted document is now converted to json, which means, it contains valid json after conversion.
		// Next comes the loading, conversion and verification of the schema-file
		//

		schemajson := ""        // declaration is required out of if-statement, as it is needed for further processing
		schemayaml := ""        // declaration is required out of if-statement for error tracing
		folded_schemayaml := "" // declaration is required out of if-statement for error tracing

		if !global_schemaisjson { // if schema is not provided in json -> schema is provided in yaml, then do the unfolding
			folded_schemayaml = load_file(global_schemapath) // load schema.yaml
			debuglog("content of provided "+global_schemapath, folded_schemayaml)

			schemajson = convert_yaml_document_to_json(folded_schemayaml) // convert schema from yaml to json and validate it // not used any further if no validation error occurs

			//
			// The yaml-formatted schema was now converted to json, which means, it contains valid json after conversion.
			// Next there can be two options:
			//   The schema is provided in json, which means there is nothing much todo concerning conversion and unfolding of the features (or "schema extensions") introduced by this software.
			//   The schema is provided in yaml, which means conversion must happen and unfolding should happen if enabled.
			//

			if !global_disable_unfolding { // if unfolding is not disabled
				schemayaml = unfold_schema(folded_schemayaml) // unfold own schemayaml-functions
				debuglog("unfolded "+global_schemapath, schemayaml)
			}

			schemajson = convert_yaml_document_to_json(schemayaml) // convert unfolded schema from yaml to json and validate it
			debuglog("content of jsonified "+global_schemapath, schemajson)

		} else { // schema is provided in json -> don't do unfolding
			schemajson = load_file(global_schemapath) // load schema.json from file
			debuglog("content of "+global_schemapath, schemajson)

			// making sure the provided json schema file contains valid json
			var validationvariable map[string]interface{}                  // not further used
			err := json.Unmarshal([]byte(schemajson), &validationvariable) // parse filecontent from json to map // map is not further used, but required for method call
			if err != nil {                                                // if error is thrown
				log.Fatal("error: ", err) // throw error
			}
		}

		//
		// The schema is now formatted in json.
		// Next comes the verification whether the schema is a valid json schema.
		//

		json_validationerrors := validate_json(schemajson, global_officialschemapath) // validate schema against officialschema

		if json_validationerrors != nil { // if validation is unsuccessful
			// Schema is not a valid json schema.
			// If schema was provided in json, display line number of error in json file
			// else (schema was provided in yaml) display line number of error in yaml file
			for _, json_validationerror := range json_validationerrors { // for each validationerror
				if global_schemaisjson {
					// schema was provided in json
					error_in_json_path(json_validationerror, schemajson)
				} else {
					// schema was provided in yaml
					error_in_yaml_path(json_validationerror, schemayaml)
				}
			}
		} else { // everything is fine
			debuglog("", "jsonified (unfolded) schema was successfully validated against provided metaschema at "+global_officialschemapath)
		}

		//
		// The yaml-formatted schema was now unfolded and
		// again converted to json and verified against the official (core-)json-schema, which means, it contains valid json after conversion.
		// It was also verified against the official json-meta-schema, which means, it contains a valid json-schema after conversion.
		// Next comes either
		//   the validation of the document against the custom schema or
		//   if only a schema is provided, output the (unfolded) jsonified schema
		//

		if !global_schema_only { // if only schema is provided and processed, don't do the validation

			json_validationerrors := validate_json(documentjson, schemajson) // validate document against schema

			if json_validationerrors != nil { // if validation is unsuccessful
				// document is not valid by schema definitions
				for _, json_validationerror := range json_validationerrors { // for each validationerror
					error_in_yaml_path(json_validationerror, documentyaml)
				}
			} else { // everything is fine
				debuglog("", global_sourcepath+" was successfully validated against provided "+global_schemapath)
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

		generatedschemayaml := generate_schema(documentyaml) // generate example schema
		fmt.Println(generatedschemayaml)                     // just print the generated schema
	}
}
