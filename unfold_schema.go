package main

import (
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func unfold_schema(schema string) string {

	/// The following block unfolds what is called "beginning already included" in the documentation.md.
	/// It autofills the fields $id, $schema, title and type at the beginning of the schema file.
	lines := strings.Split(schema, "\n")                                    // split string to lines
	if lines[0] == "---" && len(lines[1]) > 6 && lines[1][:6] == "name: " { // if beginning of schema is folded
		name := lines[1][6:]       // retrieve schema name
		lines = lines[2:]          // remove first two lines (the folded beginning)
		var prepend_lines []string // initialize variable

		prepend_lines = append(prepend_lines, "---") // readded for better visibility how the file is started
		prepend_lines = append(prepend_lines, "$id: "+name+".schema.json")
		prepend_lines = append(prepend_lines, "$schema: 'http://json-schema.org/draft-07/schema#'")
		prepend_lines = append(prepend_lines, "title: "+name)
		prepend_lines = append(prepend_lines, "type: object") // not really required as it is added by what is called "object notation" in documentation.md

		lines = append(prepend_lines, lines...) // add the unfolded beginning of schema
	}
	schema = strings.Join(lines, "\n") // concat lines to string

	/// the following block starts recursion over the yaml file and unfold several functionalities, for more information see method unfold_recursive
	mapped_schema := make(map[string]interface{}, 0)     // create empty map
	yaml.Unmarshal([]byte(schema), &mapped_schema)       // store yaml into map
	mapped_schema = unfold_recursive("", mapped_schema)  // start recursion
	unfolded_schema, err := yaml.Marshal(&mapped_schema) // store map into yaml
	if err != nil {                                      // if error is thrown
		log.Fatal(err) // throw error
	}

	return string(unfolded_schema) // return yaml
}

func unfold_recursive(parentkey string, m map[string]interface{}) map[string]interface{} {

	/// The following block unfolds what is called "object notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field properties is present.
	/// It also adds the field additionalProperties corresponding to the presence of the property ...:.
	has_properties := false           // initialize searchresult with not found
	has_type := false                 // initialize searchresult with not found
	has_additionalproperties := false // initialize searchresult with not found
	for key, _ := range m {           // search all keys within m
		if key == "properties" { // if key "properties" is found
			has_properties = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			has_type = true // set searchresult to found
		} else if key == "additionalProperties" { // if key "additionalProperties" is found
			has_additionalproperties = true // set searchresult to found
		}
	}
	if has_properties && !has_type { // if m has properties but no type
		m["type"] = "object" // add type manually
	}
	if !has_additionalproperties && has_properties { // if m does not have the key additionalproperties but should have it
		allows_additionalproperties := false                           // initialize searchresult with not found
		for key, _ := range m["properties"].(map[string]interface{}) { // search all keys within m
			if key == "..." { // if key "..." is found in properties
				allows_additionalproperties = true                    // set searchresult to found
				delete(m["properties"].(map[string]interface{}), key) // delete the key, as it is not of use any more
			}
		}
		if allows_additionalproperties { // if additionalProperties should be allowed
			m["additionalProperties"] = true // add the key-value accordingly
		} else { // if additionalproperties should be denied
			m["additionalProperties"] = false // add the key-value accordingly
		}
	}

	/// The following block unfolds what is called "array notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field items is present.
	/// It also adds the field additionalItems corresponding to the presence of the property ...:
	has_items := false           // initialize searchresult with not found
	has_type = false             // initialize searchresult with not found
	has_additionalitems := false // initialize searchresult with not found
	for key, _ := range m {      // search all keys within m
		if key == "items" { // if key "items" is found
			has_items = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			has_type = true // set searchresult to found
		}
	}
	if has_items && !has_type { // if m has items but no type
		m["type"] = "array" // add type manually
	}
	if !has_additionalitems && has_items { // if m does not have the key additionalitems but should have it
		allows_additionalitems := false                           // initialize searchresult with not found
		for key, _ := range m["items"].(map[string]interface{}) { // search all keys within m
			if key == "..." { // if key "..." is found in items
				allows_additionalitems = true                    // set searchresult to found
				delete(m["items"].(map[string]interface{}), key) // delete the key, as it is not of use any more
			}
		}
		if allows_additionalitems { // if additionalitems should be allowed
			m["additionalItems"] = true // add the key-value accordingly
		} else { // if additionalitems should be denied
			m["additionalItems"] = false // add the key-value accordingly
		}
	}

	/// The following block unfolds what is called "pattern notation" in the unfolding_documentation.md.
	/// It adds the field type automatically if the field pattern is present.
	has_pattern := false    // initialize searchresult with not found
	has_type = false        // initialize searchresult with not found
	for key, _ := range m { // search all keys within m
		if key == "pattern" { // if key "items" is found
			has_pattern = true // set searchresult to found
		} else if key == "type" { // if key "type" is found
			has_type = true // set searchresult to found
		}
	}
	if has_pattern && !has_type { // if m has items but no type
		m["type"] = "string" // add type manually
	}

	///TODO: implement additional unfolding mechanisms (for maps) here

	for key, value := range m { // for each key-value pair in m

		///TODO: implement additional unfolding mechanisms (for primitives) here

		if mapped_value, ok := value.(map[string]interface{}); ok { // if value is a map/object
			m[key] = unfold_recursive(key, mapped_value) // start recursion and set value to result
		} else if stringed_value, ok := value.(string); ok { // if value is string
			/// The following block unfolds what is called "simplified primitive type declaration" in the unfolding_documentation.md.
			/// It removes the additional layer on primitive type declaration and adds it automatically.
			if parentkey == "properties" || parentkey == "items" { // only do this when element is subkey of object or array
				primitives := []string{"null", "boolean", "number", "integer", "string", "object", "array"} // list of primitive datatypes which should be unfolded
				if sliceContains(primitives, stringed_value) {                                              // if the single subkey is "type" and its value is a primitive type
					new_map := make(map[string]interface{}) // initialize empty map
					new_map["type"] = stringed_value        // set unfolded type to folded type
					m[key] = new_map                        // set map as new value
				}
			}
		}
	}

	return m
}

func sliceContains(list []string, compare_element string) bool {
	for _, element := range list {
		if element == compare_element {
			return true
		}
	}
	return false
}
