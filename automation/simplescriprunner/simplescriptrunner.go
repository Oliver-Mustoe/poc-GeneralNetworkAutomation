package main

import (
	"flag"
	"fmt"
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

type login struct {
	Username string `yaml:"username"`
	// Having password here is temp and needs to be probably encrypted coming in that gets decrypted throughout
	Password string `yaml:"password"`
	Msg      string `yaml:"msg"`
}

func check(errormsg error) {
	if errormsg != nil {
		log.Fatal(errormsg)
	}
}

func main() {
	// Grab flags
	var fileArg string
	// Temp default value of the nokia file
	flag.StringVar(&fileArg, "file", "nokia.yml", "File that contains login and a hello msg")
	// Parse flags
	flag.Parse()

	// Open and check file contents
	fileContentsPtr, err := os.ReadFile(fileArg)
	check(err)
	//defer fileContentsPtr.Close()

	// Read contents into YAML parser -
	var yamlContents login
	yaml.Unmarshal([]byte(fileContentsPtr), &yamlContents)
	// Print out vars
	fmt.Println("Username to use is: ", yamlContents.Username)
	fmt.Println("Password to use is: ", yamlContents.Password)
	fmt.Println("Msg to use is: '", yamlContents.Msg, "'")
}
