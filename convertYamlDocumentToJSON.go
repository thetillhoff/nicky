package main

import (
	"encoding/json"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func convertYamlDocumentToJSON(yamldocument string) string {

	if len(strings.Split(yamldocument, "\n---\n")) > 1 { // if yamldocument contains more than one yaml-document
		log.Fatal("Only one yaml-document can be verified against a schema at a time.") // throw error
	}

	documentMap := make(map[string]interface{})

	err := yaml.Unmarshal([]byte(yamldocument), &documentMap) // store yaml into map
	if err != nil {                                           // if error was thrown
		log.Fatal(err) // throw the error
	}

	//debuglog("yaml-content to golang-map resulted in", documentMap)

	// jsondocument, err := json.Marshal(documentMap) // convert map to json, but unbeautified
	jsondocument, err := json.MarshalIndent(documentMap, "", "  ") // convert map to json, but beautify it with default indentation of "  "
	if err != nil {                                                // if error was thrown
		log.Fatal(err) // throw the error
	}

	return string(jsondocument)
}
