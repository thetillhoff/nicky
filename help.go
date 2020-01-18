package main

import (
	"fmt"
)

// This method is just for printing a formatted help-text.
func showhelp() {
	fmt.Printf("Nicky is a tool for schema verification of yaml files.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Println("  $> nicky --source <filename>.yaml --schema <schemafilename>.yaml [--noofficialschema|--officialschema [http[s]://|file://|/|./|../]<path>] [--disable-unfolding] [--debug]")
	fmt.Println("  $> nicky --source <filename>.yaml [--debug]")
	fmt.Println("  $> nicky --schema <schemafilename>.yaml [--noofficialschema|--officialschema [http[s]://|file://|/|./|../]<path>] [--disable-unfolding] [--debug]")
	fmt.Printf("\nPlease note, that the path to the officialschema (if provided) must look like one of the following.\nExamples:\n")
	fmt.Println("- url\n  http://path/to/file")
	fmt.Println("- unix uri\n  file:///path/to/file")
	fmt.Println("- windows uri\n  file://H:/path/to/file")
	fmt.Println("- path\n  [/|./|../]path/to/file")
	fmt.Println("") // empty line for better readability
}
