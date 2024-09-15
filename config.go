package gocommon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func InitialiseConfig[T any](workingDir, configEnvar string) *T {
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

	config := new(T)
	if err = json.NewDecoder(file).Decode(config); err != nil {
		panic(fmt.Sprint("could not parse config file: ", err))
	}

	return config
}
