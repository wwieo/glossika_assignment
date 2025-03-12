package config

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"glossika/service/internal/flags"
	"os"
)

func NewGlossika() Glossika {
	glossika := Glossika{}
	buildConfig(&glossika)
	return glossika
}

func buildConfig(cfg interface{}) {
	var err error
	path := flags.GetConfigPath()

	viper.SetConfigType("json")
	viper.SetConfigFile(path)

	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var rawConfig map[string]interface{}
	if err = json.Unmarshal(fileContent, &rawConfig); err != nil {
		panic(fmt.Sprintf("Error unmarshalling JSON: %v", err))
	}
	decoderConfigOption := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
		),
		Result: cfg,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfigOption)
	if err != nil {
		panic(fmt.Sprintf("Error creating decoder: %v", err))
	}

	if err = decoder.Decode(rawConfig); err != nil {
		panic(fmt.Sprintf("Error decoding config: %v", err))
	}
}
