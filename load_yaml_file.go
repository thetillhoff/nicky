package main

import (
	"io/ioutil"
	"log"
)

func load_file(filepath string) string {

	if filepath[:1] != "/" && filepath[:2] != "./" && filepath[:3] != "../" { // if filepath is no absolute filepath and no full relative path
		filepath = "./" + filepath // add current folder as relative path
	}

	chart_content, err := ioutil.ReadFile(filepath) // read contents from path into variable
	if err != nil {                                 // if error was thrown
		log.Fatal(err, filepath) // throw the error
	}
	return string(chart_content)
}
