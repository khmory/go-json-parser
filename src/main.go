package main

import (
	"fmt"
	"github.com/koheimorii/go-json-parser/src/parser"
	"io/ioutil"
	"os"
)

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("please input json file from stdin")
	}
	json := (string(body))
	parsed,err := parser.Parse(json)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(parsed)
}
