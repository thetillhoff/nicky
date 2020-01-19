package main

import (
	"log"

	"gopkg.in/yaml.v3"
)

func errorInYAMLPath(jsonValidationerror JSONValidationError, document string) {

	var documentNode yaml.Node
	err := yaml.Unmarshal([]byte(document), &documentNode)
	if err != nil { // if error was thrown
		log.Fatal(err) // throw the error
	}

	// fmt.Println("linenumber: ", yaml_pathfinder(json_valdiationerror.Field(), documentNode))

	log.Fatal("\"", jsonValidationerror.errortype, "\"-error in yaml at object\n  ", jsonValidationerror.field, ": ", jsonValidationerror.description) // throw error

	///TODO note that linenumber is taken after unfolding, so there might be a problem with the error mapping
}

// func yaml_pathfinder(path string, documentNode yaml.Node) int {
// 	for index := range documentNode.Content {
// 		if documentNode[index].Value == strings.Split(path, ".")[0] { // correct childnode
// 		}
// 	}
// }
