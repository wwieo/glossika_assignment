package flags

import (
	"flag"
)

var (
	configPath = flag.String("config", "../../../config/localTest.json", "config path")
)

func GetConfigPath() string {
	return *configPath
}

func Parse() {
	if flag.Parsed() {
		return
	}

	flag.Parse()
}
