package main

import (
	"log"

	"gopkg.in/yaml.v3"
)

func error_in_yaml_path(json_validationerror JSONValidationError, document string) {

	var document_node yaml.Node
	err := yaml.Unmarshal([]byte(document), &document_node)
	if err != nil { // if error was thrown
		log.Fatal(err) // throw the error
	}

	// fmt.Println("linenumber: ", yaml_pathfinder(json_valdiationerror.Field(), document_node))

	log.Fatal("\"", json_validationerror.errortype, "\"-error in yaml at object\n  ", json_validationerror.field, ": ", json_validationerror.description) // throw error

	///TODO note that linenumber is taken after unfolding, so there might be a problem with the error mapping
}

// func yaml_pathfinder(path string, document_node yaml.Node) int {
// 	for index := range document_node.Content {
// 		if document_node[index].Value == strings.Split(path, ".")[0] { // correct childnode
// 		}
// 	}
// }
