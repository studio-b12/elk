package elk_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/studio-b12/elk"
)

var (
	ErrorReadFile      = elk.ErrorCode("files:failed-reading-file")
	ErrorParsingConfig = elk.ErrorCode("config:failed-parsing")
	ErrorReadingConfig = elk.ErrorCode("config:failed-reading")
)

func readFile() ([]byte, error) {
	data, err := os.ReadFile("does/not/exist")
	if err != nil {
		return nil, elk.Wrap(ErrorReadFile, err, "failed reading file")
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
			elk.Wrap(ErrorReadFile, err, "failed reading config file")
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return configModel{},
			elk.Wrap(ErrorParsingConfig, err, "failed parsing config data")
	}

	return cfg, nil
}

func Example_detailedError() {
	_, err := parseConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("config parsing failed: %s", err)
	}
}
