package whoops_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/studio-b12/whoops"
)

var (
	ErrorReadFile      = errors.New("files:failed-reading-file")
	ErrorParsingConfig = errors.New("config:failed-parsing")
	ErrorReadingConfig = errors.New("config:failed-reading")
)

func readFile() ([]byte, error) {
	data, err := os.ReadFile("does/not/exist")
	if err != nil {
		return nil, whoops.WrapMessage(ErrorReadFile, "failed reading file", err)
	}
	return data, nil
}

type configModel struct {
	BindAddress string
	LogLevel    int
}

func parseConfig() (cfg configModel, err error) {
	data, err := readFile()
	if err != nil {
		return configModel{},
			whoops.WrapMessage(ErrorReadFile, "failed reading config file", err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return configModel{},
			whoops.WrapMessage(ErrorParsingConfig, "failed parsing config data", err)
	}

	return cfg, nil
}

func Example_detailedError() {
	_, err := parseConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatal("config parsing failed:", whoops.Format(err))
	}
}
