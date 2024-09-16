package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Parses the config json and decodes it into an arbitrary struct.
//
// Panics if an error is encountered.
//
// Additionally, a validation is performed on the structs' fields, with the following characteristics:
//
// - Fields cannot be zero-valued, unless explicitly allowed to be
//
// - `value` must be a `struct` or its pointer
//
// - The following tags are accepted:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> for JSONField only, allows field to not be set
//	`nullable:"true"` -> for JSONField only, allows field to be set to null
//	`zeroable:"true"` -> for JSONField only, allows field to be set to zero-value
func Parse[T any](workingDir, configEnvar string, configStruct *T) {
	filePath := os.Getenv(configEnvar)
	if filePath == "" {
		fmt.Printf("$%s not set, using default config\n", configEnvar)
		filePath = filepath.Join(workingDir, "default.json")
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprint("could not open config file: ", err))
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(configStruct); err != nil {
		panic(fmt.Sprint("could not parse config file: ", err))
	}

	fields, err := structFromJSON(configStruct)
	if err != nil {
		panic(err)
	}
	if len(fields) != 0 {
		panic(fmt.Sprint("missing required config field(s): ", fields))
	}
}
