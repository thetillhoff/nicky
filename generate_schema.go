package main

import (
	"log"
	"reflect"

	"gopkg.in/yaml.v3"
)

func generate_schema(source string) string {

	/// the following block starts recursion over the yaml file and generates an example schema, for more information see method generate_recursive
	mapped_source := make(map[string]interface{}, 0)          // create empty map
	yaml.Unmarshal([]byte(source), &mapped_source)            // store yaml into map
	mapped_source = generate_recursive_map("", mapped_source) // start recursion
	generated_schema, err := yaml.Marshal(&mapped_source)     // store map into yaml
	if err != nil {                                           // if error is thrown
		log.Fatal(err) // throw error
	}

	return string(generated_schema) // return yaml
}

func generate_recursive_map(parentkey string, m map[string]interface{}) map[string]interface{} {

	new_map := make(map[string]interface{}) // initialize empty map

	for key, value := range m { // search all keys within m
		switch value.(type) { // compare type of value
		case bool:
			new_map[key] = "boolean" // set variable type to boolean
		case int:
			new_map[key] = "integer" // set variable type to integer
		case string:
			string_map := make(map[string]interface{})
			string_map["pattern"] = value
			new_map[key] = string_map
		case map[string]interface{}: // object
			if mapped_value, ok := value.(map[string]interface{}); ok { // if value is a map/object
				object_map := make(map[string]interface{})
				object_map["properties"] = generate_recursive_map(key, mapped_value) // start recursion and set value to result
				new_map[key] = object_map
			}
		case []interface{}: //list
			if interfaced_value, ok := value.([]interface{}); ok { // if value is a slice/list/array
				new_map[key] = generate_recursive_interface(key, interfaced_value) // start recursion and add new map to local map
			}
		default:
			log.Fatal("The type '" + reflect.TypeOf(value).String() + "' is not implemented. This problem occured on a key with the name '" + key + "'.") //unkown type error
		}
	}

	return new_map
}

func generate_recursive_interface(parentkey string, i []interface{}) map[string]interface{} {

	//new_map := make(map[string]interface{})      // initialize empty map
	new_interface := make([]interface{}, 0) // initialize empty map

	for _, value := range i { // search all keys within i
		switch value.(type) { // compare type of value
		case bool:
			new_map := make(map[string]interface{})        // initialize empty map
			new_map["type"] = "boolean"                    // set variable type to boolean
			new_interface = append(new_interface, new_map) // add new map to local interface
		case int:
			new_map := make(map[string]interface{})        // initialize empty map
			new_map["type"] = "integer"                    // set variable type to integer
			new_interface = append(new_interface, new_map) // add new map to local interface
		case string:
			string_map := make(map[string]interface{})
			//string_map["type"] = "string"                     // set variable type to string ///TODO: optional?
			string_map["pattern"] = value                     // set pattern to value to match exactly
			new_interface = append(new_interface, string_map) // add new map to local interface
		case map[string]interface{}: // object
			if mapped_value, ok := value.(map[string]interface{}); ok { // if value is a map/object
				object_map := make(map[string]interface{})
				object_map["properties"] = generate_recursive_map("", mapped_value) // start recursion and set value to result
				new_interface = append(new_interface, object_map)
			}
		case []interface{}: //list
			if interfaced_value, ok := value.([]interface{}); ok { // if value is a slice/list/array
				new_interface = append(new_interface, generate_recursive_interface("", interfaced_value)) // start recursion and add new map to local interface
			}
		default:
			log.Fatal("The type '" + reflect.TypeOf(value).String() + "' is not implemented. This problem occured on a interface with the parent named '" + parentkey + "'.") //unkown type error
		}
	}

	oneof_map := make(map[string]interface{}) // initialize empty map
	oneof_map["oneOf"] = new_interface        // make new interface a member of a mapped key "oneOf"

	items_map := make(map[string]interface{}) // initialize empty map
	items_map["items"] = oneof_map            // make oneof-map a member of a mapped key "items"

	return items_map
}
