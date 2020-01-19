package main

import (
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func unfoldSchema(schema string) string {

	/// The following block unfolds what is called "beginning already included" in the documentation.md.
	/// It autofills the fields $id, $schema, title and type at the beginning of the schema file.
	lines := strings.Split(schema, "\n")                                    // split string to lines
	if lines[0] == "---" && len(lines[1]) > 6 && lines[1][:6] == "name: " { // if beginning of schema is folded
		name := lines[1][6:]      // retrieve schema name
		lines = lines[2:]         // remove first two lines (the folded beginning)
		var prependLines []string // initialize variable

		prependLines = append(prependLines, "---") // readded for better visibility how the file is started
		prependLines = append(prependLines, "$id: "+name+".schema.json")
		prependLines = append(prependLines, "$schema: 'http://json-schema.org/draft-07/schema#'")
		prependLines = append(prependLines, "title: "+name)
		prependLines = append(prependLines, "type: object") // not really required as it is added by what is called "object notation" in documentation.md

		lines = append(prependLines, lines...) // add the unfolded beginning of schema
	}
	schema = strings.Join(lines, "\n") // concat lines to string

	/// the following block starts recursion over the yaml file and unfold several functionalities, for more information see method unfoldRecursive
	mappedSchema := make(map[string]interface{}, 0)    // create empty map
	yaml.Unmarshal([]byte(schema), &mappedSchema)      // store yaml into map
	mappedSchema = unfoldRecursive("", mappedSchema)   // start recursion
	unfoldedSchema, err := yaml.Marshal(&mappedSchema) // store map into yaml
	if err != nil {                                    // if error is thrown
		log.Fatal(err) // throw error
	}

	return string(unfoldedSchema) // return yaml
}

func unfoldRecursive(parentkey string, m map[string]interface{}) map[string]interface{} {

	/// The following block unfolds what is called "object notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field properties is present.
	/// It also adds the field additionalProperties corresponding to the presence of the property ...:.
	hasProperties := false           // initialize searchresult with not found
	hasType := false                 // initialize searchresult with not found
	hasAdditionalproperties := false // initialize searchresult with not found
	for key := range m {             // search all keys within m
		if key == "properties" { // if key "properties" is found
			hasProperties = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			hasType = true // set searchresult to found
		} else if key == "additionalProperties" { // if key "additionalProperties" is found
			hasAdditionalproperties = true // set searchresult to found
		}
	}
	if hasProperties && !hasType { // if m has properties but no type
		m["type"] = "object" // add type manually
	}
	if !hasAdditionalproperties && hasProperties { // if m does not have the key additionalproperties but should have it
		allowsAdditionalproperties := false                         // initialize searchresult with not found
		for key := range m["properties"].(map[string]interface{}) { // search all keys within m
			if key == "..." { // if key "..." is found in properties
				allowsAdditionalproperties = true                     // set searchresult to found
				delete(m["properties"].(map[string]interface{}), key) // delete the key, as it is not of use any more
			}
		}
		if allowsAdditionalproperties { // if additionalProperties should be allowed
			m["additionalProperties"] = true // add the key-value accordingly
		} else { // if additionalproperties should be denied
			m["additionalProperties"] = false // add the key-value accordingly
		}
	}

	/// The following block unfolds what is called "array notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field items is present.
	/// It also adds the field additionalItems corresponding to the presence of the property ...:
	hasItems := false           // initialize searchresult with not found
	hasType = false             // initialize searchresult with not found
	hasAdditionalitems := false // initialize searchresult with not found
	for key := range m {        // search all keys within m
		if key == "items" { // if key "items" is found
			hasItems = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			hasType = true // set searchresult to found
		}
	}
	if hasItems && !hasType { // if m has items but no type
		m["type"] = "array" // add type manually
	}
	if !hasAdditionalitems && hasItems { // if m does not have the key additionalitems but should have it
		allowsAdditionalitems := false                         // initialize searchresult with not found
		for key := range m["items"].(map[string]interface{}) { // search all keys within m
			if key == "..." { // if key "..." is found in items
				allowsAdditionalitems = true                     // set searchresult to found
				delete(m["items"].(map[string]interface{}), key) // delete the key, as it is not of use any more
			}
		}
		if allowsAdditionalitems { // if additionalitems should be allowed
			m["additionalItems"] = true // add the key-value accordingly
		} else { // if additionalitems should be denied
			m["additionalItems"] = false // add the key-value accordingly
		}
	}

	/// The following block unfolds what is called "pattern notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field pattern is present.
	hasPattern := false  // initialize searchresult with not found
	hasType = false      // initialize searchresult with not found
	for key := range m { // search all keys within m
		if key == "pattern" { // if key "items" is found
			hasPattern = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			hasType = true // set searchresult to found
		}
	}
	if hasPattern && !hasType { // if m has items but no type
		m["type"] = "string" // add type manually
	}

	///TODO: implement additional unfolding mechanisms (for maps) here

	for key, value := range m { // for each key-value pair in m

		///TODO: implement additional unfolding mechanisms (for primitives) here

		if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map/object
			m[key] = unfoldRecursive(key, mappedValue) // start recursion and set value to result
		} else if stringedValue, ok := value.(string); ok { // if value is string
			/// The following block unfolds what is called "simplified primitive type declaration" in the unfolding_documentation.md.
			/// It removes the additional layer on primitive type declaration and adds it automatically.
			if parentkey == "properties" || parentkey == "items" { // only do this when element is subkey of object or array
				primitives := []string{"null", "boolean", "number", "integer", "string", "object", "array"} // list of primitive datatypes which should be unfolded
				if sliceContains(primitives, stringedValue) {                                               // if the single subkey is "type" and its value is a primitive type
					newMap := make(map[string]interface{}) // initialize empty map
					newMap["type"] = stringedValue         // set unfolded type to folded type
					m[key] = newMap                        // set map as new value
				}
			}
		}
	}

	return m
}

func sliceContains(list []string, compareElement string) bool {
	for _, element := range list {
		if element == compareElement {
			return true
		}
	}
	return false
}
