package main

import (
	"fmt"
)

func debuglog(name string, content interface{}) {
	if global_debug { // if global debugging is enabled
		fmt.Println("=========================") // print separator
		if name != "" {                          // if name is provided
			fmt.Println(name + ":")                  // print name and
			fmt.Println("=========================") // print separator
		}
		fmt.Println(content)                     // print content
		fmt.Println("=========================") // print separator
	}
}
