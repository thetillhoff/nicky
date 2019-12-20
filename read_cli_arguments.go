package main

import (
	"flag"
	"log"
	"strconv"
)

func read_cli_arguments() (string, string, string, bool, bool, bool, bool, bool) {
	sourcepath := ""           // initialize variable for source file path with ""
	schemapath := ""           // initialize variable for own schema path with ""
	noofficialschema := false  // initialize variable for disabling schemaverification with false
	officialschemapath := ""   // initialize variable for official schema path with ""
	schema_only := false       // initialize variable for only conversion of schema to json with false
	source_only := false       // initialize variable for only generation an example schema from a provided source file
	schemaisjson := false      // initialize variable for disabling schema conversion (and unfolding)
	disable_unfolding := false // initialize variable for disabling unfolding feature with false
	debug := false             // initialize variable for debugging mode with false
	help := false              // initialize variable for help text with false

	flag.StringVar(&sourcepath, "source", "", "the path to the source-file") // read flag -source, without the need to dereference the pointer, default to ""
	flag.StringVar(&schemapath, "schema", "", "the path to the schema-file") // read flag -schema, without the need to dereference the pointer, default to ""
	flag.BoolVar(&noofficialschema, "noofficialschema", false, "[optional] disable verification of schema against json-meta-schema")
	flag.StringVar(&officialschemapath, "officialschema", "", " [optional] path to local json-meta-schema validation file") // read flag for official schema file, default to http url
	//flag.StringVar(&officialschemapath, "officialschema", "https://json-schema.org/draft-07/schema#", " [optional] path to local json-meta-schema validation file") // read flag for official schema file, default to http url
	flag.BoolVar(&disable_unfolding, "disable-unfolding", false, "[optional] set flag, if you want to disable the unfolding feature")
	flag.BoolVar(&debug, "debug", false, "[optional] set flag, if you need a more verbose output while running")
	flag.BoolVar(&help, "help", false, "[optional] set flag, if you need help")
	flag.Parse()        // actually get all command line flags, which were described up to this point
	tail := flag.Args() // get all remaining arguments, which do not belong to any flag
	if len(tail) > 0 {  // if there are unparsed arguments
		showhelp()
		log.Fatal("Too many arguments:\n", tail) // raise an argument error
	}

	/// check for validity of provided flag combination of 'source' and 'schema' and set variables accordingly
	if sourcepath != "" && schemapath != "" {
		// validate provided file against provided schema
	} else if sourcepath == "" && schemapath != "" { // only work on schema
		// no sourcepath, but schemapath -> set schema_only true
		schema_only = true
	} else if sourcepath != "" && schemapath == "" {
		// sourcepath, but no schemapath -> set source_only true
		source_only = true
	} else { // unexpected flag combination
		showhelp()                                                                // show help
		log.Fatal("Unallowed flag combination concerning 'source' and 'schema'.") // throw error
	}

	if !source_only { // only if schema is provided
		/// check for validity of provided flag combination of 'officialschema' and 'noofficialschema' and set variables accordingly
		if officialschemapath == "" && noofficialschema { // validate schema against official meta-schema
		// no officialschemapath set and noofficialschema -> set officialschemapath to ""
		} else if officialschemapath != "" && !noofficialschema { // do not validate schema against official meta-schema
		// officialschema set and no noofficalschema -> no need to interfere
		} else { // unexpected flag combination
			showhelp()                                                                                  // show help
			log.Fatal("Unallowed flag combination concerning 'officialschema' and 'noofficialschema'.") // throw error
		}

		/// detect schema-file-extension for [yaml|json] to disable conversion on [json]
		if schemapath[len(schemapath)-5:] == ".json" { // if schema is in json format
			schemaisjson = true
		} else if schemapath[len(schemapath)-5:] == ".yaml" || schemapath[len(schemapath)-4:] == ".yml" {
			// schema is in yaml format
		} else { // unkown file extension of schema
			log.Fatal("Unkown file-extension on --schema. Allowed extensions are .yaml, .yml and .json") // throw error
		}
	} else { // no schema was provided
		if disable_unfolding { // if no schema was provided but unfolding was disabled via flag
			log.Print("Warning: Strange flag combination. If no schema is provided, no schema is processed and thus no unfolding will take place.")
		}
	}

	/// still checking for validity of provided flag combination but on a warning level
	if disable_unfolding && schemaisjson { // if schema is json and unfolding is disabled via flag
		log.Print("Warning: Strange flag combination. If schema is provided in json, unfolding is skipped by default.") // print warning
	}

	if help { // if helpflag is set
		showhelp() // show help
	}

	return sourcepath, schemapath, officialschemapath, schema_only, source_only, schemaisjson, disable_unfolding, debug
}

// This method is only for debugging. It prints the formatted flag values.
func print_cli_arguments() {
	values := ""
	values = values + "source-path:           "
	values = values + global_sourcepath + "\n"
	values = values + "schema-path:           "
	values = values + global_schemapath + "\n"
	values = values + "officialschema-path:   "
	values = values + global_officialschemapath + "\n"
	values = values + "schema-only-boolean:       "
	values = values + strconv.FormatBool(global_schema_only) + "\n"
	values = values + "source-only-boolean:       "
	values = values + strconv.FormatBool(global_source_only) + "\n"
	values = values + "disable-unfolding-boolean: "
	values = values + strconv.FormatBool(global_disable_unfolding) + "\n"
	values = values + "schemaisjson-boolean:      "
	values = values + strconv.FormatBool(global_schemaisjson) + "\n"
	values = values + "debug-boolean:             "
	values = values + strconv.FormatBool(global_debug) + ""
	debuglog("global variables:", values)
}
