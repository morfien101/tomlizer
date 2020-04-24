package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

var (
	version = "0.0.0"

	flagVersion    = flag.Bool("v", false, "Show the version.")
	flagHelp       = flag.Bool("h", false, "Show this help menu.")
	flagSTDIN      = flag.Bool("s", false, "Read from STDIN, ignores env vars input if set.")
	flagSTDOUT     = flag.Bool("o", false, "Write to STDOUT. Ignores file output if set.")
	flagOutputFile = flag.String("f", "", "Output to path.")
	flagEnvName    = flag.String("e", "", "Read from environment vars.")
)

func showStoppers() {
	if *flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if *flagHelp {
		fmt.Println(helpMessage())
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func helpMessage() string {
	return `
This application is used to read in JSON and return TOML.
You can use to create toml configuration files from JSON
structures that are easy to store.
The input MUST be valid JSON.
`
}

func main() {
	flag.Parse()
	showStoppers()

	var jsonInput []byte

	if *flagSTDIN {
		jsonBytes, err := readSTDIN()
		if err != nil {
			terminate(err.Error(), 1)
		}
		if len(jsonBytes) == 0 {
			terminate("No input detected", 1)
		}
		jsonInput = jsonBytes
	} else {
		if len(*flagEnvName) == 0 {
			terminate("No input source given", 1)
		}
		jsonBytes, err := readEnv(*flagEnvName)
		if err != nil {
			terminate(err.Error(), 1)
		}
		jsonInput = jsonBytes
	}

	tomlString, err := jSONreader(jsonInput)
	if err != nil {
		terminate(err.Error(), 1)
	}
	if *flagSTDOUT {
		fmt.Println(tomlString)
	} else {
		if len(*flagOutputFile) == 0 {
			terminate("No output method given", 1)
		}
		err = writeToFile([]byte(tomlString), *flagOutputFile)
		if err != nil {
			terminate(err.Error(), 1)
		}
	}
}

func terminate(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func readEnv(varName string) ([]byte, error) {
	value, ok := os.LookupEnv(varName)
	if !ok {
		return nil, fmt.Errorf("Could not find environment variable named %s", varName)
	}
	return []byte(value), nil
}

func readSTDIN() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

func writeToFile(output []byte, path string) error {
	return ioutil.WriteFile(path, output, 0640)
}

func jSONreader(jsonBytes []byte) (string, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return "", err
	}

	tree, err := toml.TreeFromMap(jsonMap)
	if err != nil {
		return "", err
	}
	return mapToTOML(tree)
}

func mapToTOML(t *toml.Tree) (string, error) {
	tomlBytes, err := t.ToTomlString()
	if err != nil {
		return "", err
	}
	return string(tomlBytes), nil
}
