package main

import (
	"flag"
	"log"
	"strconv"
)

func read_cli_arguments() (string, string, string, bool, bool, bool, bool) {
	filepath := ""             // initialize variable for source file path with ""
	schemapath := ""           // initialize variable for own schema path with ""
	officialschemapath := ""   // initialize variable for official schema path with ""
	schema_only := false       // initialize variable for only conversion of schema to json with false
	schemaisjson := false      // initialize variable for disabling schema conversion (and unfolding)
	disable_unfolding := false // initialize variable for disabling unfolding feature with false
	debug := false             // initialize variable for debugging mode with false
	help := false              // initialize variable for help text with false

	flag.StringVar(&filepath, "source", "", "the path to the source-file")                                                                                          // read flag -source, without the need to dereference the pointer, default to ""
	flag.StringVar(&schemapath, "schema", "", "the path to the schema-file")                                                                                        // read flag -schema, without the need to dereference the pointer, default to ""
	flag.StringVar(&officialschemapath, "officialschema", "https://json-schema.org/draft-07/schema#", " [optional] path to local json-meta-schema validation file") // read flag for official schema file, default to http url
	flag.BoolVar(&disable_unfolding, "disable-unfolding", false, "[optional] set flag, if you want to disable the unfolding feature.")
	flag.BoolVar(&debug, "debug", false, "[optional] set flag, if you need a more verbose output while running.")
	flag.BoolVar(&help, "help", false, "[optional] set flag, if you need help.")
	flag.Parse()        // actually get all command line flags, which were described up to this point
	tail := flag.Args() // get all remaining arguments, which do not belong to any flag
	if len(tail) > 0 {  // if there are unparsed arguments
		showhelp()
		log.Fatal("Too many arguments:\n", tail) // raise an argument error
	}

	/// check for validity of provided flag combination and set variables accordingly
	if filepath != "" && schemapath != "" && !schema_only {
		// validate provided file against provided schema
	} else if filepath == "" && schemapath != "" { // only work on schema
		// no filepath, but schemapath -> set schema_only true
		schema_only = true
	} else { // unexpected flag combination
		showhelp()                               // show help
		log.Fatal("Unallowed flag combination.") // throw error
	}

	/// still checking for validity of provided flag combination but on a warning level
	if disable_unfolding && schemaisjson { // if schema is json and unfolding is disabled via flag
		log.Print("Warning: Strange flag combination. If schema is provided in json, unfolding is skipped by default.") // print warning
	}

	if help { // if helpflag is set
		showhelp() // show help
	}

	/// detect schema-file-extension for [yaml|json] to disable conversion on [json]
	if schemapath[len(schemapath)-5:] == ".json" { // if schema is in json format
		schemaisjson = true
	} else if schemapath[len(schemapath)-5:] == ".yaml" || schemapath[len(schemapath)-4:] == ".yml" {
		// schema is in yaml format
	} else { // unkown file extension of schema
		log.Fatal("Unkown file-extension on --schema. Allowed extensions are .yaml, .yml and .json") // throw error
	}

	return filepath, schemapath, officialschemapath, schema_only, schemaisjson, disable_unfolding, debug
}

// This method is only for debugging. It prints the formatted flag values.
func print_cli_arguments() {
	values := ""
	values = values + "source-filepath:           "
	values = values + global_filepath + "\n"
	values = values + "schema-filepath:           "
	values = values + global_schemapath + "\n"
	values = values + "officialschema-filepath:   "
	values = values + global_officialschemapath + "\n"
	values = values + "schema-only-boolean:       "
	values = values + strconv.FormatBool(global_schema_only) + "\n"
	values = values + "disable-unfolding-boolean: "
	values = values + strconv.FormatBool(global_disable_unfolding) + "\n"
	values = values + "schemaisjson-boolean:      "
	values = values + strconv.FormatBool(global_schemaisjson) + "\n"
	values = values + "debug-boolean:             "
	values = values + strconv.FormatBool(global_debug) + ""
	debuglog("global variables:", values)
}
